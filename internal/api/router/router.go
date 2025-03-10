package router

import (
	"alle-task-manager-gunish/internal/api/handler"
	"alle-task-manager-gunish/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(taskHandler *handler.TaskHandler) *gin.Engine {
	router := gin.New()

	router.Use(middleware.Logging())
	router.Use(middleware.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	taskHandler.RegisterRoutes(router)

	return router
}
