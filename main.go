package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"xBrowser/app"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/bear.png
var icon []byte

func main() {
	// Create an instance of the app structure
	w, h := getWindowSize()
	myApp := app.NewApp()

	err := wails.Run(&options.App{
		Title:              "xBrowser",
		Width:              w * 4 / 5,
		Height:             h * 4 / 5,
		MaxWidth:           w,
		MaxHeight:          h,
		MinWidth:           w / 2,
		MinHeight:          h / 2,
		Assets:             assets,
		LogLevelProduction: logger.DEBUG,
		Logger:             myApp.Logger,
		OnStartup:          myApp.Startup,
		OnDomReady:         myApp.DomReady,
		OnShutdown:         myApp.Shutdown,
		OnBeforeClose:      myApp.BeforeClose,
		LogLevel:           logger.DEBUG,
		Bind: []interface{}{
			myApp,
		},

		Windows: &windows.Options{},

		Mac: &mac.Options{
			//TitleBar:             mac.TitleBarHiddenInset(),
			//WebviewIsTransparent: true,
			//WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "xBrowser",
				Message: "A oss browser written by wails.",
				Icon:    icon,
			},
		},
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}
