package db

import "time"

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
	ERROR
)

// Unfinished Task
type UploadTask struct {
	AccountId    string    `json:"accountId"`
	TaskId       string    `json:"taskId"`
	Key          string    `json:"key"`
	Source       string    `json:"source"`
	Size         int64     `json:"size"`
	UploadedSize int       `json:"uploadedSize"`
	UploadId     string    `json:"uploadId"`
	IsMultipart  bool      `json:"isMultipart"`
	PartSize     int64     `json:"partSize"`
	Status       int       `json:"status"`
	ModifiedTime time.Time `json:"modifiedTime"`
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
