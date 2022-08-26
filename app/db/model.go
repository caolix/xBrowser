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
	UseSSL    bool      `json:"useSSL"`
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

// Unfinished Task
type DownloadTask struct {
	AccountId     string                   `json:"accountId"`
	TaskId        string                   `json:"taskId"`
	Bucket        string                   `json:"bucket"`
	Key           string                   `json:"key"`
	Name          string                   `json:"name"`
	Destination   string                   `json:"dest"`
	Size          int64                    `json:"size"`
	HumanSize     string                   `json:"humanSize"`
	CompletedPart []*CompletedDownloadPart `json:"completedPart" gorm:"-"`
	Status        int                      `json:"status"`
}

type CompletedDownloadPart struct {
	TaskId   string `json:"taskId"`
	Offset   int64  `json:"offset"`
	PartSize int64  `json:"partSize"`
}

type CompletedPart struct {
	// Entity tag returned when the part was uploaded.
	ETag string `json:"etag"`

	// Part number that identifies the part. This is a positive integer between
	// 1 and 10,000.
	PartNumber int64 `json:"partNumber"`
}

const (
	MaxUploadConcurrency      = 10
	MaxDownloadConcurrency    = 10
	MaxUploadPartsConcurrency = 10
	MaxPartSizeMB             = 5120
)

type Settings struct {
	AccountId              string `gorm:"primaryKey" json:"accountId"`
	PartSizeMB             int    `json:"partSize"`
	UploadPartsConcurrency int    `json:"uploadPartsConcurrency"`
	UploadConcurrency      int    `json:"uploadConcurrency"`
	DownloadConcurrency    int    `json:"downloadConcurrency"`
}

type SettingsOption func(settings *Settings)

func NewSettings(accountId string, opt ...SettingsOption) *Settings {
	s := NewDefaultSettings()
	for _, f := range opt {
		f(s)
	}
	s.AccountId = accountId
	return s
}

func NewDefaultSettings() *Settings {
	return &Settings{
		PartSizeMB:             5,
		UploadPartsConcurrency: 2,
		UploadConcurrency:      2,
		DownloadConcurrency:    1,
	}
}

func (s *Settings) SetDefault() {
	if s.PartSizeMB < 1 || s.PartSizeMB > MaxPartSizeMB {
		s.PartSizeMB = 5
	}
	if s.UploadConcurrency < 1 || s.UploadConcurrency > MaxUploadConcurrency {
		s.UploadConcurrency = 2
	}
	if s.UploadPartsConcurrency < 1 || s.UploadPartsConcurrency > MaxUploadPartsConcurrency {
		s.UploadPartsConcurrency = 2
	}
	if s.DownloadConcurrency < 1 || s.DownloadConcurrency > MaxDownloadConcurrency {
		s.DownloadConcurrency = 1
	}
}

func (s *Settings) Validate() (string, bool) {
	if s.UploadConcurrency < 1 || s.UploadConcurrency > MaxUploadConcurrency {
		return "UploadConcurrency", false
	}

	if s.UploadPartsConcurrency < 1 || s.UploadPartsConcurrency > MaxUploadPartsConcurrency {
		return "UploadPartsConcurrency", false
	}
	if s.DownloadConcurrency < 1 || s.DownloadConcurrency > MaxDownloadConcurrency {
		return "DownloadConcurrency", false
	}
	if s.PartSizeMB < 1 || s.PartSizeMB > MaxPartSizeMB {
		return "PartSizeMB", false
	}
	return "", true
}
