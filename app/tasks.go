package app

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"xBrowser/app/db"
	"xBrowser/app/util"
)

const (
	WorkerPending = iota
	WorkerRunning
	WorkerStopping
	WorkerStopped
)

type UploadTaskWrapper struct {
	task              *db.UploadTask
	disabelUpdateList bool
	readSeeker        io.ReadSeeker
	resCh             chan error
}

type uploadWorker struct {
	num    int
	ctx    context.Context
	status int32
	taskCh chan *UploadTaskWrapper
	stopCh chan struct{}
}

var TaskCancelFunc map[string]context.CancelFunc
var lock sync.Mutex

func GenerateEventProgressName(eventType string, key string) string {
	rand := &util.Random{}
	return eventType + key + rand.String(32, util.Alphanumeric)
}

var EventListObjectsQ = *util.NewQueue()

func (a *App) ListenListObjectsEvent() {
	for {
		if EventListObjectsQ.Size() == 0 {
			time.Sleep(time.Second)
			continue
		}
		EventListObjectsQ.Pop()
		runtime.EventsEmit(a.ctx, EventBackend, Event{
			Type: TypeListObjectsEvent,
			Args: nil,
		})
	}
}

func (u *uploadWorker) start(a *App) {
	for {
		if u.status == WorkerStopping {
			runtime.LogDebugf(a.ctx, "worker %d stopped. status: %d", u.num, u.status)
			u.setStatus(WorkerStopped)
			break
		}
		u.setStatus(WorkerPending)
		select {
		case <-u.ctx.Done():
			runtime.LogDebugf(a.ctx, "recieve context cancel on worker %d", u.num)
			u.setStatus(WorkerStopping)
		case <-u.stopCh:
			runtime.LogDebugf(a.ctx, "recieve stopCh on worker %d", u.num)
			u.setStatus(WorkerStopping)
		case wrapper := <-u.taskCh:
			runtime.EventsEmit(a.ctx, EventBackend, Event{
				Type: TypeListenUploadTaskEvent,
				Args: []string{wrapper.task.TaskId},
			})
			u.setStatus(WorkerRunning)
			uploadCtx, cancel := context.WithCancel(u.ctx)
			lock.Lock()
			TaskCancelFunc[wrapper.task.TaskId] = cancel
			lock.Unlock()
			go a.doPut(uploadCtx, wrapper)
			uploadErr := <-wrapper.resCh
			if uploadErr != nil {
				runtime.LogErrorf(a.ctx, "upload task %s %s err: %s",
					wrapper.task.Bucket, wrapper.task.Key, uploadErr.Error())
				a.ErrLogger.Error(fmt.Sprintf("upload task %s %s err: %s",
					wrapper.task.Bucket, wrapper.task.Key, uploadErr.Error()))
			} else {
				if EventListObjectsQ.Size() < 2 {
					EventListObjectsQ.Push(struct{}{})
				}
			}
			runtime.EventsOff(a.ctx, wrapper.task.TaskId)
		}
	}
}

func (u *uploadWorker) setStatus(status int32) {
	atomic.StoreInt32(&u.status, status)
}

func (a *App) LoadAllUploadTasks() []db.UploadTask {
	tasks, err := db.GlobalAppDB.ListAllUploadTasks(a.AccountId)
	if err != nil {
		runtime.LogErrorf(a.ctx, "db.ListAllUploadTasks err: %s", err.Error())
		return nil
	}
	runtime.LogDebugf(a.ctx, "db.ListAllUploadTasks res: %v", tasks)
	return tasks
}

func (a *App) ResumeUploadTask(u db.UploadTask) ObjectHandlerResult {
	f, err := os.Open(u.Source)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Open file %s err: %s ", u.Source, err)
		return ObjectHandlerResult{Err: err.Error()}
	}
	fInfo, err := f.Stat()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Stat file %s err: %s ", u.Source, err)
		return ObjectHandlerResult{Err: err.Error()}
	}
	runtime.LogDebugf(a.ctx, "ResumeUploadTask source: %s bucket: %s, key: %s", u.Source, u.Bucket, u.Key)
	if u.IsMultipart {
		completePart, maxPartNum, err := a.getUploadedParts(&u)
		if err != nil {
			return ObjectHandlerResult{Err: err.Error()}
		}
		u.UploadedSize = maxPartNum * u.PartSize
		u.CompletedPart = completePart
		runtime.LogDebugf(a.ctx, "getMaxUploadedPartNumber: %d", maxPartNum)
		if err != nil {
			return ObjectHandlerResult{Err: err.Error()}
		}
	}
	p := NewProgress(a.ctx, u.TaskId, fInfo.Size())
	wrapper := &UploadTaskWrapper{
		task:              &u,
		disabelUpdateList: true,
		readSeeker:        NewUploadProgressReader(f, p),
		resCh:             make(chan error),
	}
	a.uploadTaskWaitQ[a.AccountId].Push(wrapper)
	return ObjectHandlerResult{}
}

