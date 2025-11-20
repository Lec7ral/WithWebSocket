package logger

import (
	"log/slog"
	"os"
)

// New creates and configures a new structured logger (slog).
func New(level string) *slog.Logger {
	var logLevel slog.Level

	// Set the logging level based on the input string.
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create a new JSON handler that writes to standard output.
	// This format is ideal for production as it's machine-readable.
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	logger := slog.New(handler)
	return logger
}
