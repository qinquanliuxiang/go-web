package middleware

import (
	"fmt"
	"qqlx/base/apierr"
	"qqlx/base/constant"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ZapMiddleware 日志中间件
func ZapMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latency := endTime.Sub(startTime).Milliseconds()
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		query := c.Request.URL.RawQuery
		method := c.Request.Method
		path := c.Request.URL.Path

		traceID := GetTraceID(c)
		fields := []zap.Field{
			zap.String("traceID", traceID),
			zap.Int("status", statusCode),
			zap.String("clientIP", clientIP),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("latency", fmt.Sprintf("%dms", latency)),
		}

		_err, ok := c.Get(constant.LogErrMidwareKey)
		if ok {
			switch e := _err.(type) {
			case *apierr.ApiError:
				caller := e.Stack
				zap.Int("code", e.Code)
				m, exist := apierr.CodeMsg[e.Code]
				if exist {
					fields = append(fields, zap.String("type", m))
				}
				fields = append(fields,
					zap.String("error", e.Error()),
					zap.String("caller", caller),
				)
			case error:
				fields = append(fields,
					zap.String("error", e.Error()),
				)
			}
		}

		if ok {
			zap.L().Error("request failed", fields...)
		} else {
			zap.L().Info("request success", fields...)
		}
	}
}
