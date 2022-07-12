package app

import (
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
	return res
}
