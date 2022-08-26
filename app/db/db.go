package db

type DB_TYPE string

var GlobalAppDB AppDB

const (
	TYPE_SQLITE = "sqlite"
)

const DB_NAME = "xbrowser_db"

type AppDB interface {
	Init(addr string) (err error)
	Close()
	UpsertLoginInfo(l *LoginInfo) (err error)
	GetLatestLoginInfo() (*LoginInfo, error)
	ListAllLoginInfo() ([]LoginInfo, error)

	ListAllUploadTasks(accountId string) ([]UploadTask, error)
	UpsertUploadTask(u *UploadTask) error
	DeleteUploadTask(accountId string, taskId string)

	ListAllDownloadTasks(accountId string) ([]DownloadTask, error)
	CreateDownloadTask(u *DownloadTask) error
	CreateDownloadPart(taskId string, partNum, partSize int64) error
	DeleteDownloadTask(accountId string, taskId string)
	GetAllDownloadParts(taskId string) ([]*CompletedDownloadPart, error)

	UpdateSettings(settings *Settings) error
	LoadSettings(accountId string) (*Settings, error)
}
