package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"time"
	"xBrowser/app"
	logger2 "xBrowser/app/logger"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/bear.png
var icon []byte

const AppName = "xBrowser"

func main() {
	// Create an instance of the app structure
	w, h := getWindowSize()
	app := app.NewApp()
	// Create application with options
	logName := AppName + "-" + time.Now().Local().Format("2006-01-02") + ".log"
	err := wails.Run(&options.App{
		Title:              "xBrowser",
		Width:              w * 4 / 5,
		Height:             h * 4 / 5,
		MaxWidth:           w,
		MaxHeight:          h,
		MinWidth:           w / 2,
		MinHeight:          h / 2,
		Assets:             assets,
		Logger:             logger2.NewFileLogger(logName),
		LogLevelProduction: logger.DEBUG,
		OnStartup:          app.Startup,
		OnDomReady:         app.DomReady,
		OnShutdown:         app.Shutdown,
		OnBeforeClose:      app.BeforeClose,
		LogLevel:           logger.DEBUG,
		Bind: []interface{}{
			app,
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
