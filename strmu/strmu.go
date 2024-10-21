package strmu

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

const (
	Digits         = "0123456789"
	AsciiLowercase = "abcdefghijklmnopqrstuvwxyz"
	AsciiUppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Punctuation    = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	Asciis         = Digits + AsciiLowercase + AsciiUppercase + Punctuation
)

// Contains 包含
func Contains(s string, a ...string) bool {
	for _, x := range a {
		if strings.Contains(s, x) {
			return true
		}
	}
	return false
}

// Join 连接
func Join[T any](elems []T, sep ...string) (s string) {
	_sep := append(sep, "")[0]
	for _, v := range elems {
		s += fmt.Sprintf("%v%s", v, _sep)
	}
	return strings.TrimSuffix(s, _sep)
}

// TrimPrefix 移除前缀
func TrimPrefix(a string, prefix ...string) (s string) {
	for _, x := range prefix {
		if strings.HasPrefix(a, x) {
			s = strings.TrimPrefix(a, x)
		}
	}
	return s
}

// TrimSuffix 移除后缀
func TrimSuffix(a string, suffix ...string) (s string) {
	for _, x := range suffix {
		if strings.HasPrefix(a, x) {
			s = strings.TrimPrefix(a, x)
		}
	}
	return s
}

// HasPrefix 包含前缀
func HasPrefix(s string, prefix ...string) bool {
	for _, x := range prefix {
		if strings.HasSuffix(s, x) {
			return true
		}
	}
	return false
}

// HasSuffix 包含后缀
func HasSuffix(s string, suffix ...string) bool {
	for _, x := range suffix {
		if strings.HasSuffix(s, x) {
			return true
		}
	}
	return false
}

// Rand 随机字符串
//
// n 长度, symbols 包含符号
func Rand(n int, symbols ...bool) (s string) {
	data := Asciis
	if symbols == nil {
		data = Digits + AsciiLowercase + AsciiUppercase
	}

	length := len(data)
	for i := 0; i < n; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
		s += fmt.Sprintf("%c", data[x.Int64()])
	}

	return
}

// SHA256 sha256
func SHA256(text string) string {
	b := sha256.Sum256([]byte(text))
	return strings.ToUpper(fmt.Sprintf("%x", b))
}
