package db

import (
	"errors"
	"time"
)

type AppMemory struct {
	loginInfo       map[string]*LoginInfo
	uploadTask      map[int]map[string]*UploadTask
	latestLoginInfo *LoginInfo
}

func (a *AppMemory) Init(addr string) (err error) {
	a.loginInfo = make(map[string]*LoginInfo)
	a.uploadTask = make(map[int]map[string]*UploadTask)
	return nil
}

func (a *AppMemory) Close() {
	return
}

func (a *AppMemory) UpsertLoginInfo(l *LoginInfo) (err error) {
	key := l.AccessKey + l.SecretKey + l.Endpoint
	v, ok := a.loginInfo[key]
	if !ok {
		a.loginInfo[key] = l
		a.latestLoginInfo = l
		return nil
	}
	v.LoginTime = time.Now().Local()
	v.Remark = l.Remark
	v.Prepath = l.Prepath
	a.latestLoginInfo = v
	return nil
}

func (a *AppMemory) GetLatestLoginInfo() (*LoginInfo, error) {
	if a.latestLoginInfo != nil {
		return a.latestLoginInfo, nil
	}
	return nil, errors.New("not found")
}

func (a *AppMemory) ListAllLoginInfo() ([]LoginInfo, error) {
	infos := []LoginInfo{}
	for _, v := range a.loginInfo {
		infos = append(infos, *v)
	}
	return infos, nil
}

func (a *AppMemory) ListAllUploadTasks(accountId int) ([]UploadTask, error) {
	tasks := []UploadTask{}
	for _, v := range a.uploadTask[accountId] {
		tasks = append(tasks, *v)
	}
	return tasks, nil
}

func (a *AppMemory) UpsertUploadTask(u *UploadTask) error {
	v, ok := a.uploadTask[u.AccountId]
	if !ok {
		tasks := make(map[string]*UploadTask)
		tasks[u.TaskId] = u
		a.uploadTask[u.AccountId] = tasks
		return nil
	}
	t, ok := v[u.TaskId]
	if !ok {
		a.uploadTask[u.AccountId][u.TaskId] = u
		return nil
	}

	t.ModifiedTime = time.Now().Local()
	t.Status = u.Status
	t.UploadedSize = u.UploadedSize
	return nil
}
