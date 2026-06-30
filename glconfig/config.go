package glconfig

import (
	"os"
	"strconv"

	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v3"
)

var jsonAPI = jsoniter.ConfigCompatibleWithStandardLibrary

// LoadJSON 从 JSON 文件加载配置到目标值。
func LoadJSON(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return jsonAPI.Unmarshal(data, out)
}

// LoadYAML 从 YAML 文件加载配置到目标值。
func LoadYAML(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

// Env 读取环境变量，缺失时返回默认值。
func Env(key string, def string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return value
}

// EnvInt 读取整数环境变量，缺失或解析失败时返回默认值。
func EnvInt(key string, def int) int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	got, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return got
}

// EnvBool 读取布尔环境变量，缺失或解析失败时返回默认值。
func EnvBool(key string, def bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	got, err := strconv.ParseBool(value)
	if err != nil {
		return def
	}
	return got
}