func (a *App) RemoveUploadTask(u db.UploadTask) ObjectHandlerResult {
	defer func() {
		runtime.LogDebugf(a.ctx, "CancelUploadTask: %s %s %s", u.Bucket, u.Key, u.UploadId)
		runtime.LogDebugf(a.ctx, "DeleteUploadTask: %s %s", u.AccountId, u.TaskId)
		db.GlobalAppDB.DeleteUploadTask(u.AccountId, u.TaskId)
	}()
	if cancel, ok := TaskCancelFunc[u.TaskId]; ok {
		cancel()
		lock.Lock()
		delete(TaskCancelFunc, u.TaskId)
		lock.Unlock()
	}
	if u.IsMultipart && u.Status != db.FINISH {
		runtime.LogDebugf(a.ctx, "AbortMultiPartUpload: %s %s %s", u.Bucket, u.Key, u.UploadId)
		err := a.S3Client.AbortMultiPartUpload(u.Bucket, u.Key, u.UploadId)
		if err != nil {
			runtime.LogErrorf(a.ctx, "AbortMultiPartUpload err: %v", err)
			return ObjectHandlerResult{Err: err.Error()}
		}
	}
	return ObjectHandlerResult{}
}

func (a *App) getUploadedParts(u *db.UploadTask) ([]*db.CompletedPart, int64, error) {
	var partNumberMarker int64 = 0
	var completePart = []*db.CompletedPart{}
	for {
		out, err := a.S3Client.ListMultipartUploadParts(u.Bucket, u.Key, u.UploadId, partNumberMarker)
		if err != nil {
			return nil, 0, err
		}

		for _, p := range out.Parts {
			if *p.PartNumber == partNumberMarker+1 {
				partNumberMarker = *p.PartNumber
				completePart = append(completePart, &db.CompletedPart{ETag: *p.ETag, PartNumber: partNumberMarker})
				continue
			} else {
				return completePart, partNumberMarker, nil
			}
		}

		if !*out.IsTruncated {
			break
		}
	}

	return completePart, partNumberMarker, nil
}

type DownloadTaskWrapper struct {
	task      *db.DownloadTask
	w         io.WriterAt
	resCh     chan error
	requestCh chan error
}

type downloadWorker struct {
	num    int
	ctx    context.Context
	status int32
	taskCh chan *DownloadTaskWrapper
	stopCh chan struct{}
}

func (u *downloadWorker) start(a *App) {
	for {
		if u.status == WorkerStopping {
			runtime.LogDebugf(a.ctx, "download worker %d stopped. status: %d", u.num, u.status)
			u.setStatus(WorkerStopped)
			break
		}
		u.setStatus(WorkerPending)
		select {
		case <-u.ctx.Done():
			runtime.LogDebugf(a.ctx, "recieve context cancel on download worker %d", u.num)
			u.setStatus(WorkerStopping)
		case <-u.stopCh:
			runtime.LogDebugf(a.ctx, "recieve stopCh on download worker %d", u.num)
			u.setStatus(WorkerStopping)
		case wrapper := <-u.taskCh:
			u.setStatus(WorkerRunning)
			downloadCtx, cancel := context.WithCancel(u.ctx)
			lock.Lock()
			TaskCancelFunc[wrapper.task.TaskId] = cancel
			lock.Unlock()
			go a.doGet(downloadCtx, wrapper)
			uploadErr := <-wrapper.resCh
			if uploadErr != nil {
				runtime.LogErrorf(a.ctx, "download task %s %s err: %s",
					wrapper.task.Bucket, wrapper.task.Key, uploadErr.Error())
			}
			wrapper.requestCh <- uploadErr
		}
	}
}

func (u *downloadWorker) setStatus(status int32) {
	atomic.StoreInt32(&u.status, status)
}

func (a *App) LoadAllDownloadTasks() []db.DownloadTask {
	tasks, err := db.GlobalAppDB.ListAllDownloadTasks(a.AccountId)
	if err != nil {
		runtime.LogErrorf(a.ctx, "db.ListAllDownloadTasks err: %s", err.Error())
		return nil
	}
	runtime.LogDebugf(a.ctx, "db.ListAllDownloadTasks res: %v", tasks)
	return tasks
}

func (a *App) ResumeDownloadTask(u db.DownloadTask) ObjectHandlerResult {
	f, err := os.Create(u.Destination)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Open file %s err: %s ", u.Destination, err)
		return ObjectHandlerResult{Err: err.Error()}
	}
	runtime.LogDebugf(a.ctx, "ResumeDownloadTask source: %s bucket: %s, key: %s", u.Destination, u.Bucket, u.Key)
	p := NewProgress(a.ctx, u.TaskId, u.Size)
	wrapper := &DownloadTaskWrapper{
		task:      &u,
		w:         NewDownloadProgressWriterAt(f, p),
		requestCh: make(chan error),
		resCh:     make(chan error),
	}
	a.downloadTaskQ <- wrapper
	if err = <-wrapper.requestCh; err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	return ObjectHandlerResult{}
}

func (a *App) RemoveDownloadTask(u db.DownloadTask) ObjectHandlerResult {
	defer func() {
		runtime.LogDebugf(a.ctx, "CancelDownloadTask: %s %s %s", u.Bucket, u.Key, u.Destination)
		runtime.LogDebugf(a.ctx, "DeleteDownloadTask: %s %s", u.AccountId, u.TaskId)
		db.GlobalAppDB.DeleteDownloadTask(u.AccountId, u.TaskId)
	}()
	if cancel, ok := TaskCancelFunc[u.TaskId]; ok {
		cancel()
		lock.Lock()
		delete(TaskCancelFunc, u.TaskId)
		lock.Unlock()
	}
	return ObjectHandlerResult{}
}
