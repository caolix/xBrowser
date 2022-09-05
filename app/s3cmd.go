package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"xBrowser/app/db"
	"xBrowser/app/util"
)

type ListBucketResult struct {
	Buckets []string `json:"buckets"`
	Err     string   `json:"err"`
}

func (a *App) ListBuckets() ListBucketResult {
	res := ListBucketResult{}
	buckets, err := a.S3Client.ListBuckets()
	if err != nil {
		runtime.LogErrorf(a.ctx, "ListBuckets err: %s ", err)
		res.Err = err.Error()
		return res
	}
	res.Buckets = buckets
	runtime.LogDebug(a.ctx, "ListBuckets success.")
	return res
}

func (a *App) MakeBucket(bucket string) string {
	err := a.S3Client.MakeBucket(bucket)
	if err != nil {
		runtime.LogErrorf(a.ctx, "MakeBucket %s err: %s ", bucket, err)
		return err.Error()
	}
	return ""
}

func (a *App) DeleteBucket(bucket string) string {
	err := a.S3Client.DeleteBucket(bucket)
	if err != nil {
		runtime.LogErrorf(a.ctx, "DeleteBucket %s err: %s ", bucket, err)
		return err.Error()
	}
	return ""
}

type Object struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	HumanSize    string    `json:"humanSize"`
	LastModified time.Time `json:"lastModified"`
}

type ListObjectResult struct {
	Contents    []Object `json:"contents"`
	Prefixes    []string `json:"prefixes"`
	NextMarker  string   `json:"nextMarker"`
	IsTruncated bool     `json:"isTruncated"`
	Err         string   `json:"err"`
}

func (a *App) ListObjects(bucketName, marker, prefix string, maxKeys int64) ListObjectResult {
	out, err := a.S3Client.ListObjects(bucketName, marker, prefix, maxKeys, "/")
	if err != nil {
		runtime.LogErrorf(a.ctx, "ListObjects %s %s %s %d %s err: %s ", bucketName, marker, prefix, maxKeys, "/", err)
		return ListObjectResult{Err: err.Error()}
	}
	res := ListObjectResult{}

	for _, content := range out.Contents {
		o := Object{
			Key:          *content.Key,
			Size:         *content.Size,
			HumanSize:    util.IBytes(uint64(*content.Size)),
			LastModified: *content.LastModified,
		}
		res.Contents = append(res.Contents, o)
	}
	for _, prefix := range out.CommonPrefixes {
		res.Prefixes = append(res.Prefixes, *prefix.Prefix)
	}
	sort.Slice(res.Prefixes, func(i, j int) bool {
		return res.Prefixes[i] < res.Prefixes[j]
	})
	if out.NextMarker != nil {
		res.NextMarker = *out.NextMarker
	}
	if out.IsTruncated != nil {
		res.IsTruncated = *out.IsTruncated
	}
	return res
}

type ObjectHandlerResult struct {
	Err string `json:"err"`
}

func getUploadName(key string) string {
	sp := strings.Split(key, string(os.PathSeparator))
	return sp[len(sp)-1]
}

func getDownloadName(key string) string {
	sp := strings.Split(key, "/")
	return sp[len(sp)-1]
}

type SelectUploadFolderResult struct {
	Path string `json:"path"`
	Err  string `json:"err"`
}

func (a *App) SelectUploadFolder(prefix string, bucketName string) ObjectHandlerResult {
	root, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "OpenDirectoryDialog err: %s ", err)
		return ObjectHandlerResult{Err: err.Error()}
	}

	if len(root) == 0 {
		return ObjectHandlerResult{}
	}

	folderName := getUploadName(root)
	runtime.EventsEmit(a.ctx, EventBackend, Event{
		Type: TypeShowTasksEvent,
		Args: []string{"upload"},
	})

	go filepath.Walk(root, func(fp string, info fs.FileInfo, err error) error {
		kp := strings.Replace(fp, string(os.PathSeparator), "/", -1)
		key := prefix + folderName + kp[len(root):]
		if info.IsDir() {
			res := a.PutDir(bucketName, folderName, prefix)
			if res.Err != "" {
				return errors.New(res.Err)
			}
			return nil
		}

		a.DoPutObject(bucketName, key, fp, GenerateEventProgressName("upload", key))
		return nil
	})

	fmt.Println("SelectUploadFolder FINISH")
	return ObjectHandlerResult{}
}

type SelectedUploadFile struct {
	Object
	AccountId  string `json:"accountId"`
	Type       string `json:"type"`
	SourcePath string `json:"source"`
	Name       string `json:"name"`
}

func (a *App) SelectUploadFiles(prefix string, bucketName string) ObjectHandlerResult {
	filePaths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	if len(filePaths) == 0 {
		return ObjectHandlerResult{}
	}

	runtime.EventsEmit(a.ctx, EventBackend, Event{
		Type: TypeShowTasksEvent,
		Args: []string{"upload"},
	})
	go func() {
		for _, fp := range filePaths {
			key := prefix + getUploadName(fp)
			a.DoPutObject(bucketName, key, fp, GenerateEventProgressName("upload", key))
		}
	}()

	fmt.Println("SelectUploadFiles FINISH")
	return ObjectHandlerResult{}
}

