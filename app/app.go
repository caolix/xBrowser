package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"xBrowser/app/db"
	"xBrowser/app/s3lib"
)

// App struct
type App struct {
	ctx      context.Context
	S3Client *s3lib.S3Client
	DB       db.AppDB
	Config   *AppConfig

	// App setup status
	LoadDbErr error
}

type AppConfig struct {
	DbType  db.DB_TYPE
	Address string
	// TODO: add logger
}

// NewApp creates a new App application struct
func NewApp() *App {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	DefaultConfig := &AppConfig{
		DbType:  db.TYPE_SQLITE,
		Address: dir,
	}
	return &App{
		Config: DefaultConfig,
	}
}

// startup is called at application startup
// startup 在应用程序启动时调用
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置
	a.ctx = ctx
	// TODO: use OS_ENV to choose db config, load it and then set db type
	switch a.Config.DbType {
	case db.TYPE_SQLITE:
		a.DB = &db.AppSqlite{}
	default:
		a.LoadDbErr = errors.New("db type not supported")
	}
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) DomReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
	if a.DB != nil {
		a.LoadDbErr = a.DB.Init(a.Config.Address)
		if a.LoadDbErr != nil {
			a.DB = nil
		}
	}
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) Shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
	if a.DB != nil {
		a.DB.Close()
	}
}
