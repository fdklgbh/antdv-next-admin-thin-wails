package main

import (
	"embed"
	"encoding/json"
	"log"

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

//go:embed wails.json
var projectConfig []byte

func main() {
	var metadata struct {
		Name string `json:"name"`
		Info struct {
			ProductName string `json:"productName"`
		} `json:"info"`
	}
	if err := json.Unmarshal(projectConfig, &metadata); err != nil {
		log.Fatalf("解析内嵌 wails.json 失败: %v", err)
	}
	if metadata.Name == "" || metadata.Info.ProductName == "" {
		log.Fatal("wails.json 必须设置 name 和 info.productName")
	}
	// Create an instance of the app structure
	app := NewApp()
	if stopWindowRestore := configureWindowRestore(); stopWindowRestore != nil {
		defer stopWindowRestore()
	}

	// Create application with options
	err := wails.Run(&options.App{
		Title:     metadata.Info.ProductName,
		Width:     1024,
		Height:    768,
		Frameless: true,
		Linux: &linux.Options{
			Icon:             appIcon,
			ProgramName:      metadata.Name,
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
