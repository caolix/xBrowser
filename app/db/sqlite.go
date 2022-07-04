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
	err = s.DB.AutoMigrate(&LoginInfo{})
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
	query := &LoginInfo{
		Endpoint:  l.Endpoint,
		AccessKey: l.AccessKey,
		SecretKey: l.SecretKey,
	}
	res := s.DB.First(&query)
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
	query.LoginTime = time.Now().Local()
	query.Remark = l.Remark
	query.Prepath = l.Prepath
	res = s.DB.Where("access_key = ? AND secret_key = ? AND endpoint = ?",
		query.AccessKey, query.SecretKey, query.Endpoint).Save(&query)
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
