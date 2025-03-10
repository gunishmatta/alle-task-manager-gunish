package loggingtype

import (
	"log/slog"
	"os"
	"sync"
)

type Logger struct {
	*slog.Logger
}

var (
	instance *Logger
	once     sync.Once
)

func NewLogger() *Logger {
	once.Do(func() {
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		handler := slog.NewJSONHandler(os.Stdout, opts)
		logger := slog.New(handler)
		instance = &Logger{
			Logger: logger,
		}
	})

	return instance
}
func GetLogger() *Logger {
	if instance == nil {
		return NewLogger()
	}
	return instance
}
