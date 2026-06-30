package glerror

import (
	"errors"
)

// CodeError 表示带业务错误码的错误。
type CodeError struct {
	// Code 表示业务错误码。
	Code int
	// Message 表示错误消息。
	Message string
	// Err 表示被包装的底层错误。
	Err error
}

// New 创建带错误码和消息的错误。
func New(code int, message string) *CodeError {
	return &CodeError{
		Code:    code,
		Message: message,
	}
}

// Wrap 创建带错误码和消息的包装错误。
func Wrap(code int, message string, err error) *CodeError {
	return &CodeError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error 返回错误消息。
func (e *CodeError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// Unwrap 返回被包装的底层错误。
func (e *CodeError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// From 从错误链中提取 CodeError。
func From(err error) (*CodeError, bool) {
	var codeErr *CodeError
	if errors.As(err, &codeErr) {
		return codeErr, true
	}
	return nil, false
}

// IsCode 判断错误链中是否包含指定错误码。
func IsCode(err error, code int) bool {
	codeErr, ok := From(err)
	return ok && codeErr.Code == code
}
