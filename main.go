package main

import (
	"embed"

	"antdv-next-admin-thin-wails-v2/internal/system"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	// Create an instance of the app structure
	app := NewApp()
	if stopWindowRestore := configureWindowRestore(); stopWindowRestore != nil {
		defer stopWindowRestore()
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "Antdv Next Thin V2",
		Width:     1024,
		Height:    768,
		Frameless: true,
		Linux: &linux.Options{
			Icon:             appIcon,
			ProgramName:      "antdv-next-admin-thin-wails",
			WebviewGpuPolicy: linux.WebviewGpuPolicyNever,
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			&system.Service{},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
