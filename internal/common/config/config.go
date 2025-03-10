package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Kafka    KafkaConfig
	Database DBConfig
}

type ServerConfig struct {
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type DBConfig struct {
	Driver             string
	DSN                string
	Path               string
	AutoMigrate        bool
	LogLevel           string
	MaxIdleConnections int
	MaxOpenConnections int
	ConnMaxLifetime    time.Duration
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnvInt("SERVER_PORT", 8080),
			ReadTimeout:  getEnvInt("SERVER_READ_TIMEOUT", 10),
			WriteTimeout: getEnvInt("SERVER_WRITE_TIMEOUT", 10),
		},
		Kafka: KafkaConfig{
			Brokers: getEnvStringSlice("KAFKA_BROKERS", []string{"localhost:9092"}),
			Topic:   getEnvString("KAFKA_TOPIC", "task-events"),
			GroupID: getEnvString("KAFKA_GROUP_ID", "task-management-group"),
		},
		Database: DBConfig{
			Driver:             getEnvString("DB_DRIVER", "sqlite"),
			Path:               getEnvString("SQLITE_DB_PATH", "tasks.db"),
			AutoMigrate:        getEnvBool("DB_AUTO_MIGRATE", true),
			LogLevel:           getEnvString("DB_LOG_LEVEL", "warn"),
			MaxIdleConnections: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConnections: getEnvInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime:    getEnvDuration("DB_CONN_MAX_LIFETIME", time.Hour),
		},
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return value == "true" || value == "1"
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvString(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		values := strings.Split(value, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}
		return values
	}
	return defaultValue
}
