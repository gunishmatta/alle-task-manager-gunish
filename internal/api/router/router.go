package router

import (
	"alle-task-manager-gunish/internal/api/handler"
	"alle-task-manager-gunish/internal/api/middleware"
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"github.com/gin-gonic/gin"
)

func SetupRouter(taskHandler *handler.TaskHandler, logger *loggingtype.Logger) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Logging(logger))
	router.Use(middleware.Recovery(logger))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	taskHandler.RegisterRoutes(router)

	return router
}
