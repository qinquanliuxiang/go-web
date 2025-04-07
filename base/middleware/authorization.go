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
	"strings"
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
		roleName, err := receive.cache.GetString(c, key)
		if err != nil {
			permissionDenied(c, map[string]any{
				"error": err.Error(),
				"code":  http.StatusForbidden,
				"msg":   AuthFailed,
			}, apierr.Unauthorized().WithMsg("get cache failed").WithErr(err).WithStack())
			return
		}
		_roleName := make([]string, 0)
		if roleName == "" {
			user, err := receive.userStore.Query(c, userstore.Name(claims.UserName))
			if err != nil {
				permissionDenied(c, map[string]any{
					"error": err.Error(),
					"code":  http.StatusForbidden,
					"msg":   AuthFailed,
				}, apierr.Unauthorized().WithMsg("get user failed").WithErr(err).WithStack())
				return
			}

			if len(user.Roles) > 0 {
				for _, role := range user.Roles {
					_roleName = append(_roleName, role.Name)
				}
				roleName = strings.Join(_roleName, ",")
			}

			_ = receive.cache.SetString(c, key, roleName, &cache.NeverExpires)
			if roleName == "" {
				permissionDenied(c, map[string]any{
					"error": reason.ErrRoleEmpty.Error(),
					"code":  http.StatusForbidden,
					"msg":   AuthFailed,
				}, apierr.Forbidden().WithMsg("role name is empty").WithErr(reason.ErrRoleEmpty).WithStack())
				return
			}
		}

		// 判断是否有权限D
		for i := range _roleName {
			ok, err := receive.authorizer.EnforceWithCtx(c, _roleName[i], c.Request.URL.Path, c.Request.Method)
			if err != nil {
				permissionDenied(c, map[string]any{
					"error": err.Error(),
					"code":  http.StatusForbidden,
					"msg":   AuthFailed,
				}, apierr.Forbidden().WithMsg("enforce failed").WithErr(err).WithStack())
				return
			}
			if !ok {
				permissionDenied(c, map[string]any{
					"error": reason.ErrPermission.Error(),
					"code":  http.StatusForbidden,
					"msg":   AuthFailed,
				}, apierr.Forbidden().WithMsg("no permission").WithErr(reason.ErrPermission).WithStack())
				logger.WithContext(c, true).Errorf("permission denied, userName: '%s', roleName: '%s', action: '%s', resource: '%s'", claims.UserName, roleName, c.Request.Method, c.Request.URL.Path)
				return
			}
		}

		c.Next()
	}
}

func permissionDenied(c *gin.Context, res map[string]any, err error) {
	c.Set(constant.LogErrMidwareKey, err)
	c.JSON(http.StatusForbidden, res)
	c.Abort()
}
