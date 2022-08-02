package db

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

type AppSqlite struct {
	Address string
	DB      *gorm.DB
}

func (s *AppSqlite) Init(dir string) (err error) {
	dbAddress := dir + string(os.PathSeparator) + DB_NAME
	s.DB, err = gorm.Open(sqlite.Open(dbAddress), &gorm.Config{})
	if err != nil {
		return err
	}
	s.Address = dbAddress
	err = s.DB.AutoMigrate(&LoginInfo{}, &UploadTask{}, &Settings{})
	if err != nil {
		s.DB = nil
		return err
	}
	return nil
}

func (s *AppSqlite) Close() {
	return
}

func (s *AppSqlite) UpsertLoginInfo(l *LoginInfo) (err error) {
	query := LoginInfo{
		Endpoint:  l.Endpoint,
		AccessKey: l.AccessKey,
		SecretKey: l.SecretKey,
	}
	info := LoginInfo{}
	res := s.DB.Where(&query).First(&info)
	if res.Error != nil {
		// Insert login info if not exist
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			res = s.DB.Create(l)
			if res.Error != nil {
				return res.Error
			}
			return nil
		}
		return res.Error
	}
	// Update login time if login info exist
	info.LoginTime = time.Now().Local()
	info.Remark = l.Remark
	info.Prepath = l.Prepath
	res = s.DB.Where("access_key = ? AND secret_key = ? AND endpoint = ?",
		l.AccessKey, l.SecretKey, l.Endpoint).Save(&info)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *AppSqlite) GetLatestLoginInfo() (*LoginInfo, error) {
	query := &LoginInfo{}
	res := s.DB.Order("login_time desc").First(&query)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return query, nil
}

func (s *AppSqlite) ListAllLoginInfo() ([]LoginInfo, error) {
	querys := []LoginInfo{}
	res := s.DB.Find(&querys)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return querys, nil
}

func (s *AppSqlite) ListAllUploadTasks(accountId int) ([]UploadTask, error) {
	querys := []UploadTask{}
	res := s.DB.Where("account_id = ?", accountId).Find(&querys)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, res.Error
	}
	return querys, nil
}

func (s *AppSqlite) UpsertUploadTask(u *UploadTask) error {
	query := UploadTask{
		AccountId: u.AccountId,
	}
	task := UploadTask{}
	res := s.DB.Where(&query).First(&task)
	if res.Error != nil {
		// Insert login info if not exist
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			res = s.DB.Create(u)
			if res.Error != nil {
				return res.Error
			}
			return nil
		}
		return res.Error
	}
	// Update login time if login info exist
	task.ModifiedTime = time.Now().Local()
	task.Status = u.Status
	task.UploadedSize = u.UploadedSize

	res = s.DB.Where("status = ? AND uploaded_size = ? AND modified_time = ?",
		u.Status, u.UploadedSize, u.ModifiedTime).Save(&task)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
