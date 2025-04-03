package middleware

import (
	"net/http"
	"qqlx/base/apierr"
	"qqlx/base/constant"
	"qqlx/base/reason"
	"qqlx/pkg/jwt"

	"strings"

	"github.com/gin-gonic/gin"
)

const auth = "auth failed"

// Authentication 基于JWT的认证中间件
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			permissionDenied(c, map[string]any{
				"code":  http.StatusUnauthorized,
				"error": reason.ErrHeaderEmpty.Error(),
			}, apierr.Unauthorized().WithMsg(auth).WithErr(reason.ErrHeaderEmpty).WithStack())
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			permissionDenied(c, map[string]any{
				"code":  http.StatusUnauthorized,
				"error": reason.ErrHeaderMalformed.Error(),
			}, apierr.Unauthorized().WithMsg(auth).WithErr(reason.ErrHeaderMalformed).WithStack())
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			permissionDenied(c, map[string]any{
				"code":  http.StatusUnauthorized,
				"error": reason.ErrTokenInvalid.Error(),
			}, apierr.Unauthorized().WithMsg(auth).WithErr(reason.ErrTokenInvalid).WithStack())
			return
		}
		c.Set(constant.AuthMidwareKey, mc)
		c.Next()
	}
}
