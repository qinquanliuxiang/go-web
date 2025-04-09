package apierr

import (
	"fmt"
	"net/http"
	"runtime"
)

const (
	DBErrCode = iota + 1000
	CasbinErrCode
	RedisErrCode
	LdapErrCode
	AuthErrCode
	JwtErrCode
	ForbiddenErrCode
	ParamsErrCode
	ServiceErrCode
	SonyflakeErrCode
)

var CodeMsg = map[int]string{
	DBErrCode:        "database error",
	CasbinErrCode:    "casbin error",
	RedisErrCode:     "redis error",
	LdapErrCode:      "ldap error",
	AuthErrCode:      "auth error",
	JwtErrCode:       "jwt error",
	ForbiddenErrCode: "permission denied",
	ParamsErrCode:    "params error",
	ServiceErrCode:   "service error",
	SonyflakeErrCode: "sonyflake error",
}

type ApiError struct {
	httpCode int
	Code     int
	Msg      string
	Err      error
	Stack    string
}

func (receive *ApiError) Error() string {
	return fmt.Sprintf("%s, %v", receive.Msg, receive.Err.Error())
}

// Unwrap 实现 Unwrap 方法，允许递归解包底层错误
func (receive *ApiError) Unwrap() error {
	return receive.Err
}

func InternalServer() *ApiError {
	_, file, line, _ := runtime.Caller(1)
	stack := fmt.Sprintf("%s:%d", file, line)
	return &ApiError{
		httpCode: http.StatusInternalServerError,
		Stack:    stack,
	}
}

func Unauthorized() *ApiError {
	_, file, line, _ := runtime.Caller(1)
	stack := fmt.Sprintf("%s:%d", file, line)
	return &ApiError{
		httpCode: http.StatusUnauthorized,
		Stack:    stack,
	}
}

func Forbidden() *ApiError {
	_, file, line, _ := runtime.Caller(1)
	stack := fmt.Sprintf("%s:%d", file, line)
	return &ApiError{
		httpCode: http.StatusForbidden,
		Stack:    stack,
	}
}

func BadRequest() *ApiError {
	_, file, line, _ := runtime.Caller(1)
	stack := fmt.Sprintf("%s:%d", file, line)
	return &ApiError{
		httpCode: http.StatusBadRequest,
		Stack:    stack,
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

func (receive *ApiError) WithCode(code int) *ApiError {
	receive.Code = code
	return receive
}

func (receive *ApiError) GetHttpCode() (httpCode int) {
	return receive.httpCode
}

func (receive *ApiError) Set(code int, msg string, err error) *ApiError {
	receive.Code = code
	receive.Msg = msg
	receive.Err = err
	return receive
}
