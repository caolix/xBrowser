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
	a.AccountId = util.GenAccountId(l.AccessKey, l.Endpoint)
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
	a.registerWorkers()
	a.loadTaskQueues()
	runtime.LogDebugf(a.ctx, "Login account id: %s", a.AccountId)
	return ""
}

func (a *App) loadTaskQueues() {
	a.loadTaskRunnningQ()
	a.loadTaskWaitQ()
}

func (a *App) loadTaskWaitQ() {
	// TODO: Load from db
	if _, ok := a.uploadTaskWaitQ[a.AccountId]; !ok {
		a.uploadTaskWaitQ[a.AccountId] = util.NewQueue()
		go a.listenTaskWaitQ()
	}
}

func (a *App) listenTaskWaitQ() {
	for {
		if a.uploadTaskWaitQ[a.AccountId].Size() == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// control running queue size
		if a.uploadTaskRunningQ[a.AccountId].Size() > 100 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		w := a.uploadTaskWaitQ[a.AccountId].Pop()
		wrapper := w.(*UploadTaskWrapper)
		if !wrapper.disabelUpdateList {
			runtime.EventsEmit(a.ctx, EventAddToUploadList, *wrapper.task)
		}
		a.uploadTaskRunningQ[a.AccountId].Push(w)
	}
}

func (a *App) loadTaskRunnningQ() {
	if _, ok := a.uploadTaskRunningQ[a.AccountId]; !ok {
		a.uploadTaskRunningQ[a.AccountId] = util.NewQueue()
		go a.listenTaskRunningQ()
	}
}

func (a *App) listenTaskRunningQ() {
	for {
		if a.uploadTaskRunningQ[a.AccountId].Size() == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		w := a.uploadTaskRunningQ[a.AccountId].Pop()
		wrapper := w.(*UploadTaskWrapper)
		a.uploadTaskCh <- wrapper
	}
}

func (a *App) GetAccountId() string {
	return a.AccountId
}

func (a *App) registerWorkers() {
	TaskCancelFunc = make(map[string]context.CancelFunc)
	a.uploadCtx, a.uploadCancelFunc = context.WithCancel(context.Background())
	go a.ListenListObjectsEvent()
	for i := 0; i < db.MaxUploadConcurrency; i++ {
		var ctx context.Context
		ctx, _ = context.WithCancel(a.uploadCtx)
		worker := &uploadWorker{
			num:    i,
			ctx:    ctx,
			taskCh: a.uploadTaskCh,
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

	a.downloadCtx, a.downloadCancelFunc = context.WithCancel(context.Background())
	for i := 0; i < db.MaxDownloadConcurrency; i++ {
		var ctx context.Context
		ctx, _ = context.WithCancel(a.downloadCtx)
		worker := &downloadWorker{
			num:    i,
			ctx:    ctx,
			taskCh: a.downloadTaskQ,
			stopCh: make(chan struct{}, 1),
			status: WorkerStopped,
		}
		if i < a.Config.AppSettings.UploadConcurrency {
			worker.setStatus(WorkerRunning)
			go worker.start(a)
		}
		runtime.LogDebugf(a.ctx, "registerWorkers id: %d status: %d", worker.num, worker.status)
		a.downloadWorkers = append(a.downloadWorkers, worker)
	}
}

func (a *App) Logout() {
	a.uploadCancelFunc()
	a.downloadCancelFunc()
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
