package db

import "time"

type LoginInfo struct {
	Endpoint  string    `json:"endpoint"`
	AccessKey string    `json:"accessKey"`
	SecretKey string    `json:"secretKey"`
	Remark    string    `json:"remark"`
	Prepath   string    `json:"prepath"`
	LoginTime time.Time `json:"loginTime"`
}
