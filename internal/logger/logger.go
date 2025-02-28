package logger

import (
	"log/slog"
	"os"
)

// Logger defines the logging interface.
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	With(args ...any) Logger
}

// slogLogger is a wrapper around *slog.Logger that implements Logger.
type slogLogger struct {
	slog *slog.Logger
}

// Info delegates to slog.Logger's Info method.
func (l *slogLogger) Info(msg string, args ...any) {
	l.slog.Info(msg, args...)
}

// Warn delegates to slog.Logger's Warn method.
func (l *slogLogger) Warn(msg string, args ...any) {
	l.slog.Warn(msg, args...)
}

// Error delegates to slog.Logger's Error method.
func (l *slogLogger) Error(msg string, args ...any) {
	l.slog.Error(msg, args...)
}

// Debug delegates to slog.Logger's Debug method.
func (l *slogLogger) Debug(msg string, args ...any) {
	l.slog.Debug(msg, args...)
}

// With returns a new Logger with additional context.
func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{slog: l.slog.With(args...)}
}

// NewLogger creates a new Logger instance.
func NewLogger() Logger {
	sl := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return &slogLogger{slog: sl}
}