func (a *App) doPut(ctx context.Context, wrapper *UploadTaskWrapper) {
	_, err := a.S3Client.UploadObject(ctx, wrapper.readSeeker, wrapper.task)
	if err != nil {
		wrapper.resCh <- err
		return
	}
	wrapper.resCh <- nil
}

func (a *App) DoPutObject(bucketName, key, filePath, eventProgress string) ObjectHandlerResult {
	f, err := os.Open(filePath)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Open file %s err: %s ", filePath, err)
		return ObjectHandlerResult{Err: err.Error()}
	}
	fInfo, err := f.Stat()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Stat file %s err: %s ", filePath, err)
		return ObjectHandlerResult{Err: err.Error()}
	}

	task := &db.UploadTask{
		AccountId:    a.AccountId,
		TaskId:       eventProgress,
		Bucket:       bucketName,
		Key:          key,
		Name:         getUploadName(key),
		Source:       filePath,
		Size:         fInfo.Size(),
		HumanSize:    util.IBytes(uint64(fInfo.Size())),
		Status:       db.PAUSE,
		ModifiedTime: time.Now().Local(),
	}

	p := NewProgress(a.ctx, eventProgress, fInfo.Size())
	wrapper := &UploadTaskWrapper{
		task:       task,
		readSeeker: NewUploadProgressReader(f, p),
		resCh:      make(chan error),
	}

	a.uploadTaskWaitQ[a.AccountId].Push(wrapper)
	return ObjectHandlerResult{}
}

func (a *App) PutDir(bucketName, dirName, prefix string) ObjectHandlerResult {
	key := prefix + dirName + "/"
	r := bytes.NewReader([]byte(""))
	err := a.S3Client.PutObject(bucketName, key, r)
	if err != nil {
		runtime.LogErrorf(a.ctx, "PutDir %s in bucket %s err: %s ", key, bucketName, err)
		return ObjectHandlerResult{Err: err.Error()}
	}
	return ObjectHandlerResult{}
}

type SelectDownloadPathResult struct {
	AccountId string `json:"accountId"`
	Path      string `json:"path"`
	Err       string `json:"err"`
}

func (a *App) SelectDownloadPath() SelectDownloadPathResult {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "OpenDirectoryDialog err: %s ", err)
		return SelectDownloadPathResult{Err: err.Error()}
	}
	return SelectDownloadPathResult{
		AccountId: a.AccountId,
		Path:      path,
	}
}

func (a *App) DoGetObject(bucketName, key, destPath string, size int64, eventProgress string, override bool) ObjectHandlerResult {
	downloadFileName := getDownloadName(key) + ".download"
	downloadFilePath := destPath + string(os.PathSeparator) + downloadFileName
	filePath := destPath + string(os.PathSeparator) + getDownloadName(key)
	// TODO: implement override
	var f *os.File
	f, err := os.Create(downloadFilePath)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Create file %s err: %s ", downloadFilePath, err)
		return ObjectHandlerResult{Err: err.Error()}
	}

	task := &db.DownloadTask{
		AccountId:   a.AccountId,
		TaskId:      eventProgress,
		Bucket:      bucketName,
		Key:         key,
		Name:        getDownloadName(key),
		Destination: filePath,
		Size:        size,
		HumanSize:   util.IBytes(uint64(size)),
		Status:      db.PAUSE,
	}

	p := NewProgress(a.ctx, eventProgress, size)
	wrapper := &DownloadTaskWrapper{
		task:      task,
		w:         NewDownloadProgressWriterAt(f, p),
		resCh:     make(chan error),
		requestCh: make(chan error),
	}

	a.downloadTaskQ <- wrapper
	if err = <-wrapper.requestCh; err != nil {
		f.Close()
		return ObjectHandlerResult{Err: err.Error()}
	}
	f.Close()
	err = os.Rename(downloadFilePath, filePath)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	return ObjectHandlerResult{}
}

func (a *App) doGet(ctx context.Context, wrapper *DownloadTaskWrapper) {
	_, err := a.S3Client.DownloadObject(ctx, wrapper.w, wrapper.task)
	if err != nil {
		wrapper.resCh <- err
		return
	}
	wrapper.resCh <- nil
}

func (a *App) DeleteObject(bucketName, key string, selectedType string,
	eventDeleteSuccess, eventDeleteCount string) ObjectHandlerResult {
	task := &DeleteTask{
		a:                  a,
		bucketName:         bucketName,
		eventDeleteSuccess: eventDeleteSuccess,
		eventDeleteCount:   eventDeleteCount,
		delCh:              make(chan DeleteKey, 100),
		keys: []DeleteKey{
			{
				Key:     key,
				KeyType: selectedType,
			},
		},
		wg: &sync.WaitGroup{},
	}
	task.Start()
	task.wg.Wait()
	return ObjectHandlerResult{}
}

type SelectedObject struct {
	Type string
	Key  string
}

func (a *App) DeleteObjects(bucketName string, keys []DeleteKey,
	eventDeleteSuccess, eventDeleteCount string) ObjectHandlerResult {
	task := &DeleteTask{
		a:                  a,
		bucketName:         bucketName,
		eventDeleteSuccess: eventDeleteSuccess,
		eventDeleteCount:   eventDeleteCount,
		delCh:              make(chan DeleteKey, 100),
		keys:               keys,
		wg:                 &sync.WaitGroup{},
	}

	task.Start()
	task.wg.Wait()
	return ObjectHandlerResult{}
}
