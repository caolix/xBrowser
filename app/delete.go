package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"sync"
	"sync/atomic"
)

const (
	TypeFolder string = "Folder"
	TypeObject string = "Object"
)

type DeleteKey struct {
	Key     string `json:"key"`
	KeyType string `json:"keyType"`
}

type DeleteTask struct {
	a                                    *App
	bucketName                           string
	total                                int64
	success                              int64
	eventDeleteSuccess, eventDeleteCount string // for event listener
	delCh                                chan DeleteKey
	keys                                 []DeleteKey
	wg                                   *sync.WaitGroup
}

func (t *DeleteTask) doDelete() {
	for {
		select {
		case d := <-t.delCh:
			err := t.a.S3Client.DeleteObject(t.bucketName, d.Key)
			if err != nil {
				runtime.LogErrorf(t.a.ctx, "DeleteObject %s in bucket %s err: %s ", d.Key, t.bucketName, err)
				continue
			}
			atomic.AddInt64(&t.success, 1)
			runtime.EventsEmit(t.a.ctx, t.eventDeleteSuccess, t.success)
			t.wg.Done()
		}
	}
}

func (t *DeleteTask) deleteFolder(key *DeleteKey) error {
	marker := ""
	for {
		res, err := t.a.S3Client.ListObjects(t.bucketName, marker, key.Key, 1000, "")
		if err != nil {
			return err
		}
		atomic.AddInt64(&t.total, int64(len(res.Contents)))
		runtime.EventsEmit(t.a.ctx, t.eventDeleteCount, t.total)
		t.wg.Add(len(res.Contents))
		for _, c := range res.Contents {
			delKey := DeleteKey{
				Key:     *c.Key,
				KeyType: TypeFolder,
			}
			t.delCh <- delKey
		}
		if !*res.IsTruncated {
			break
		}
		marker = *res.NextMarker
	}
	return nil
}

func (t *DeleteTask) Start() {
	go t.doDelete()
	if t.a == nil {
		return
	}
	for _, k := range t.keys {
		if k.KeyType == TypeObject {
			t.wg.Add(1)
			atomic.AddInt64(&t.total, 1)
			runtime.EventsEmit(t.a.ctx, t.eventDeleteCount, t.total)
			t.delCh <- k
		} else if k.KeyType == TypeFolder {
			err := t.deleteFolder(&k)
			if err != nil {
				runtime.LogErrorf(t.a.ctx, "DeleteFolder %s in bucket %s err: %s ", k.Key, t.bucketName, err)
			}
		} else {
			continue
		}
	}
}
