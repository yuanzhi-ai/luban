package comm

import (
	"crypto/md5"
	"fmt"
)

// Md5Encode 对一个字符串进行md5加密
func Md5Encode(data string) string {
	has := md5.Sum([]byte(data))
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
