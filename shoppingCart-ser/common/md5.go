package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)
	return hex.EncodeToString(tempStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

func MakePassword(plainpwd, salt string) string {
	return MD5Encode(plainpwd + salt)
}

func VaildPassword(plainpwd, salt string, password string) bool {
	md := Md5Encode(plainpwd + salt)
	fmt.Println(md + "      " + password)
	return md == password
}
