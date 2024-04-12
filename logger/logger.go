package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func NewLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Level(getLevelFromEnv())
	return logger
}

func getLevelFromEnv() zerolog.Level {
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.InfoLevel
	}
}
