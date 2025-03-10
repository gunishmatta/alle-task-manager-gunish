package middleware

import (
	"alle-task-manager-gunish/internal/common/logging"
	"github.com/gin-gonic/gin"
	"time"
)

func Logging() gin.HandlerFunc {
	logger := loggingtype.GetLogger()
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		if statusCode >= 500 {
			logger.Error("Request",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency.String(),
				"client_ip", c.ClientIP(),
			)
		} else {
			logger.Info("Request",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", latency.String(),
				"client_ip", c.ClientIP(),
			)
		}
	}
}
