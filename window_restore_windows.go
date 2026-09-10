package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"golang.org/x/sys/windows"
)

var setWindowPos = windows.NewLazySystemDLL("user32.dll").NewProc("SetWindowPos")

func configureWindowRestore(window *application.WebviewWindow) {
	window.OnWindowEvent(events.Windows.WindowUnMinimise, func(_ *application.WindowEvent) {
		// Run after the restore message has finished: WebView2 may have used
		// the parked, minimised window's monitor scale for its first resize.
		application.InvokeAsync(func() {
			handle := window.NativeWindow()
			if handle == nil || window.IsMinimised() {
				return
			}

			// SWP_FRAMECHANGED sends WM_NCCALCSIZE / WM_SIZE again with the
			// restored client bounds, without moving/resizing/activating the window.
			const flags = 0x0001 | 0x0002 | 0x0004 | 0x0010 | 0x0020
			ok, _, err := setWindowPos.Call(uintptr(handle), 0, 0, 0, 0, 0, flags)
			if ok == 0 {
				log.Printf("refresh window bounds after restore: %v", err)
			}
		})
	})
}
