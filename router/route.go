package router

import (
	"qqlx/base/middleware"
	"qqlx/controller"

	"github.com/gin-gonic/gin"
)

type ApiRoute struct {
	userCtrl   *controller.UserCtrl
	roleCtrl   *controller.RoleCtrl
	policyCtrl *controller.PolicyCtrl
}

func NewApiRoute(
	userContr *controller.UserCtrl,
	roleContr *controller.RoleCtrl,
	policyController *controller.PolicyCtrl,
) *ApiRoute {
	return &ApiRoute{
		userCtrl:   userContr,
		roleCtrl:   roleContr,
		policyCtrl: policyController,
	}
}

func (a *ApiRoute) RegisterApiUserRoute(r *gin.RouterGroup, authorization *middleware.AuthorizationMiddleware) {
	userGroup := r.Group("/users")
	{
		userGroup.POST("create", a.userCtrl.RegisterHandler)
		userGroup.POST("/login", a.userCtrl.LoginHandler)
		userGroup.Use(middleware.Authentication())
		{
			userGroup.POST("/logout", a.userCtrl.LogoutHandler)
			userGroup.GET("", authorization.Authorization(), a.userCtrl.ListHandler)
			userGroup.PATCH("", a.userCtrl.UpdatePasswordHandler)
			userGroup.PUT("", a.userCtrl.UpdateHandler)
			userGroup.GET("/info", a.userCtrl.InfoHandler)
			userGroup.GET("/:id", authorization.Authorization(), a.userCtrl.GetUserInfoHandler)
			userGroup.DELETE("/:id", authorization.Authorization(), a.userCtrl.DisableHandler)
			userGroup.PUT("/enable/:id", a.userCtrl.EnableHandler)
			userGroup.PUT("/:id/roles", authorization.Authorization(), a.userCtrl.AddRoleHandler)
			userGroup.POST("/:id/roles", authorization.Authorization(), a.userCtrl.RemoveRoleHandler)
		}
	}
}

func (a *ApiRoute) RegisterApiRoleRoute(r *gin.RouterGroup, authorization *middleware.AuthorizationMiddleware) {
	roleGroup := r.Group("/roles")
	roleGroup.Use(middleware.Authentication(), authorization.Authorization())
	roleGroup.GET("", a.roleCtrl.ListHandler)
	roleGroup.POST("", a.roleCtrl.CreateHandler)
	roleGroup.PUT("/:id", a.roleCtrl.UpdateInfoHandler)
	roleGroup.GET("/:id", a.roleCtrl.GetHandler)
	roleGroup.DELETE("/:id", a.roleCtrl.DeleteHandler)
	roleGroup.PUT("/:id/polices", a.roleCtrl.AddRoleByPolicyHandler)
	roleGroup.POST("/:id/polices", a.roleCtrl.DeleteRoleByPolicyHandler)
}

func (a *ApiRoute) RegisterApiPolicyRoute(r *gin.RouterGroup, authorization *middleware.AuthorizationMiddleware) {
	poliyGroup := r.Group("/polices")
	poliyGroup.Use(middleware.Authentication(), authorization.Authorization())
	poliyGroup.GET("", a.policyCtrl.ListHandler)
	poliyGroup.POST("", a.policyCtrl.CreateHandler)
	poliyGroup.GET("/:id", a.policyCtrl.GetHandler)
	poliyGroup.PUT("/:id", a.policyCtrl.UpdateHandler)
	poliyGroup.DELETE("/:id", a.policyCtrl.DeleteHandler)
}
