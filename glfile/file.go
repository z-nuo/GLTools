package glfile

import (
	"os"
	"path/filepath"
)

// Exists 判断路径是否存在。
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsFile 判断路径是否存在且为普通文件。
func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

// IsDir 判断路径是否存在且为目录。
func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// EnsureDir 确保目录存在。
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

// ReadText 读取文本文件内容。
func ReadText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteText 写入文本文件内容，并在需要时创建父目录。
func WriteText(path string, content string) error {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// Ext 返回路径的扩展名。
func Ext(path string) string {
	return filepath.Ext(path)
}

// Join 拼接路径片段。
func Join(elem ...string) string {
	return filepath.Join(elem...)
}
