package util

import (
	"fmt"
	"golang.org/x/crypto/sha3"
)

// SHA256 对文本进行SHA256加密
func SHA256(text string) string {
	str := fmt.Sprintf("%x", sha3.New256().Sum([]byte(text)))
	return str
}
