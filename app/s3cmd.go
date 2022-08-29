package app

import (
	"bytes"
	"context"
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

func (a *App) SelectUploadFolder() SelectUploadFolderResult {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "OpenDirectoryDialog err: %s ", err)
		return SelectUploadFolderResult{Err: err.Error()}
	}
	return SelectUploadFolderResult{
		Path: path,
	}

}

func (a *App) DoUploadFolder(prefix, root, eventWalkPath string) {
	folderName := getUploadName(root)
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		path = strings.Replace(path, string(os.PathSeparator), "/", -1)
		key := prefix + folderName + path[len(root):]
		var t string
		if info.IsDir() {
			t = TypeFolder
		} else {
			t = TypeObject
		}
		var sf = SelectedUploadFile{
			AccountId:  a.AccountId,
			SourcePath: path,
			Name:       getUploadName(key),
			Type:       t,
		}
		sf.Key = key
		sf.Size = info.Size()
		sf.HumanSize = util.IBytes(uint64(info.Size()))
		runtime.EventsEmit(a.ctx, eventWalkPath, sf)
		return nil
	})
}

type SelectedUploadFile struct {
	Object
	AccountId  string `json:"accountId"`
	Type       string `json:"type"`
	SourcePath string `json:"source"`
	Name       string `json:"name"`
}

type SelectUploadFilesResult struct {
	SelectedFile []SelectedUploadFile `json:"files"`
	Err          string               `json:"err"`
}

func (a *App) SelectUploadFiles(prefix string) SelectUploadFilesResult {
	filePaths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return SelectUploadFilesResult{Err: err.Error()}
	}
	if len(filePaths) == 0 {
		return SelectUploadFilesResult{}
	}

	var res SelectUploadFilesResult
	for _, fp := range filePaths {
		f, err := os.Open(fp)
		if err != nil {
			runtime.LogErrorf(a.ctx, "Open file %s err: %s ", fp, err)
			return SelectUploadFilesResult{Err: err.Error()}
		}
		fInfo, err := f.Stat()
		if err != nil {
			runtime.LogErrorf(a.ctx, "Stat file %s err: %s ", fp, err)
			return SelectUploadFilesResult{Err: err.Error()}
		}

		fName := getUploadName(fp)
		var sf = SelectedUploadFile{
			AccountId:  a.AccountId,
			SourcePath: fp,
			Name:       fName,
		}
		sf.Key = prefix + getUploadName(fp)
		sf.Size = fInfo.Size()
		sf.HumanSize = util.IBytes(uint64(fInfo.Size()))
		res.SelectedFile = append(res.SelectedFile, sf)
	}
	return res
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
		Status:       db.PENDING,
		ModifiedTime: time.Now().Local(),
	}

	p := NewProgress(a.ctx, eventProgress, fInfo.Size())
	wrapper := &UploadTaskWrapper{
		task:       task,
		readSeeker: NewUploadProgressReader(f, p),
		requestCh:  make(chan error),
		resCh:      make(chan error),
	}

	a.uploadTaskQ <- wrapper
	if err = <-wrapper.requestCh; err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
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

	defer f.Close()

	task := &db.DownloadTask{
		AccountId:   a.AccountId,
		TaskId:      eventProgress,
		Bucket:      bucketName,
		Key:         key,
		Name:        getDownloadName(key),
		Destination: filePath,
		Size:        size,
		HumanSize:   util.IBytes(uint64(size)),
		Status:      db.PENDING,
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
		return ObjectHandlerResult{Err: err.Error()}
	}
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
