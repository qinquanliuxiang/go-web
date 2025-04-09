package middleware

import (
	"net/http"
	"qqlx/base/apierr"
	"qqlx/base/constant"
	"qqlx/base/helpers"
	"qqlx/base/logger"
	"qqlx/base/reason"
	"qqlx/model"
	"qqlx/pkg/jwt"
	"qqlx/store"
	"qqlx/store/cache"
	"qqlx/store/userstore"

	"github.com/gin-gonic/gin"
)

const authFailed = "authentication failed"

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
			user    *model.User
		)
		claims, err := jwt.GetMyClaims(c)
		if err != nil {
			permissionDenied(c, apierr.Unauthorized().Set(apierr.ForbiddenErrCode, authFailed, err))
			return
		}
		key := helpers.GetRoleCacheKey(claims.UserName)
		roleName, err := receive.cache.GetSet(c, key)
		if err != nil {
			permissionDenied(c, apierr.Unauthorized().Set(apierr.ForbiddenErrCode, authFailed, err))
			return
		}
		_roleName := make([]any, 0, 10)
		if len(roleName) == 0 {
			user, err = receive.userStore.Query(c, userstore.Name(claims.UserName), userstore.LoadRoles())
			if err != nil {
				permissionDenied(c, apierr.Unauthorized().Set(apierr.ForbiddenErrCode, authFailed, err))
				return
			}

			if len(user.Roles) == 0 {
				permissionDenied(c, apierr.Unauthorized().Set(apierr.ForbiddenErrCode, authFailed, reason.ErrRoleNotFound))
				return
			}

			for _, role := range user.Roles {
				_roleName = append(_roleName, role.Name)
				roleName = append(roleName, role.Name)
			}
			_ = receive.cache.SetSet(c, key, _roleName, &cache.NeverExpires)
			logger.WithContext(c, true).Debugf("user: %s, set roles: %v", user.Name, _roleName)
		}

		// 判断是否有权限D
		for _, role := range roleName {
			allowed, err = receive.authorizer.EnforceWithCtx(c, role, c.Request.URL.Path, c.Request.Method)
			if err != nil {
				permissionDenied(c, apierr.Forbidden().Set(apierr.ForbiddenErrCode, "unknown error", err))
				return
			}
			if allowed {
				c.Next()
				return
			}
		}
		// 所有角色都无权限，最终拒绝
		permissionDenied(c, apierr.Forbidden().Set(apierr.ForbiddenErrCode, authFailed, reason.ErrPermission))
		logger.WithContext(c, true).Errorf("permission denied: user=%s, roles=%v, path=%s, method=%s",
			claims.UserName, roleName, c.Request.URL.Path, c.Request.Method)
	}
}

func permissionDenied(c *gin.Context, err error) {
	c.Set(constant.LogErrMidwareKey, err)
	c.JSON(http.StatusForbidden, newRes(apierr.ForbiddenErrCode))
	c.Abort()
}

func newRes(code int) map[string]any {
	return map[string]any{
		"code": code,
		"msg":  authFailed,
		"data": nil,
	}
}
