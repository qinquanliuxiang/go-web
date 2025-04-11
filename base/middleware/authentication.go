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
			authenticationDenied(c, apierr.Unauthorized().Set(apierr.AuthErrCode, auth, reason.ErrHeaderEmpty))
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			authenticationDenied(c, apierr.Unauthorized().Set(apierr.AuthErrCode, auth, reason.ErrHeaderMalformed))
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			authenticationDenied(c, apierr.Unauthorized().Set(apierr.AuthErrCode, auth, err))
			return
		}
		c.Set(constant.AuthMidwareKey, mc)
		c.Next()
	}
}

func authenticationDenied(c *gin.Context, err error) {
	c.Set(constant.LogErrMidwareKey, err)
	c.JSON(http.StatusUnauthorized, newRes(apierr.AuthErrCode))
	c.Abort()
}
