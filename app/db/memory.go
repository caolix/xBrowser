package db

//
//import (
//	"errors"
//	"sync"
//	"time"
//)
//
//type AppMemory struct {
//	loginInfo       sync.Map
//	uploadTask      sync.Map
//	settings        sync.Map
//	latestLoginInfo *LoginInfo
//}
//
//func (a *AppMemory) Init(addr string) (err error) {
//	return nil
//}
//
//func (a *AppMemory) Close() {
//	return
//}
//
//func (a *AppMemory) UpsertLoginInfo(l *LoginInfo) (err error) {
//	key := l.AccessKey + l.SecretKey + l.Endpoint
//	val, ok := a.loginInfo.Load(key)
//	if !ok {
//		a.loginInfo.Store(key, l)
//		a.latestLoginInfo = l
//		return nil
//	}
//	v := val.(*LoginInfo)
//	v.LoginTime = time.Now().Local()
//	v.Remark = l.Remark
//	v.Prepath = l.Prepath
//	a.latestLoginInfo = v
//	return nil
//}
//
//func (a *AppMemory) GetLatestLoginInfo() (*LoginInfo, error) {
//	if a.latestLoginInfo != nil {
//		return a.latestLoginInfo, nil
//	}
//	return nil, errors.New("not found")
//}
//
//func (a *AppMemory) ListAllLoginInfo() ([]LoginInfo, error) {
//	infos := []LoginInfo{}
//	a.loginInfo.Range(func(key, value any) bool {
//		v := value.(*LoginInfo)
//		infos = append(infos, *v)
//		return true
//	})
//	return infos, nil
//}
//
//func (a *AppMemory) ListAllUploadTasks(accountId string) ([]UploadTask, error) {
//	tasks := []UploadTask{}
//	t, ok := a.uploadTask.Load(accountId)
//	if ok {
//		tt := t.()
//	}
//	for _, v := range a.uploadTask[accountId] {
//		tasks = append(tasks, *v)
//	}
//	return tasks, nil
//}
//
//func (a *AppMemory) CreateUploadTask(u *UploadTask) error {
//	v, ok := a.uploadTask.Load(u.AccountId)
//	if !ok {
//		tasks := make(map[string]*UploadTask)
//		tasks[u.TaskId] = u
//		a.uploadTask.Store(u.AccountId, tasks)
//		return nil
//	}
//	t, ok := v[u.TaskId]
//	if !ok {
//		a.uploadTask[u.AccountId][u.TaskId] = u
//		return nil
//	}
//
//	t.ModifiedTime = time.Now().Local()
//	t.Status = u.Status
//	t.UploadedSize = u.UploadedSize
//	return nil
//}
//
//func (a *AppMemory) DeleteUploadTask(accountId string, taskId string) {
//	v, ok := a.uploadTask[accountId]
//	if !ok {
//		return
//	}
//	delete(v, taskId)
//}
//
//func (a *AppMemory) UpdateSettings(settings *Settings) error {
//
//}
//
//func (a *AppMemory) LoadSettings() (*Settings, error) {
//
//}
