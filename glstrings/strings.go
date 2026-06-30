package glstrings

import (
	"strings"
	"unicode"
)

// IsBlank 判断字符串在去除首尾空白后是否为空。
func IsBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Trim 去除字符串首尾空白字符。
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// Truncate 按 rune 数量截断字符串，避免截断多字节字符。
func Truncate(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}

	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes])
}

// SnakeToCamel 将 snake_case 字符串转换为 lowerCamelCase。
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	var b strings.Builder
	first := true

	for _, part := range parts {
		if part == "" {
			continue
		}
		runes := []rune(strings.ToLower(part))
		if len(runes) == 0 {
			continue
		}
		if first {
			b.WriteString(string(runes))
			first = false
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}

	return b.String()
}

// CamelToSnake 将 camelCase 或 PascalCase 字符串转换为 snake_case。
func CamelToSnake(s string) string {
	runes := []rune(s)
	if len(runes) == 0 {
		return ""
	}

	var b strings.Builder
	for i, r := range runes {
		if unicode.IsUpper(r) {
			prevIsLowerOrDigit := i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1]))
			nextIsLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			if i > 0 && (prevIsLowerOrDigit || nextIsLower) {
				b.WriteRune('_')
			}
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
