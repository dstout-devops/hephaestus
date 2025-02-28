package logger

import (
	"log/slog"
	"os"
)

type PrefixedLogger struct {
	prefix string
	logger *slog.Logger
}

func (l *PrefixedLogger) Info(msg string, args ...any) {
	l.logger.Info(l.prefix+msg, args...)
}

func (l *PrefixedLogger) Warn(msg string, args ...any) {
	l.logger.Warn(l.prefix+msg, args...)
}

func (l *PrefixedLogger) Error(msg string, args ...any) {
	l.logger.Error(l.prefix+msg, args...)
}

func (l *PrefixedLogger) Debug(msg string, args ...any) {
	l.logger.Debug(l.prefix+msg, args...)
}

func (l *PrefixedLogger) With(args ...any) Logger {
	return &PrefixedLogger{
		prefix: l.prefix,
		logger: l.logger.With(args...),
	}
}

// NewPrefixedLogger creates a logger with a custom prefix.
func NewPrefixedLogger(prefix string) Logger {
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	return &PrefixedLogger{
		prefix: prefix,
		logger: slogLogger,
	}
}
