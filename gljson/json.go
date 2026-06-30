package gljson

import jsoniter "github.com/json-iterator/go"

var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

// Marshal 将值编码为 JSON 字节。
func Marshal(v any) ([]byte, error) {
	return jsonAPI.Marshal(v)
}

// Unmarshal 将 JSON 字节解码到目标值。
func Unmarshal(data []byte, v any) error {
	return jsonAPI.Unmarshal(data, v)
}

// Valid 判断字符串是否为合法 JSON。
func Valid(s string) bool {
	return jsonAPI.Valid([]byte(s))
}

// Pretty 将值编码为带两个空格缩进的 JSON 字符串。
func Pretty(v any) (string, error) {
	data, err := jsonAPI.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
