package app

import (
	"bytes"
	"context"
	"errors"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	. "oBrowser/app/models"
	"oBrowser/app/util"
)

func (a *App) ListBuckets() ListBucketsResult {
	res := ListBucketsResult{}
	buckets, err := a.S3Client.ListBuckets()
	if err != nil {
		runtime.LogErrorf(a.ctx, "ListBuckets err: %s ", err)
		res.Err = err.Error()
		return res
	}
	for _, bucket := range buckets.Buckets{
		b := Buckets{
			Bucket: *bucket.Name,
			CreateTime: (*bucket.CreationDate).Local().Format("2006-01-02 15:04:05"),
			//CreateTime: (*bucket.CreationDate).String(),
		}
		res.Buckets = append(res.Buckets, b)
	}
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

func (a *App) PutBucketVersioning(bucket string, status string) string {
	err := a.S3Client.PutBucketVersioning(bucket, status, "")
	if err != nil {
		runtime.LogErrorf(a.ctx, "PutBucketVersioning %s status %s err: %s", bucket, status, err)
		return err.Error()
	}
	return ""
}

func (a *App) GetBucketDetailResult(bucket string) GetBucketVersioningResult {
	versioning, _, err := a.S3Client.GetBucketVersioning(bucket)
	if err != nil {
		runtime.LogErrorf(a.ctx, "GetBucketVersioning %s err: %s", bucket, err)
		return GetBucketVersioningResult{Err: err.Error()}
	}
	res := GetBucketVersioningResult{}
	res.Versioning = versioning
	return res
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

// Windows will return '\'
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

func (a *App) SelectUploadFolder(prefix string, bucketName string) ErrResult {
	root, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		runtime.LogErrorf(a.ctx, "OpenDirectoryDialog err: %s ", err)
		return ErrResult{Err: err.Error()}
	}

	if len(root) == 0 {
		return ErrResult{}
	}

	folderName := getUploadName(root)
	runtime.EventsEmit(a.ctx, AppEventBus, Event{
		Name: EventShowTasks,
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

	return ErrResult{}
}

func (a *App) SelectUploadFiles(prefix string, bucketName string) ErrResult {
	filePaths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ErrResult{Err: err.Error()}
	}
	if len(filePaths) == 0 {
		return ErrResult{}
	}

	runtime.EventsEmit(a.ctx, AppEventBus, Event{
		Name: EventShowTasks,
		Args: []string{"upload"},
	})

	go func() {
		for _, fp := range filePaths {
			key := prefix + getUploadName(fp)
			a.DoPutObject(bucketName, key, fp, GenerateEventProgressName("upload", key))
		}
	}()

	return ErrResult{}
}

func (a *App) doPut(ctx context.Context, wrapper *UploadTaskWrapper) {
	_, err := a.S3Client.UploadObject(ctx, wrapper.readSeeker, wrapper.task)
	if err != nil {
		wrapper.resCh <- err
		return
	}
	wrapper.resCh <- nil
}

func (a *App) DoPutObject(bucketName, key, filePath, eventProgress string) ErrResult {
	f, err := os.Open(filePath)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Open file %s err: %s ", filePath, err)
		return ErrResult{Err: err.Error()}
	}
	fInfo, err := f.Stat()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Stat file %s err: %s ", filePath, err)
		return ErrResult{Err: err.Error()}
	}
	// NOTE: App file on MacOS is a DIRECTORY !
	if fInfo.IsDir() {
		go filepath.Walk(filePath, func(fp string, info fs.FileInfo, err error) error {
			kp := strings.Replace(fp, string(os.PathSeparator), "/", -1)
			k := key + kp[len(filePath):]
			if info.IsDir() {
				res := a.PutDir(bucketName, k, "")
				if res.Err != "" {
					return errors.New(res.Err)
				}
				return nil
			}
			a.DoPutObject(bucketName, k, fp, GenerateEventProgressName("upload", key))
			return nil
		})
	} else {
		task := &UploadTask{
			AccountId:    a.AccountId,
			TaskId:       eventProgress,
			Bucket:       bucketName,
			Key:          key,
			Name:         getUploadName(key),
			Source:       filePath,
			Size:         fInfo.Size(),
			HumanSize:    util.IBytes(uint64(fInfo.Size())),
			Status:       WAITING,
			ModifiedTime: time.Now().Local(),
		}

		a.pushUploadTask(task, f)
		return ErrResult{}
	}

	return ErrResult{}
}

