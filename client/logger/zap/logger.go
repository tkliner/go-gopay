package zap

import (
    "context"

    "go.uber.org/zap"
)

// ZapLogger is an implementation of the Logger interface using Uber's zap.Logger.
type ZapLogger struct {
    zap *zap.Logger
}

// NewZapLogger creates a new ZapLogger instance.
func NewZapLogger(zap *zap.Logger) *ZapLogger {
    return &ZapLogger{zap: zap}
}

// Info logs informational messages.
func (l *ZapLogger) Info(ctx context.Context, msg string, args ...any) {
    fields := make([]zap.Field, 0, len(args)/2)
    for i := 0; i < len(args); i += 2 {
        if i+1 < len(args) {
            fields = append(fields, zap.Any(args[i].(string), args[i+1]))
        }
    }
    l.zap.Info(msg, fields...)
}

// Error logs error messages.
func (l *ZapLogger) Error(ctx context.Context, msg string, args ...any) {
    fields := make([]zap.Field, 0, len(args)/2)
    for i := 0; i < len(args); i += 2 {
        if i+1 < len(args) {
            fields = append(fields, zap.Any(args[i].(string), args[i+1]))
        }
    }
    l.zap.Error(msg, fields...)
}

// Warn logs warning messages.
func (l *ZapLogger) Warn(ctx context.Context, msg string, args ...any) {
    fields := make([]zap.Field, 0, len(args)/2)
    for i := 0; i < len(args); i += 2 {
        if i+1 < len(args) {
            fields = append(fields, zap.Any(args[i].(string), args[i+1]))
        }
    }
    l.zap.Warn(msg, fields...)
}

// Debug logs debug-level messages.
func (l *ZapLogger) Debug(ctx context.Context, msg string, args ...any) {
    fields := make([]zap.Field, 0, len(args)/2)
    for i := 0; i < len(args); i += 2 {
        if i+1 < len(args) {
            fields = append(fields, zap.Any(args[i].(string), args[i+1]))
        }
    }
    l.zap.Debug(msg, fields...)
}

// Trace logs trace-level messages for detailed debugging.
func (l *ZapLogger) Trace(ctx context.Context, msg string, args ...any) {
    fields := make([]zap.Field, 0, len(args)/2)
    for i := 0; i < len(args); i += 2 {
        if i+1 < len(args) {
            fields = append(fields, zap.Any(args[i].(string), args[i+1]))
        }
    }
    // Zap does not have a Trace level, so we use Debug for trace messages.
    l.zap.Debug("[TRACE] "+msg, fields...)
}