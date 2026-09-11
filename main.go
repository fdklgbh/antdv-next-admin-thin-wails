package main

import (
	"embed"

	"log"
	"time"

	"antdv-next-admin-thin-wails/internal/system"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gopkg.in/yaml.v3"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

//go:embed build/config.yml
var projectConfig []byte

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	var metadata struct {
		Info struct {
			ProductName string `yaml:"productName"`
			Description string `yaml:"description"`
		} `yaml:"info"`
		Packaging struct {
			AppName string `yaml:"appName"`
		} `yaml:"packaging"`
	}
	if err := yaml.Unmarshal(projectConfig, &metadata); err != nil {
		log.Fatalf("解析内嵌 build/config.yml 失败: %v", err)
	}
	if metadata.Packaging.AppName == "" || metadata.Info.ProductName == "" {
		log.Fatal("build/config.yml 必须设置 packaging.appName 和 info.productName")
	}

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        metadata.Packaging.AppName,
		Description: metadata.Info.Description,
		Icon:        appIcon,
		Linux: application.LinuxOptions{
			ProgramName: metadata.Packaging.AppName,
		},
		Services: []application.Service{
			application.NewService(&system.Service{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: metadata.Info.ProductName,
		// Window sized to the golden ratio (1000 / 618 ≈ 1.618).
		Width:     1024,
		Height:    768,
		Frameless: true,
		Linux: application.LinuxWindow{
			Icon: appIcon,
			// Avoid corrupted composited layers with Linux WebKit GPU drivers.
			WebviewGpuPolicy: application.WebviewGpuPolicyNever,
		},
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})
	configureWindowRestore(window)

	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
