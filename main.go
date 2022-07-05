package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"xBrowser/app"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/bear.png
var icon []byte

func main() {
	// Create an instance of the app structure
	app := app.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:      "xBrowser",
		Width:      1024,
		Height:     1024,
		Assets:     assets,
		OnStartup:  app.Startup,
		OnDomReady: app.DomReady,
		OnShutdown: app.Shutdown,
		LogLevel:   logger.DEBUG,
		Bind: []interface{}{
			app,
		},

		Mac: &mac.Options{
			//TitleBar:             mac.TitleBarHiddenInset(),
			//WebviewIsTransparent: true,
			//WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "Wails Template Vue",
				Message: "A Wails template based on Vue and Vue-Router",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		println("Error:", err)
	}
}
