package app

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"xBrowser/app/db"
)

func (a *App) UpdateSettings(s db.Settings) ObjectHandlerResult {
	settings := &s
	settings.AccountId = a.AccountId
	if _, ok := settings.Validate(); !ok {
		settings.SetDefault()
	}
	runtime.LogDebugf(a.ctx, "update settings: %v", *settings)
	err := db.GlobalAppDB.UpdateSettings(settings)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	a.Config.AppSettings = settings
	for i, w := range a.uploadWorkers {
		if i < a.Config.AppSettings.UploadConcurrency {
			if w.status == WorkerStopping {
				w.setStatus(WorkerRunning)
			}
			if w.status == WorkerRunning {
				continue
			}
			if w.status == WorkerStopped {
				// restart worker
				runtime.LogDebugf(a.ctx, "restart worker: %d", w.num)
				w.ctx, _ = context.WithCancel(a.uploadCtx)
				w.setStatus(WorkerRunning)
				go w.start(a)
			}
		} else {
			if w.status == WorkerPending {
				w.stopCh <- struct{}{}
				continue
			}
			if w.status == WorkerRunning {
				w.setStatus(WorkerStopping)
			}

		}
	}
	return ObjectHandlerResult{}
}

func (a *App) LoadSettings() db.Settings {
	return *a.Config.AppSettings
}
