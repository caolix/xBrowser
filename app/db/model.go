package db

import "time"

type LoginInfo struct {
	Endpoint  string    `json:"endpoint"`
	AccessKey string    `json:"ak"`
	SecretKey string    `json:"sk"`
	Remark    string    `json:"remark"`
	Prepath   string    `json:"prepath"`
	LoginTime time.Time `json:"loginTime"`
}
