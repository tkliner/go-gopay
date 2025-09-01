package logger

import "context"

// Logger is an interface for structured logging with context support.
type Logger interface {
    // Info logs informational messages.
    Info(ctx context.Context, msg string, args ...any)
    // Error logs error messages.
    Error(ctx context.Context, msg string, args ...any)
    // Warn logs warning messages.
    Warn(ctx context.Context, msg string, args ...any)
    // Debug logs debug-level messages.
    Debug(ctx context.Context, msg string, args ...any)
    // Trace logs trace-level messages for detailed debugging.
    Trace(ctx context.Context, msg string, args ...any)
}