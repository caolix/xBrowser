package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"os"
	"sort"
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

func (a *App) GetObject(bucketName, key string, override bool) ObjectHandlerResult {

	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	out, err := a.S3Client.GetObjectOutPut(bucketName, key)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	filePath := path + string(os.PathSeparator) + key
	fileExist, err := PathExists(filePath)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	if fileExist {
		return ObjectHandlerResult{Err: "File already exists."}
	}
	f, err := os.Create(filePath)
	if err != nil {
		return ObjectHandlerResult{Err: err.Error()}
	}
	_, err = io.Copy(f, out.Body)
	if err != nil {
		delErr := os.Remove(filePath)
		return ObjectHandlerResult{Err: err.Error() + delErr.Error()}
	}
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
