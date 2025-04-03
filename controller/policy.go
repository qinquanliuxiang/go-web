package controller

import (
	"qqlx/base/handler"
	"qqlx/schema"
	"qqlx/service"

	"github.com/gin-gonic/gin"
)

type PolicyCtrl struct {
	policySvc *service.PolicySVC
	res       handler.BindResponseInterface
}

func NewPolicyCtrl(policySvc *service.PolicySVC, res *handler.BindRequest) *PolicyCtrl {
	return &PolicyCtrl{
		policySvc: policySvc,
		res:       res,
	}
}

func (receive *PolicyCtrl) GetHandler(c *gin.Context) {
	req := new(schema.PolicyIDRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri()) {
		return
	}
	res, err := receive.policySvc.GetPolicy(c, req)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}

func (receive *PolicyCtrl) CreateHandler(c *gin.Context) {
	req := new(schema.PolicyCreateRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckJson()) {
		return
	}
	if err := receive.policySvc.CreatePolicy(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *PolicyCtrl) DeleteHandler(c *gin.Context) {
	req := new(schema.PolicyIDRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri()) {
		return
	}
	if err := receive.policySvc.DeletePolicy(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *PolicyCtrl) UpdateHandler(c *gin.Context) {
	req := new(schema.PolicyUpdateRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri(), handler.WithCheckJson()) {
		return
	}
	if err := receive.policySvc.UpdatePolicy(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *PolicyCtrl) ListHandler(c *gin.Context) {
	req := new(schema.PolicyListRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckQuery()) {
		return
	}
	res, err := receive.policySvc.List(c, req)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}
