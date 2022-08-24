package app

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
	"xBrowser/app/db"
	"xBrowser/app/s3lib"
	"xBrowser/app/util"
)

// Greet returns a greeting for the given name
func (a *App) Login(l db.LoginInfo, needSave bool, isHttps bool) string {
	runtime.LogDebugf(a.ctx, "Login info: %s %s %s %v %v", l.Endpoint, l.AccessKey, l.SecretKey, needSave, isHttps)
	a.S3Client = s3lib.NewS3(l.Endpoint, l.AccessKey, l.SecretKey, isHttps)
	_, err := a.S3Client.ListBuckets()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Login failed. err: %s", err)
		return err.Error()
	}
	runtime.LogInfof(a.ctx, "Login success.")
	if needSave {
		// Update if the database has the same record, otherwise insert
		err = db.GlobalAppDB.UpsertLoginInfo(&db.LoginInfo{
			Endpoint:  l.Endpoint,
			AccessKey: l.AccessKey,
			SecretKey: l.SecretKey,
			Remark:    l.Remark,
			Prepath:   l.Prepath,
			LoginTime: time.Now().Local(),
			UseSSL:    isHttps,
		})
		if err != nil {
			runtime.LogWarningf(a.ctx, "Insert login info failed. err: %s", err)
		}
	}
	a.AccountId = util.GenAccountId(l.AccessKey, l.Endpoint)
	a.registerWorkers()
	runtime.LogDebugf(a.ctx, "Login account id: %s", a.AccountId)
	return ""
}

func (a *App) registerWorkers() {
	UploadTaskCancelFunc = make(map[string]context.CancelFunc)
	a.uploadCtx, a.cancelFunc = context.WithCancel(context.Background())
	for i := 0; i < db.MaxUploadConcurrency; i++ {
		var ctx context.Context
		ctx, _ = context.WithCancel(a.uploadCtx)
		worker := &uploadWorker{
			num:    i,
			ctx:    ctx,
			taskCh: a.uploadTaskQ,
			stopCh: make(chan struct{}, 1),
			status: WorkerStopped,
		}
		if i < a.Config.AppSettings.UploadConcurrency {
			worker.setStatus(WorkerRunning)
			go worker.start(a)
		}
		runtime.LogDebugf(a.ctx, "registerWorkers id: %d status: %d", worker.num, worker.status)
		a.uploadWorkers = append(a.uploadWorkers, worker)
	}
}

func (a *App) Logout() {
	a.cancelFunc()
}

func (a *App) CheckDbError() string {
	if a.LoadDbErr != nil {
		return a.LoadDbErr.Error()
	}
	return ""
}

func (a *App) LoadLatestLoginInfo() db.LoginInfo {
	info := db.LoginInfo{}
	res, err := db.GlobalAppDB.GetLatestLoginInfo()
	if err != nil {
		return info
	}
	if res != nil {
		return *res
	}
	return info
}

func (a *App) ListAllLoginInfo() []db.LoginInfo {
	infos := []db.LoginInfo{}
	res, err := db.GlobalAppDB.ListAllLoginInfo()
	if err != nil {
		return infos
	}
	return res
}
