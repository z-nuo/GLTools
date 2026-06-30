package glhttp

import "github.com/z-nuo/GLTools/glerror"

// Response 表示统一 HTTP JSON 响应。
type Response[T any] struct {
	// Code 表示业务状态码。
	Code int `json:"code"`
	// Message 表示响应消息。
	Message string `json:"message"`
	// Data 表示响应数据。
	Data T `json:"data"`
}

// Success 创建成功响应。
func Success[T any](data T) Response[T] {
	return Response[T]{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

// Fail 创建失败响应。
func Fail(code int, message string) Response[any] {
	err := glerror.New(code, message)
	return Response[any]{
		Code:    err.Code,
		Message: err.Message,
		Data:    nil,
	}
}
