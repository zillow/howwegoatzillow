package server

import (
	"context"
)

var _ Logger = NoopLogger{}

// Logger is a local interface for logging functionality
type Logger interface {
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
}

// NoopLogger is a noop logger implementation.
type NoopLogger struct{}

// Debug ...
func (n NoopLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// Error ...
func (n NoopLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// Info ...
func (n NoopLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {}

// DPanic ...
func (n NoopLogger) DPanic(ctx context.Context, msg string, keysAndValues ...interface{}) {}
