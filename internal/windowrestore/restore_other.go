//go:build !windows

package windowrestore

import "github.com/wailsapp/wails/v3/pkg/application"

// Configure is a no-op on platforms that do not need the Windows workaround.
func Configure(_ *application.WebviewWindow) {}
