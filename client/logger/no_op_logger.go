package logger

import "context"

// NoOpLogger is a logger implementation that does nothing.
// Useful for disabling logging in tests or production.
type NoOpLogger struct{}

// NewNoOpLogger creates a new instance of NoOpLogger.
func NewNoOpLogger() *NoOpLogger {
    return &NoOpLogger{}
}

// Info does nothing. Implements Logger interface.
func (l *NoOpLogger) Info(ctx context.Context, msg string, args ...any) {}

// Error does nothing. Implements Logger interface.
func (l *NoOpLogger) Error(ctx context.Context, msg string, args ...any) {}

// Warn does nothing. Implements Logger interface.
func (l *NoOpLogger) Warn(ctx context.Context, msg string, args ...any) {}

// Debug does nothing. Implements Logger interface.
func (l *NoOpLogger) Debug(ctx context.Context, msg string, args ...any) {}

// Trace does nothing. Implements Logger interface.
func (l *NoOpLogger) Trace(ctx context.Context, msg string, args ...any) {}