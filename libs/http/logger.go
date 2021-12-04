package http

import (
	"context"

	"github.com/hashicorp/go-retryablehttp"
)

// Logger conforms to our.Logger interface
type Logger interface {
	Debugw(ctx context.Context, msg string, keysAndValues ...interface{})
	Infow(ctx context.Context, msg string, keysAndValues ...interface{})
	Errorw(ctx context.Context, msg string, keysAndValues ...interface{})
	Warnw(ctx context.Context, msg string, keysAndValues ...interface{})
}

// LeveledLogger conforms to retryablehttp.LeveledLogger interface.
// This has no context support, so all log messages are logged without contextual information.
// application scoped log fields will still be added.
type LeveledLogger struct {
	logger Logger
}

var _ retryablehttp.LeveledLogger = LeveledLogger{}

func NewLeveledLogger(l Logger) LeveledLogger {
	return LeveledLogger{l}
}

func (l LeveledLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(context.Background(), msg, keysAndValues...)
}
func (l LeveledLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(context.Background(), msg, keysAndValues...)
}
func (l LeveledLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(context.Background(), msg, keysAndValues...)
}
func (l LeveledLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(context.Background(), msg, keysAndValues...)
}
