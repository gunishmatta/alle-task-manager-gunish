package database

import (
	"alle-task-manager-gunish/internal/common/config"
	"alle-task-manager-gunish/internal/common/logging"
	"alle-task-manager-gunish/internal/domain/model"
	"context"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type Database struct {
	Db     *gorm.DB
	logger *loggingtype.Logger
	config config.DBConfig
}

func NewDatabase(ctx context.Context, config config.DBConfig) (*Database, error) {
	var dialector gorm.Dialector
	logger := loggingtype.GetLogger()

	switch config.Driver {
	case "sqlite":
		if config.Path == "" {
			logger.Error("sqlite database path is required")
			return nil, errors.New("sqlite database path is required")
		}
		if _, err := os.Stat(config.Path); os.IsNotExist(err) {
			file, err := os.Create(config.Path)
			if err != nil {
				logger.Error("failed to create SQLite database file", "error", err)
				return nil, errors.New("failed to create SQLite database file: " + err.Error())
			}
			err = file.Close()
			if err != nil {
				return nil, err
			}
			logger.Info("SQLite database file created", "path", config.Path)
		}
		dialector = sqlite.Open(config.Path)
	default:
		logger.Error("unsupported database driver", "driver", config.Driver)
		return nil, errors.New("unsupported database driver: " + config.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		logger.Error("failed to open database connection", "error", err)
		return nil, errors.New("failed to open database connection: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("failed to get database instance", "error", err)
		return nil, errors.New("failed to get database instance: " + err.Error())
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(config.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	database := &Database{
		Db:     db,
		logger: logger,
		config: config,
	}

	if config.AutoMigrate {
		if err := database.Db.AutoMigrate(&model.Task{}); err != nil {
			logger.Error("failed to migrate database schema", "error", err)
			return nil, errors.New("failed to migrate database schema: " + err.Error())
		}
		logger.Info("Database schema migrated successfully")
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Error("failed to ping database", "error", err)
		return nil, errors.New("failed to ping database: " + err.Error())
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
		d.logger.Error("failed to get database instance for closing", "error", err)
		return errors.New("failed to get database instance for closing: " + err.Error())
	}

	if err := sqlDB.Close(); err != nil {
		d.logger.Error("failed to close database connection", "error", err)
		return errors.New("failed to close database connection: " + err.Error())
	}

	d.logger.Info("Database connection closed")
	return nil
}
