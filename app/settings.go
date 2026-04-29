package app

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"oBrowser/app/db"
	. "oBrowser/app/models"
)

func (a *App) reloadWorkers() {
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

	for i, w := range a.downloadWorkers {
		if i < a.Config.AppSettings.DownloadConcurrency {
			if w.status == WorkerStopping {
				w.setStatus(WorkerRunning)
			}
			if w.status == WorkerRunning {
				continue
			}
			if w.status == WorkerStopped {
				// restart worker
				runtime.LogDebugf(a.ctx, "restart download worker: %d", w.num)
				w.ctx, _ = context.WithCancel(a.downloadCtx)
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
}

func (a *App) UpdateSettings(s Settings) ErrResult {
	settings := &s
	settings.AccountId = a.AccountId
	if _, ok := settings.Validate(); !ok {
		settings.SetDefault()
	}
	runtime.LogDebugf(a.ctx, "update settings: %v", *settings)
	err := db.GlobalAppDB.UpdateSettings(settings)
	if err != nil {
		return ErrResult{Err: err.Error()}
	}
	a.Config.AppSettings = settings
	a.reloadWorkers()
	return ErrResult{}
}

func (a *App) LoadSettings() Settings {
	return *a.Config.AppSettings
}
