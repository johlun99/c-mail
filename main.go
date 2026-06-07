package main

import (
	"embed"
	"os"
	"runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// WebKitGTK's native Wayland surface mishandles compositor-driven resizes
	// under tiling WMs (e.g. Hyprland: the window renders full size, then shrinks
	// on the next repaint). Running via XWayland avoids it. Set GDK_BACKEND
	// yourself to override.
	if runtime.GOOS == "linux" && os.Getenv("GDK_BACKEND") == "" {
		_ = os.Setenv("GDK_BACKEND", "x11")
	}

	// Create an instance of the app structure
	svc := NewMailService()
	app := NewApp(svc)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "cmail",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			svc,
		},
		// Linux/Wayland: a stable program name → stable window app_id ("cmail"),
		// so tiling WMs like Hyprland identify and tile the window consistently
		// (without it, `wails dev`'s temp binary yields an unstable app_id).
		Linux: &linux.Options{
			ProgramName: "cmail",
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
