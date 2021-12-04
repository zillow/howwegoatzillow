package worker

import (
	"context"
)

var _ Logger = NoopLogger{}

// Logger ...
type Logger interface {
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Sync()
}

// NoopLogger ...
type NoopLogger struct{}

// DPanicw ...
func (n NoopLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
}

// Sync ...
func (n NoopLogger) Sync() {
}
