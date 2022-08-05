package db

import (
	"time"
)

type LoginInfo struct {
	AccountId string    `gorm:"primaryKey" json:"accountId"`
	Endpoint  string    `json:"endpoint"`
	AccessKey string    `json:"ak"`
	SecretKey string    `json:"sk"`
	Remark    string    `json:"remark"`
	Prepath   string    `json:"prepath"`
	LoginTime time.Time `json:"loginTime"`
}

const (
	PENDING int = iota
	PAUSE
	ERROR
	FINISH
)

// Unfinished Task
type UploadTask struct {
	AccountId     string           `json:"accountId"`
	TaskId        string           `json:"taskId"`
	Bucket        string           `json:"bucket"`
	Key           string           `json:"key"`
	Name          string           `json:"name"`
	Source        string           `json:"source"`
	Size          int64            `json:"size"`
	HumanSize     string           `json:"humanSize"`
	UploadedSize  int64            `json:"uploadedSize"`
	UploadId      string           `json:"uploadId"`
	IsMultipart   bool             `json:"isMultipart"`
	PartSize      int64            `json:"partSize"`
	Status        int              `json:"status"`
	CompletedPart []*CompletedPart `json:"completedPart" gorm:"-"`
	ModifiedTime  time.Time        `json:"modifiedTime"`
}

type CompletedPart struct {
	// Entity tag returned when the part was uploaded.
	ETag string `json:"etag"`

	// Part number that identifies the part. This is a positive integer between
	// 1 and 10,000.
	PartNumber int64 `json:"partNumber"`
}

type Settings struct {
	PartSizeMB             int `json:"partSize"`
	UploadPartsConcurrency int `json:"uploadPartsConcurrency"`
	UploadConcurrency      int `json:"uploadConcurrency"`
	DownloadConcurrency    int `json:"downloadConcurrency"`
}

type SettingsOption func(settings *Settings)

func NewSettings(opt ...SettingsOption) *Settings {
	s := NewDefaultSettings()
	for _, f := range opt {
		f(s)
	}
	return s
}

func NewDefaultSettings() *Settings {
	return &Settings{
		PartSizeMB:             5,
		UploadPartsConcurrency: 10,
		UploadConcurrency:      10,
		DownloadConcurrency:    10,
	}
}
