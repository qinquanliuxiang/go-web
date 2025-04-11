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

// ResponseFailure 返回失败
func (receive *BindRequest) ResponseFailure(c *gin.Context, err error) {
	c.Set(constant.LogErrMidwareKey, err)
	httpCode, res := receive.getResp(err)
	if httpCode == 0 {
		httpCode = http.StatusInternalServerError
	}
	c.JSON(httpCode, res)
}

type CheckReqOptions func(c *gin.Context, req any) error

func WithCheckJson() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBindJSON(req); err != nil {
			return apierr.BadRequest().Set(apierr.ParamsErrCode, "invalid JSON params", err)
		}
		return nil
	}
}

func WithCheckForm() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBind(req); err != nil {
			return apierr.BadRequest().Set(apierr.ParamsErrCode, "invalid form params", err)
		}
		return nil
	}
}

func WithCheckUri() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBindUri(req); err != nil {
			return apierr.BadRequest().Set(apierr.ParamsErrCode, "invalid uri params", err)
		}
		return nil
	}
}

func WithCheckQuery() CheckReqOptions {
	return func(c *gin.Context, req any) error {
		if err := c.ShouldBindQuery(req); err != nil {
			return apierr.BadRequest().Set(apierr.ParamsErrCode, "invalid query params", err)
		}
		return nil
	}
}

func (receive *BindRequest) BindAndCheck(c *gin.Context, req any, opts ...CheckReqOptions) bool {
	// 遍历所有的绑定方式，只要有一个失败就返回
	for _, opt := range opts {
		if err := opt(c, req); err != nil {
			receive.ResponseFailure(c, apierr.BadRequest().Set(apierr.ParamsErrCode, "invalid params", err))
			return true
		}
	}

	// 进行业务逻辑校验
	errMsg, err := receive.check.CheckReq(c, req)
	if err != nil {
		receive.ResponseFailure(c, apierr.BadRequest().Set(apierr.ParamsErrCode, "invalid params", err))
		return true
	}
	if len(errMsg) > 0 {
		receive.ResponseFailure(c, apierr.BadRequest().Set(apierr.ParamsErrCode, errMsg, reason.ErrParams))
		return true
	}
	return false
}

// getResp 转换错误后设置响应体
func (receive *BindRequest) getResp(err error) (httpCode int, res *response) {
	var ae *apierr.ApiError
	ok := errors.As(err, &ae)
	if !ok {
		return http.StatusInternalServerError, &response{
			Code: http.StatusInternalServerError,
		}
	}
	return ae.GetHttpCode(), &response{
		Code: ae.Code,
		Msg:  ae.Msg,
	}
}

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
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
