package database

import (
	"alle-task-manager-gunish/internal/common/config"
	"alle-task-manager-gunish/internal/common/logging"
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	Db     *gorm.DB
	logger *loggingtype.Logger
	config config.DBConfig
}

func NewDatabase(ctx context.Context, config config.DBConfig) (*Database, error) {
	var dialector gorm.Dialector
	switch config.Driver {
	case "sqlite":
		if config.Path == "" {
			return nil, errors.New("sqlite database path is required")
		}
		dialector = sqlite.Open(config.Path)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	logger := loggingtype.GetLogger()
	database := &Database{
		Db:     db,
		logger: logger,
		config: config,
	}

	if config.AutoMigrate {
		if err := database.Db.AutoMigrate(); err != nil {
			return nil, fmt.Errorf("failed to migrate database schema: %w", err)
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Connected to database",
		"driver", config.Driver,
		"maxOpenConns", config.MaxOpenConnections,
		"connMaxLifetime", config.ConnMaxLifetime)

	return database, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.Db.DB()
	if err != nil {
		logErr := fmt.Errorf("failed to get database instance for closing: %w", err)
		d.logger.Error(logErr.Error())
		return logErr
	}
	return sqlDB.Close()
}
