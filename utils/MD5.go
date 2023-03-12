package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	temp := h.Sum(nil)
	return hex.EncodeToString(temp)
}

func MakePassword(oldPassword string, salt string) string {
	return Md5Encode(oldPassword + salt)
}

func CheckPassword(nowPassword, oldPassword, salt string) bool {
	return MakePassword(nowPassword, salt) == oldPassword
}
