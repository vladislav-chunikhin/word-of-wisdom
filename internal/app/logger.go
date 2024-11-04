package app

import (
	"log"
	"log/slog"
	"os"
)

// SetLogger sets the logger for the application
func SetLogger(logLevelAsStr string) {
	logLevel := slog.LevelDebug
	switch logLevelAsStr {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		log.Printf("unknown log level: %s, using default log level: DEBUG", logLevelAsStr)
	}

	opts := &slog.HandlerOptions{Level: logLevel}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
