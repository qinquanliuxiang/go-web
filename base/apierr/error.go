package apierr

import (
	"fmt"
	"net/http"
	"runtime"
)

type ApiError struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg,omitempty"`
	Err   error  `json:"err,omitempty"`
	Stack string `json:"stack,omitempty"`
}

func (receive *ApiError) Error() string {
	return fmt.Sprintf("%s: %v", receive.Msg, receive.Err)
}

// Unwrap 实现 Unwrap 方法，允许递归解包底层错误
func (receive *ApiError) Unwrap() error {
	return receive.Err
}

func InternalServer() *ApiError {
	return &ApiError{
		Code: http.StatusInternalServerError,
	}
}

func Unauthorized() *ApiError {
	return &ApiError{
		Code: http.StatusUnauthorized,
	}
}

func Forbidden() *ApiError {
	return &ApiError{
		Code: http.StatusForbidden,
	}
}

func BadRequest() *ApiError {
	return &ApiError{
		Code: http.StatusBadRequest,
	}
}

func (receive *ApiError) WithStack() *ApiError {
	_, file, line, _ := runtime.Caller(1)
	receive.Stack = fmt.Sprintf("%s:%d", file, line)
	return receive
}

func (receive *ApiError) WithMsg(msg string) *ApiError {
	receive.Msg = msg
	return receive
}

func (receive *ApiError) WithErr(err error) *ApiError {
	receive.Err = err
	return receive
}
