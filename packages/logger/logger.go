package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a JSON structured logger for a service.
func New(serviceName, level string) *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(level),
	})).With(
		slog.String("service", serviceName),
	)
}

// WithRequestID returns a child logger including the request ID field.
func WithRequestID(log *slog.Logger, requestID string) *slog.Logger {
	if requestID == "" {
		return log
	}
	return log.With(slog.String("request_id", requestID))
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}