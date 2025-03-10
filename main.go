package main

import (
	"alle-task-manager-gunish/internal/api/router"
	"alle-task-manager-gunish/internal/common/config"
	"alle-task-manager-gunish/internal/common/database"
	"alle-task-manager-gunish/internal/common/dependency"
	"alle-task-manager-gunish/internal/common/logging"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := loggingtype.GetLogger()
	logger.Info("Starting Task Management Service")

	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewDatabase(ctx, cfg.Database)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer func(db *database.Database) {
		err := db.Close()
		if err != nil {
			logger.Error("Failed to close database", "error", err)
		}
	}(db)

	c, err := dependency.NewContainer(
		dependency.WithConfig(cfg),
		dependency.WithDatabase(db),
	)
	if err != nil {
		logger.Error("Failed to initialize container", "error", err)
		os.Exit(1)
	}
	defer c.Close()
	r := router.SetupRouter(c.TaskHandler())

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	go func() {
		logger.Info("Task Management Service is listening", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Task Management Service failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Task Management Service")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Service forced to shutdown", "error", err)
	}

	logger.Info("Service exited gracefully")
}
