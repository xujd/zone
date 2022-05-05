package common

import (
	"fmt"
)

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewHTTPError(code int, msg string) *HttpError {
	return &HttpError{
		Code:    code,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *HttpError) Error() string {
	return fmt.Sprintf("%d:%s", e.Code, e.Message)
}

const (
	// SUCCESS 成功
	SUCCESS int = 20000
	// ERR_NOT_FOUND 未找到
	ERR_NOT_FOUND int = 40000
	// ERR_INTERNAL_SERVER_ERROR 内部错误
	ERR_INTERNAL_SERVER_ERROR int = 40003
	// ERR_BAD_REQUEST 错误请求
	ERR_BAD_REQUEST int = 40005
	// ERR_TOKEN_EXPIRED Token超时
	ERR_TOKEN_EXPIRED int = 50014
	// ERR_ILLEGAL_TOKEN 无效的Token
	ERR_ILLEGAL_TOKEN int = 50008
)
