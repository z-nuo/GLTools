package glrand

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

const (
	numericCharset = "0123456789"
	stringCharset  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// NumericCode 生成指定长度的数字随机码。
func NumericCode(length int) (string, error) {
	return StringFromCharset(length, numericCharset)
}

// String 生成指定长度的字母数字随机字符串。
func String(length int) (string, error) {
	return StringFromCharset(length, stringCharset)
}

// StringFromCharset 使用指定字符集生成指定长度的随机字符串。
func StringFromCharset(length int, charset string) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than 0")
	}
	if charset == "" {
		return "", errors.New("charset must not be empty")
	}

	runes := []rune(charset)
	var b strings.Builder
	b.Grow(length)
	max := big.NewInt(int64(len(runes)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b.WriteRune(runes[n.Int64()])
	}

	return b.String(), nil
}
