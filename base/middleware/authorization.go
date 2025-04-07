package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"qqlx/base/apierr"
	"qqlx/base/constant"
	"qqlx/base/helpers"
	"qqlx/base/logger"
	"qqlx/base/reason"
	"qqlx/pkg/jwt"
	"qqlx/store"
	"qqlx/store/cache"
	"qqlx/store/userstore"
)

const AuthFailed = "authentication failed"

type AuthorizationMiddleware struct {
	cache      store.CacheInterface
	authorizer store.Authorizer
	userStore  store.UserStoreInterface
}

func NewAuthorization(cache store.CacheInterface, authorizer store.Authorizer, userStore store.UserStoreInterface) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		cache:      cache,
		authorizer: authorizer,
		userStore:  userStore,
	}
}

// Authorization 基于 Casbin 的鉴权中间件
func (receive *AuthorizationMiddleware) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			allowed bool
		)
		claims, err := jwt.GetMyClaims(c)
		if err != nil {
			permissionDenied(c, map[string]any{
				"error": err.Error(),
				"code":  http.StatusForbidden,
				"msg":   AuthFailed,
			}, apierr.Unauthorized().WithMsg("authentication failed").WithErr(err).WithStack())
			return
		}
		key := helpers.GetRoleCacheKey(claims.UserName)
		roleName, err := receive.cache.GetSlice(c, key)
		if err != nil {
			permissionDenied(c, map[string]any{
				"error": err.Error(),
				"code":  http.StatusForbidden,
				"msg":   AuthFailed,
			}, apierr.Unauthorized().WithMsg("get cache failed").WithErr(err).WithStack())
			return
		}
		_roleName := make([]any, 0, 10)
		if len(roleName) == 0 {
			user, err := receive.userStore.Query(c, userstore.Name(claims.UserName))
			if err != nil {
				permissionDenied(c, map[string]any{
					"error": err.Error(),
					"code":  http.StatusForbidden,
					"msg":   AuthFailed,
				}, apierr.Unauthorized().WithMsg("get user failed").WithErr(err).WithStack())
				return
			}

			if len(user.Roles) == 0 {
				permissionDenied(c, map[string]any{
					"error": err.Error(),
					"code":  http.StatusForbidden,
					"msg":   AuthFailed,
				}, apierr.Unauthorized().WithMsg("get user roles failed").WithErr(reason.ErrRoleNotFound).WithStack())
				return
			}

			for _, role := range user.Roles {
				_roleName = append(_roleName, role.Name)
				roleName = append(roleName, role.Name)
			}
			_ = receive.cache.SetSlice(c, key, _roleName, &cache.NeverExpires)
			logger.WithContext(c, true).Debugf("user: %s, set roles: %v", user.Name, _roleName)
		}

		// 判断是否有权限D
		for _, role := range roleName {
			allowed, err = receive.authorizer.EnforceWithCtx(c, role, c.Request.URL.Path, c.Request.Method)
			if err != nil {
				permissionDenied(c, map[string]any{
					"error": reason.ErrPermission.Error(),
					"code":  http.StatusForbidden,
					"msg":   AuthFailed,
				}, apierr.Forbidden().WithMsg("no permission").WithErr(reason.ErrPermission).WithStack())
				return
			}
			if allowed {
				c.Next()
				return
			}
		}
		// 所有角色都无权限，最终拒绝
		permissionDenied(c, map[string]any{
			"error": reason.ErrPermission.Error(),
			"code":  http.StatusForbidden,
			"msg":   AuthFailed,
		}, apierr.Forbidden().WithMsg("no permission").WithErr(reason.ErrPermission).WithStack())

		logger.WithContext(c, true).Errorf("permission denied: user=%s, roles=%v, path=%s, method=%s",
			claims.UserName, roleName, c.Request.URL.Path, c.Request.Method)
		return
	}
}

func permissionDenied(c *gin.Context, res map[string]any, err error) {
	c.Set(constant.LogErrMidwareKey, err)
	c.JSON(http.StatusForbidden, res)
	c.Abort()
}
