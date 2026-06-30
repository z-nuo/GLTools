package gljson

import "encoding/json"

// Marshal 将值编码为 JSON 字节。
func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal 将 JSON 字节解码到目标值。
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// Valid 判断字符串是否为合法 JSON。
func Valid(s string) bool {
	return json.Valid([]byte(s))
}

// Pretty 将值编码为带两个空格缩进的 JSON 字符串。
func Pretty(v any) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
