package controller

import (
	"qqlx/base/apierr"
	"qqlx/base/handler"
	"qqlx/base/reason"
	"qqlx/pkg/jwt"
	"qqlx/schema"
	"qqlx/service"
	"regexp"

	"github.com/gin-gonic/gin"
)

type UserCtrl struct {
	userSvc *service.UserSVC
	res     handler.BindResponseInterface
}

func NewUserCtrl(userContr *service.UserSVC, res *handler.BindRequest) *UserCtrl {
	return &UserCtrl{
		userSvc: userContr,
		res:     res,
	}
}

func (receive *UserCtrl) RegisterHandler(c *gin.Context) {
	req := new(schema.UserRegistryRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckJson()) {
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(req.Name) {
		receive.res.ResponseFailure(c, apierr.BadRequest().WithErr(reason.ErrNameInvalid))
		return
	}

	if err := receive.userSvc.RegistryUser(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *UserCtrl) LoginHandler(c *gin.Context) {
	req := new(schema.UserLoginRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckJson()) {
		return
	}
	res, err := receive.userSvc.Login(c, req)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}

func (receive *UserCtrl) LogoutHandler(c *gin.Context) {
	claims, err := jwt.GetMyClaims(c)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	if err = receive.userSvc.Logout(c, claims.UserID); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *UserCtrl) UpdateHandler(c *gin.Context) {
	req := new(schema.UserUpdateRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckJson()) {
		return
	}
	claims, err := jwt.GetMyClaims(c)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	req.ID = claims.UserID
	if err = receive.userSvc.UpdateUser(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *UserCtrl) UpdatePasswordHandler(c *gin.Context) {
	req := new(schema.UserUpdatePasswordRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckJson()) {
		return
	}
	claims, err := jwt.GetMyClaims(c)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	req.ID = claims.UserID
	if err = receive.userSvc.UpdatePassword(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *UserCtrl) UpdateRoleHandler(c *gin.Context) {
	req := new(schema.UserUpdateRoleRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri()) {
		return
	}
	if err := receive.userSvc.UpdateUserRole(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *UserCtrl) DeleteHandler(c *gin.Context) {
	req := new(schema.UserIDRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri()) {
		return
	}
	if err := receive.userSvc.DeleteUser(c, req); err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, nil)
}

func (receive *UserCtrl) ListHandler(c *gin.Context) {
	req := new(schema.UserListRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckQuery()) {
		return
	}
	res, err := receive.userSvc.ListUser(c, req)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}

// GetUserInfoHandler 根据ID获取用户
func (receive *UserCtrl) GetUserInfoHandler(c *gin.Context) {
	req := new(schema.UserIDRequest)
	if receive.res.BindAndCheck(c, req, handler.WithCheckUri()) {
		return
	}
	res, err := receive.userSvc.Info(c, req.ID)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}

// InfoHandler 获取当前登录用户信息
func (receive *UserCtrl) InfoHandler(c *gin.Context) {
	claims, err := jwt.GetMyClaims(c)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	res, err := receive.userSvc.Info(c, claims.UserID)
	if err != nil {
		receive.res.ResponseFailure(c, err)
		return
	}
	receive.res.ResponseSuccess(c, res)
}
