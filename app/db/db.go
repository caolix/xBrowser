package db

type DB_TYPE string

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
}
