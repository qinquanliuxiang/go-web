package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"qqlx/base/apierr"
	"qqlx/base/constant"
	"qqlx/base/reason"
	"qqlx/base/validator"
)

type BindResponseInterface interface {
	ResponseSuccess(c *gin.Context, data any)
	ResponseFailure(c *gin.Context, err error)
	BindAndCheck(c *gin.Context, req any, opts ...CheckReqOptions) bool
}

type BindRequest struct {
	check validator.CheckReqInterface
}

func NewResponse(check validator.CheckReqInterface) *BindRequest {
	return &BindRequest{
		check: check,
	}
}

func (receive *BindRequest) ResponseSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, newResponse(data))
}

func (receive *BindRequest) ResponseFailure(c *gin.Context, err error) {
	c.Set(constant.LogErrMidwareKey, err)
	res := receive.getResp(err)
	switch res.Code {
	case http.StatusBadRequest:
		c.JSON(http.StatusBadRequest, res)
		return
	case http.StatusUnauthorized:
		c.JSON(http.StatusUnauthorized, res)
		return
	case http.StatusForbidden:
		c.JSON(http.StatusForbidden, res)
		return
	default:
		c.JSON(http.StatusInternalServerError, res)
		return
	}
}

type CheckReqOptions func(c *gin.Context, req any) error

func WithCheckJson() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBindJSON(req); err != nil {
			return apierr.BadRequest().WithMsg("invalid JSON params").WithErr(err).WithStack()
		}
		return nil
	}
}

func WithCheckForm() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBind(req); err != nil {
			return apierr.BadRequest().WithMsg("invalid Form params").WithErr(err).WithStack()
		}
		return nil
	}
}

func WithCheckUri() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBindUri(req); err != nil {
			return apierr.BadRequest().WithMsg("invalid URI params").WithErr(err).WithStack()
		}
		return nil
	}
}

func WithCheckQuery() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBindQuery(req); err != nil {
			return apierr.BadRequest().WithMsg("invalid Query params").WithErr(err).WithStack()
		}
		return nil
	}
}

func (receive *BindRequest) BindAndCheck(c *gin.Context, req any, opts ...CheckReqOptions) bool {
	// 遍历所有的绑定方式，只要有一个失败就返回
	for _, opt := range opts {
		if err := opt(c, req); err != nil {
			receive.ResponseFailure(c, apierr.BadRequest().WithMsg("invalid request").WithErr(err).WithStack())
			return true
		}
	}

	// 进行业务逻辑校验
	errMsg, err := receive.check.CheckReq(c, req)
	if err != nil {
		receive.ResponseFailure(c, apierr.BadRequest().WithMsg("invalid request").WithErr(err).WithStack())
		return true
	}
	if len(errMsg) > 0 {
		receive.ResponseFailure(c, apierr.BadRequest().WithMsg(errMsg).WithErr(reason.ErrParams).WithStack())
		return true
	}
	return false
}

// getResp 转换错误后设置响应体
func (receive *BindRequest) getResp(err error) *response {
	var ae *apierr.ApiError
	ok := errors.As(err, &ae)
	if !ok {
		return &response{
			Code: http.StatusInternalServerError,
			Err:  err.Error(),
		}
	}
	return &response{
		Code: ae.Code,
		Msg:  ae.Msg,
		Err:  ae.Err.Error(),
	}
}

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data any    `json:"data,omitempty"`
	Err  string `json:"err,omitempty"`
}

func newResponse(data any) *response {
	res := &response{
		Code: 0,
		Msg:  "success",
		Data: "",
	}
	res.Data = data
	return res
}
