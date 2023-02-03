package models

import "time"

type ListBucketsResult struct {
	Buckets []string `json:"buckets"`
	Err     string   `json:"err"`
}

type GetBucketVersioningResult struct {
	Bucket     string `json:"bucket"`
	Versioning string `json:"versioning"`
	Err        string `json:"err"`
}

type Object struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	HumanSize    string    `json:"humanSize"`
	LastModified time.Time `json:"lastModified"`
}

type ObjectVersioning struct {
	ETag         string    `json:"ETag"`
	IsLatest     bool      `json:"isLatest"`
	Key          string    `json:"key"`
	LastModified time.Time `json:"lastModified"`
	Owner        string    `json:"owner"`
	Size         int64     `json:"size"`
	StorageClass string    `json:"storageClass"`
	VersionId    string    `json:"versionId"`
}

type ListObjectResult struct {
	Contents    []Object `json:"contents"`
	Prefixes    []string `json:"prefixes"`
	NextMarker  string   `json:"nextMarker"`
	IsTruncated bool     `json:"isTruncated"`
	Err         string   `json:"err"`
}

type ErrResult struct {
	Err string `json:"err"`
}

type SelectedUploadFile struct {
	Object
	AccountId  string `json:"accountId"`
	Type       string `json:"type"`
	SourcePath string `json:"source"`
	Name       string `json:"name"`
}

type SelectDownloadPathResult struct {
	AccountId string `json:"accountId"`
	Path      string `json:"path"`
	Err       string `json:"err"`
}

type SelectedObject struct {
	Type string
	Key  string
}

type GetObjectVersionsResult struct {
	Versions []ObjectVersioning `json:"versions"`
	Err      string             `json:"err"`
}
