package middleware

import (
	"alle-task-manager-gunish/internal/api/response"
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// Recovery : To catch any panics that might occur during request handling.
func Recovery() gin.HandlerFunc {
	logger := loggingtype.GetLogger()
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Recovery from panic",
					"error", err,
					"stack", string(debug.Stack()),
				)
				response.Error(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
