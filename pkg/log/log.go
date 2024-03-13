package log

import (
	"context"
	"log/slog"
	"os"
)

type key struct{}

// WithContext returns a new context with the logger added.
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, key{}, logger)
}

// FromContext returns the logger in the context if it exists, otherwise a new logger is returned.
func FromContext(ctx context.Context) *slog.Logger {
	logger := ctx.Value(key{})
	if l, ok := logger.(*slog.Logger); ok {
		return l
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
