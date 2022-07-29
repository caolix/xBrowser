package app

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"time"
	"xBrowser/app/db"
	"xBrowser/app/s3lib"
)

// Greet returns a greeting for the given name
func (a *App) Login(l db.LoginInfo, needSave bool) string {
	runtime.LogInfof(a.ctx, "Login info: %s %s %s %v", l.Endpoint, l.AccessKey, l.SecretKey, needSave)
	a.S3Client = s3lib.NewS3(l.Endpoint, l.AccessKey, l.SecretKey)
	_, err := a.S3Client.ListBuckets()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Login failed. err: %s", err)
		return err.Error()
	}
	runtime.LogInfof(a.ctx, "Login success.")
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
		if err != nil {
			runtime.LogWarningf(a.ctx, "Insert login info failed. err: %s", err)
		}
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
		if res != nil {
			return *res
		}
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
