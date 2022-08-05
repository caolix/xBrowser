package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"xBrowser/app/db"
)

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
	p := NewProgress(a.ctx, u.TaskId, fInfo.Size())
	r := NewUploadProgressReader(f, p)
	resCh := make(chan ObjectHandlerResult)
	if u.IsMultipart {
		completePart, maxPartNum, err := a.getUploadedParts(&u)
		if err != nil {
			return ObjectHandlerResult{Err: err.Error()}
		}
		u.UploadedSize = maxPartNum * u.PartSize
		u.CompletedPart = completePart
		runtime.LogDebugf(a.ctx, "getMaxUploadedPartNumber: %d", maxPartNum)
		_, err = r.Seek(u.UploadedSize, io.SeekStart)
		if err != nil {
			return ObjectHandlerResult{Err: err.Error()}
		}
	}
	go a.doPut(u.Bucket, u.Key, r, &u, resCh)
	return <-resCh
}

func (a *App) CancelUploadTask(u db.UploadTask) ObjectHandlerResult {
	if u.IsMultipart {
		runtime.LogDebugf(a.ctx, "AbortMultiPartUpload: %s %s %s", u.Bucket, u.Key, u.UploadId)
		err := a.S3Client.AbortMultiPartUpload(u.Bucket, u.Key, u.UploadId)
		if err != nil {
			runtime.LogErrorf(a.ctx, "AbortMultiPartUpload err: %v", err)
			return ObjectHandlerResult{Err: err.Error()}
		}
	}
	runtime.LogDebugf(a.ctx, "CancelUploadTask to DeleteUploadTask: %s %s %s", u.Bucket, u.Key, u.UploadId)
	db.GlobalAppDB.DeleteUploadTask(u.AccountId, u.TaskId)
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
