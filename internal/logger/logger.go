package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

type ctxKey struct{}

func InitLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func SetGlobal(logger *zap.Logger) {
	defaultLogger = logger
}

// FromContext retrieves the logger from the context.
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*zap.Logger); ok {
		return logger
	}

	return defaultLogger
}

// ToContext sets the logger in the context.
func ToContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

// Infof logs an info message.
func Infof(ctx context.Context, format string, args ...interface{}) {
	if logger := FromContext(ctx); logger != nil {
		logger.Sugar().Infof(format, args...)
	}
}

// Errorf logs an error message.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	if logger := FromContext(ctx); logger != nil {
		logger.Sugar().Errorf(format, args...)
	}
}

// Fatalf logs a fatal message. The os.Exit(1) will be called after logging the message.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	if logger := FromContext(ctx); logger != nil {
		logger.Sugar().Fatalf(format, args...)
	}
}
