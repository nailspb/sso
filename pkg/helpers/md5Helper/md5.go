package md5Helper

import (
	"crypto/md5"
	"encoding/hex"
)

func GetMD5String(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