func (a *App) pushUploadTask(task *UploadTask, f io.ReadSeekCloser) {
	p := NewProgress(a.ctx, task.TaskId, task.Size)
	wrapper := &UploadTaskWrapper{
		task:       task,
		readSeeker: NewUploadProgressReader(f, p),
		resCh:      make(chan error),
	}
	a.uploadTaskWaitQ[a.AccountId].Push(wrapper)
}

func (a *App) PutDir(bucketName, dirName, prefix string) ErrResult {
	key := prefix + dirName + "/"
	r := bytes.NewReader([]byte(""))
	err := a.S3Client.PutObject(bucketName, key, r)
	if err != nil {
		runtime.LogErrorf(a.ctx, "PutDir %s in bucket %s err: %s ", key, bucketName, err)
		return ErrResult{Err: err.Error()}
	}
	return ErrResult{}
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

func (a *App) DoGetObject(bucketName, key, versionId, destPath string, size int64, eventProgress string, override bool) ErrResult {
	downloadFileName := getDownloadName(key) + ".download"
	downloadFilePath := destPath + string(os.PathSeparator) + downloadFileName
	filePath := destPath + string(os.PathSeparator) + getDownloadName(key)
	// TODO: implement override
	var f *os.File
	f, err := os.Create(downloadFilePath)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Create file %s err: %s ", downloadFilePath, err)
		return ErrResult{Err: err.Error()}
	}

	task := &DownloadTask{
		AccountId:   a.AccountId,
		TaskId:      eventProgress,
		Bucket:      bucketName,
		Key:         key,
		VersionId:   versionId,
		Name:        getDownloadName(key),
		Destination: filePath,
		Size:        size,
		HumanSize:   util.IBytes(uint64(size)),
		Status:      WAITING,
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
		return ErrResult{Err: err.Error()}
	}
	f.Close()
	err = os.Rename(downloadFilePath, filePath)
	if err != nil {
		return ErrResult{Err: err.Error()}
	}
	return ErrResult{}
}

func (a *App) doGet(ctx context.Context, wrapper *DownloadTaskWrapper) {
	_, err := a.S3Client.DownloadObject(ctx, wrapper.w, wrapper.task)
	if err != nil {
		wrapper.resCh <- err
		return
	}
	wrapper.resCh <- nil
}

func (a *App) DeleteObject(bucketName, key, versionId string, selectedType string,
	eventDeleteSuccess, eventDeleteCount string) ErrResult {
	task := &DeleteTask{
		a:                  a,
		bucketName:         bucketName,
		eventDeleteSuccess: eventDeleteSuccess,
		eventDeleteCount:   eventDeleteCount,
		delCh:              make(chan DeleteKey, 100),
		keys: []DeleteKey{
			{
				Key:       key,
				KeyType:   selectedType,
				VersionId: versionId,
			},
		},
		wg: &sync.WaitGroup{},
	}
	task.Start()
	task.wg.Wait()
	return ErrResult{}
}

func (a *App) DeleteObjects(bucketName string, keys []DeleteKey,
	eventDeleteSuccess, eventDeleteCount string) ErrResult {
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
	return ErrResult{}
}

func (a *App) ListObjectVersions(bucketName string, key string) GetObjectVersionsResult {
	out, err := a.S3Client.ListObjectVersions(bucketName, key)
	if err != nil {
		runtime.LogErrorf(a.ctx, "ListObjectVersions %s %s err: %s", bucketName, key, err)
	}
	res := GetObjectVersionsResult{}
	for _, object := range out.Versions {
		obj := ObjectVersioning{
			ETag:         *object.ETag,
			IsLatest:     *object.IsLatest,
			Key:          *object.Key,
			LastModified: *object.LastModified,
			Owner:        *object.Owner.DisplayName,
			Size:         *object.Size,
			StorageClass: *object.StorageClass,
			VersionId:    *object.VersionId,
		}
		res.Versions = append(res.Versions, obj)
	}
	return res
}
