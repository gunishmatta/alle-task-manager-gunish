// @title Task Management Service by Gunish - Alle
// @version 1.0.0
// @description This is a Task Management Microservice having CRUD APIs for Tasks.

// @contact.name Gunish Matta
// @contact.email gunishmatta@gmail.com

// @host localhost:8080

package main

import (
	"alle-task-manager-gunish/internal/api/handler"
	"alle-task-manager-gunish/internal/api/router"
	"alle-task-manager-gunish/internal/common/config"
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"alle-task-manager-gunish/internal/domain/repository"
	"alle-task-manager-gunish/internal/service"
	"context"
	"errors"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logger := loggingtype.NewLogger()
	logger.Info("Starting Task Management Service")

	cfg := config.LoadConfig()

	taskRepo := repository.NewMemoryTaskRepository()

	taskService := service.NewTaskService(taskRepo)

	taskHandler := handler.NewTaskHandler(taskService)

	r := router.SetupRouter(taskHandler, logger)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Info("Task Management Service is listening on port %d", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Task Management Service failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down Task Management Service")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Service forced to shutdown", "error", err)
	}
	logger.Info("Service exiting")

}
