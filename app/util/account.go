package util

import (
	"crypto/md5"
	"fmt"
)

func GenAccountId(accessKey, endpoint string) string {
	sig := accessKey + endpoint
	newSig := md5.Sum([]byte(sig)) //转成加密编码
	// 将编码转换为字符串
	return fmt.Sprintf("%x", newSig)
}
