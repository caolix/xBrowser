package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path/filepath"
	"oBrowser/app/db"
	logger2 "oBrowser/app/logger"
	. "oBrowser/app/models"
	"oBrowser/app/s3lib"
	"oBrowser/app/util"
)

// App struct
type App struct {
	ctx      context.Context
	S3Client *s3lib.S3Client
	Config   *AppConfig

	uploadTaskCh       chan *UploadTaskWrapper
	uploadTaskWaitQ    map[string]*util.Queue // accountId -> TaskQ
	uploadTaskRunningQ map[string]*util.Queue // accountId -> TaskQ
	uploadWorkers      []*uploadWorker
	uploadCtx          context.Context
	uploadCancelFunc   context.CancelFunc

	downloadTaskQ      chan *DownloadTaskWrapper
	downloadWorkers    []*downloadWorker
	downloadCtx        context.Context
	downloadCancelFunc context.CancelFunc

	Logger    logger.Logger
	ErrLogger logger.Logger

	// App setup status
	InitLoggerErr error
	LoadDbErr     error
	AccountId     string
}

type AppConfig struct {
	DbType      db.DB_TYPE
	Address     string
	AppSettings *Settings
}

// NewApp creates a new App application struct
func NewApp() *App {
	app := &App{
		uploadTaskCh:       make(chan *UploadTaskWrapper),
		uploadTaskWaitQ:    make(map[string]*util.Queue),
		uploadTaskRunningQ: make(map[string]*util.Queue),

		downloadTaskQ: make(chan *DownloadTaskWrapper),
	}
	dbDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	app.Config = &AppConfig{
		DbType:  db.TYPE_SQLITE,
		Address: dbDir,
	}

	app.Logger, app.InitLoggerErr = logger2.NewAppLogger()
	if app.InitLoggerErr != nil {
		fmt.Println(app.InitLoggerErr)
	}
	app.ErrLogger, app.InitLoggerErr = logger2.NewAppErrLogger()
	if app.InitLoggerErr != nil {
		fmt.Println(app.InitLoggerErr)
	}
	return app
}

// startup is called at application startup
// startup 在应用程序启动时调用
// 1. Init db, and then load account settings.
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置
	defer runtime.LogInfo(ctx, "Startup finished.")
	a.ctx = ctx
	// TODO: use OS_ENV to choose db config, load it and then set db type
	runtime.LogInfo(ctx, "Load db type:"+string(a.Config.DbType))
	switch a.Config.DbType {
	case db.TYPE_SQLITE:
		db.GlobalAppDB = &db.AppSqlite{Logger: a.Logger}
	default:
		runtime.LogError(ctx, "not supported:"+string(a.Config.DbType))
		a.LoadDbErr = errors.New("db type not supported")
		//db.GlobalAppDB = &db.AppMemory{}
		return
	}
	err := db.GlobalAppDB.Init(a.Config.Address)
	if err != nil {
		a.LoadDbErr = err
		return
	}
	settings, err := db.GlobalAppDB.LoadSettings(a.AccountId)
	if err != nil {
		a.LoadDbErr = err
		return
	}
	a.Config.AppSettings = settings
	runtime.LogDebugf(ctx, "LoadSettings: %v", *settings)
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
	defer runtime.LogInfo(ctx, "DomReady finished.")
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
	if sqlite, ok := db.GlobalAppDB.(*db.AppSqlite); ok && sqlite.DB != nil {
		sqlite.DB.Close()
	}
	runtime.LogInfo(ctx, "Shutdown finished.")
}

func (a *App) BeforeClose(ctx context.Context) bool {
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Quit?",
		Message:       "Are you sure you want to quit?",
		Buttons:       []string{"No", "Yes"},
		DefaultButton: "Yes",
	})
	if err != nil {
		return false
	}
	return dialog != "Yes"
}
