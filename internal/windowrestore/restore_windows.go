package windowrestore

import (
	"log"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	restoreUser32            = windows.NewLazySystemDLL("user32.dll")
	restoreSetWinEventHook   = restoreUser32.NewProc("SetWinEventHook")
	restoreUnhookWinEvent    = restoreUser32.NewProc("UnhookWinEvent")
	restoreGetMessage        = restoreUser32.NewProc("GetMessageW")
	restorePeekMessage       = restoreUser32.NewProc("PeekMessageW")
	restorePostThreadMessage = restoreUser32.NewProc("PostThreadMessageW")
	restoreSetWindowPos      = restoreUser32.NewProc("SetWindowPos")
	restoreIsIconic          = restoreUser32.NewProc("IsIconic")
)

// WinEvent callbacks require a message loop on the thread that installed them.
type restoreMessage struct {
	hwnd           uintptr
	message        uint32
	wParam, lParam uintptr
	time           uint32
	x, y           int32
	private        uint32
}

var restoreCallback = windows.NewCallback(func(_ uintptr, _ uint32, hwnd uintptr, objectID, childID int32, _ uint32, _ uint32) uintptr {
	if hwnd == 0 || objectID != 0 || childID != 0 {
		return 0
	}
	if iconic, _, _ := restoreIsIconic.Call(hwnd); iconic != 0 {
		return 0
	}
	// Recalculate the maximised client area after the restore animation has
	// finished, when MonitorFromRect can identify the actual restore monitor.
	// Keep the window's position, size, stacking order and activation unchanged.
	// SWP_ASYNCWINDOWPOS queues the work on Wails' UI thread, so this
	// listener cannot block on that thread while the application is exiting.
	const flags = 0x0001 | 0x0002 | 0x0004 | 0x0010 | 0x0020 | 0x4000
	if ok, _, err := restoreSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, flags); ok == 0 {
		log.Printf("refresh window bounds after restore: %v", err)
	}
	return 0
})

// Configure installs the restore listener and returns its shutdown function.
func Configure() func() {
	ready := make(chan uint32, 1)
	done := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		defer close(done)

		// Create the message queue before exposing the thread ID to shutdown.
		var message restoreMessage
		restorePeekMessage.Call(uintptr(unsafe.Pointer(&message)), 0, 0, 0, 0)
		const eventSystemMinimizeEnd = 0x0017
		// WINEVENT_OUTOFCONTEXT (0) delivers only this process's events here,
		// after the window's own restore processing, without subclassing Wails.
		hook, _, err := restoreSetWinEventHook.Call(eventSystemMinimizeEnd, eventSystemMinimizeEnd,
			0, restoreCallback, uintptr(windows.GetCurrentProcessId()), 0, 0)
		if hook == 0 {
			log.Printf("install window restore listener: %v", err)
			ready <- 0
			return
		}
		defer func() {
			if ok, _, err := restoreUnhookWinEvent.Call(hook); ok == 0 {
				log.Printf("remove window restore listener: %v", err)
			}
		}()
		ready <- windows.GetCurrentThreadId()
		for {
			result, _, err := restoreGetMessage.Call(uintptr(unsafe.Pointer(&message)), 0, 0, 0)
			if int32(result) == -1 {
				log.Printf("read window restore event: %v", err)
				return
			}
			if result == 0 {
				return
			}
		}
	}()
	threadID := <-ready
	if threadID == 0 {
		return nil
	}
	return func() {
		select {
		case <-done:
			return
		default:
		}
		const wmQuit = 0x0012
		if ok, _, err := restorePostThreadMessage.Call(uintptr(threadID), wmQuit, 0, 0); ok == 0 {
			log.Printf("stop window restore listener: %v", err)
			return
		}
		<-done
	}
}
