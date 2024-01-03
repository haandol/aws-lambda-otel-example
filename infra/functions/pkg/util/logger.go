package util

import (
	"log/slog"
	"os"
	"sync"
)

var logger *slog.Logger

func InitLogger(isProd bool) *slog.Logger {
	var once sync.Once
	once.Do(func() {
		if isProd {
			logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		} else {
			logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		}
	})
	logger.Info("Logger initialized", "isProd", isProd)

	return logger
}

func GetLogger() *slog.Logger {
	return logger
}
