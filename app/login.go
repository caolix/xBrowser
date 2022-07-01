package app

import (
	"fmt"
	"time"
	"xBrowser/app/db"
	"xBrowser/app/s3lib"
)

// Greet returns a greeting for the given name
func (a *App) Login(l db.LoginInfo, needSave bool) string {
	a.S3Client = s3lib.NewS3(l.Endpoint, l.AccessKey, l.SecretKey)
	_, err := a.S3Client.ListBuckets()
	if err != nil {
		return err.Error()
	}
	if needSave && a.DB != nil {
		// Update if the database has the same record, otherwise insert
		err = a.DB.UpsertLoginInfo(&db.LoginInfo{
			Endpoint:  l.Endpoint,
			AccessKey: l.AccessKey,
			SecretKey: l.SecretKey,
			Remark:    l.Remark,
			Prepath:   l.Prepath,
			LoginTime: time.Now().Local(),
		})
		// TODO: handle error
		fmt.Println("Insert err:", err)
	}
	return ""
}

func (a *App) CheckDbError() string {
	if a.LoadDbErr != nil {
		return a.LoadDbErr.Error()
	}
	return ""
}

func (a *App) LoadLatestLoginInfo() db.LoginInfo {
	info := db.LoginInfo{}
	if a.DB != nil {
		res, err := a.DB.GetLatestLoginInfo()
		if err != nil {
			return info
		}
		return *res
	}
	return info
}

func (a *App) ListAllLoginInfo() []db.LoginInfo {
	infos := []db.LoginInfo{}
	if a.DB != nil {
		res, err := a.DB.ListAllLoginInfo()
		if err != nil {
			return infos
		}
		return res
	}
	return infos
}
