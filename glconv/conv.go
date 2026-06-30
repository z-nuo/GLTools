package glconv

import (
	"fmt"
	"strconv"
	"strings"
)

// ToInt 将字符串转换为 int。
func ToInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// ToIntDefault 将字符串转换为 int，失败时返回默认值。
func ToIntDefault(s string, def int) int {
	v, err := ToInt(s)
	if err != nil {
		return def
	}
	return v
}

// ToInt64 将字符串转换为 int64。
func ToInt64(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}

// ToFloat64 将字符串转换为 float64。
func ToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

// ToBool 将字符串转换为 bool。
func ToBool(s string) (bool, error) {
	return strconv.ParseBool(strings.TrimSpace(s))
}

// String 将任意值转换为字符串。
func String(v any) string {
	return fmt.Sprint(v)
}
