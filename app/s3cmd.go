package app

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"sort"
	"strings"
	"time"
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
		res.Err = err.Error()
		return res
	}
	res.Buckets = buckets
	return res
}

func (a *App) MakeBucket(bucket string) string {
	err := a.S3Client.MakeBucket(bucket)
	if err != nil {
		return err.Error()
	}
	return ""
}

func (a *App) DeleteBucket(bucket string) string {
	err := a.S3Client.DeleteBucket(bucket)
	if err != nil {
		return err.Error()
	}
	return ""
}

type Object struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	HumanSize    string    `json:"human_size"`
	LastModified time.Time `json:"last_modified"`
}

type ListObjectResult struct {
	Contents []Object `json:"contents"`
	Prefixes []string `json:"prefixes"`
	Err      string   `json:"err"`
}

func (a *App) ListObjects(bucketName, marker, prefix string, maxKeys int64) ListObjectResult {
	out, err := a.S3Client.ListObjects(bucketName, marker, prefix, maxKeys)
	if err != nil {
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
	return res
}

type ObjectHandlerResult struct {
	Err string `json:"err"`
}

func getFileName(key string) string {
	sp := strings.Split(key, "/")
	return sp[len(sp)-1]
}

func (a *App) PutObject(bucketName, prefix, eventDialog, eventProgress string) ObjectHandlerResult {
	filePaths, err := runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	// TODO: Support mutifile upload
	if len(filePaths) == 0 {
		return ObjectHandlerResult{}
	}
	fp := filePaths[0]
	key := prefix + getFileName(fp)
	f, err := os.Open(fp)
	fInfo, err := f.Stat()
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	p := &Progress{}
	p.TotalBytes = fInfo.Size()
	r := NewProgressReader(f, p)

	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}

	// open dialog
	runtime.EventsEmit(a.ctx, eventDialog, true)
	cancelCh := make(chan bool)
	// emit progress data
	go func(ch chan bool) {
		for {
			select {
			case <-cancelCh:
				return
			default:
				if r.p.State() < 100 {
					runtime.EventsEmit(a.ctx, eventProgress, r.p.State())
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}(cancelCh)
	defer func() {
		cancelCh <- true
	}()

	err = a.S3Client.PutObjectWithOpt(bucketName, key, r)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	return ObjectHandlerResult{}
}

func (a *App) GetObject(bucketName, key string, override bool, eventDialog string, eventProgress string) ObjectHandlerResult {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	out, err := a.S3Client.GetObjectOutPut(bucketName, key)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}

	fName := getFileName(key)
	filePath := path + string(os.PathSeparator) + fName
	fileExist, err := PathExists(filePath)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	if fileExist {
		return ObjectHandlerResult{Err: "File already exists."}
	}
	f, err := os.Create(filePath)
	fmt.Println("GetObject:" + filePath)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}

	p := &Progress{}
	if out.ContentLength != nil {
		p.TotalBytes = *out.ContentLength
	}

	r := NewProgressReader(out.Body, p)

	// open dialog
	runtime.EventsEmit(a.ctx, eventDialog, true)

	cancelCh := make(chan bool)
	// emit progress data
	go func(ch chan bool) {
		for {
			select {
			case <-cancelCh:
				return
			default:
				runtime.EventsEmit(a.ctx, eventProgress, r.p.State())
				time.Sleep(10 * time.Millisecond)
			}
		}
	}(cancelCh)
	defer func() {
		cancelCh <- true
	}()
	_, err = io.Copy(f, r)
	if err != nil {
		delErr := os.Remove(filePath)
		return ObjectHandlerResult{Err: err.Error() + delErr.Error()}
	}
	runtime.EventsEmit(a.ctx, eventProgress, 100)
	return ObjectHandlerResult{}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (a *App) DeleteObject(bucketName, key string) ObjectHandlerResult {
	err := a.S3Client.DeleteObject(bucketName, key)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	return ObjectHandlerResult{}
}
