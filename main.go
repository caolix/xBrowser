package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"os"
	"os/exec"
	os_runtime "runtime"
	"xBrowser/app"
	logger2 "xBrowser/app/logger"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/bear.png
var icon []byte

func SetupMenu() *menu.Menu {
	AppMenu := menu.NewMenu()
	if os_runtime.GOOS == "darwin" {
		AppMenu.Append(menu.AppMenu())
		AppMenu.Append(menu.EditMenu()) // on macos platform, we should append EditMenu to enable Cmd+C,Cmd+V,Cmd+Z... shortcut
	}
	FileMenu := AppMenu.AddSubmenu("Debug")
	FileMenu.AddText("OpenLogDir", keys.CmdOrCtrl("1"), func(_ *menu.CallbackData) {
		cmdName, args := openDirCmd(logger2.AppLogDir)
		cmd := exec.Command(cmdName, args...)
		fmt.Println(cmd.Run())
	})
	FileMenu.AddSeparator()
	FileMenu.AddText("OpenErrorLog", keys.CmdOrCtrl("2"), func(_ *menu.CallbackData) {
		cmdName, args := openTextFileCmd(logger2.AppLogDir + string(os.PathSeparator) + logger2.AppErrLogName)
		cmd := exec.Command(cmdName, args...)
		fmt.Println(cmd.String())
		fmt.Println(cmd.Run())
	})
	return AppMenu
}

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
		Menu:               SetupMenu(),
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
