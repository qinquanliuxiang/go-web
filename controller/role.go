package controller

import (
	"qqlx/base/handler"
	"qqlx/schema"
	"qqlx/service"

	"github.com/gin-gonic/gin"
)

type RoleCtrl struct {
	roleSvc *service.RoleSVC
	res     handler.BindResponseInterface
}

func NewRoleCtrl(roleSvc *service.RoleSVC, res *handler.BindRequest) *RoleCtrl {
	return &RoleCtrl{
		roleSvc: roleSvc,
		res:     res,
	}
}

func (receive *RoleCtrl) GetHandler(c *gin.Context) {
	req := new(schema.RoleIDRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri()) {
		return
	}
	res, err := receive.roleSvc.GetRole(c, req)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}

func (receive *RoleCtrl) CreateHandler(c *gin.Context) {
	req := new(schema.RoleCreateRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckJson()) {
		return
	}
	if err := receive.roleSvc.CreateRole(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *RoleCtrl) DeleteHandler(c *gin.Context) {
	req := new(schema.RoleIDRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri()) {
		return
	}
	if err := receive.roleSvc.DeleteRole(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *RoleCtrl) UpdateInfoHandler(c *gin.Context) {
	req := new(schema.RoleUpdateRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri(), handler.WithCheckJson()) {
		return
	}
	if err := receive.roleSvc.UpdateRoleDesc(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *RoleCtrl) AddRoleByPolicyHandler(c *gin.Context) {
	req := new(schema.RolePolicyRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri(), handler.WithCheckJson()) {
		return
	}
	if err := receive.roleSvc.AddByPolicy(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *RoleCtrl) DeleteRoleByPolicyHandler(c *gin.Context) {
	req := new(schema.RolePolicyRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri(), handler.WithCheckJson()) {
		return
	}
	if err := receive.roleSvc.DeleteByPolicy(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *RoleCtrl) ListHandler(c *gin.Context) {
	req := new(schema.RoleListRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckQuery()) {
		return
	}
	res, err := receive.roleSvc.ListRole(c, req)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}
