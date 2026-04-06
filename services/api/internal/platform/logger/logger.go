package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey string

const keyLogger contextKey = "logger"

// New creates an slog.Logger based on environment.
// env="production" → JSON handler; any other value → TextHandler.
func New(env, level string) *slog.Logger {
	var lvl slog.Level
	_ = lvl.UnmarshalText([]byte(level))

	opts := &slog.HandlerOptions{Level: lvl}

	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

// WithContext injects logger into context.
func WithContext(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, l)
}

// FromContext extracts logger from context; returns default if absent.
func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(keyLogger).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
