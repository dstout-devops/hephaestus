package logger

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// recordHandler is a custom slog.Handler that captures log records for testing.
type recordHandler struct {
	records []slog.Record
}

func (h *recordHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true // Enable all log levels for testing
}

func (h *recordHandler) Handle(_ context.Context, r slog.Record) error {
	h.records = append(h.records, r)
	return nil
}

func (h *recordHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h // Simplified for testing; real impl might store attrs
}

func (h *recordHandler) WithGroup(_ string) slog.Handler {
	return h // Simplified for testing
}

// TestLogger tests the Logger interface implementation.
func TestLogger(t *testing.T) {
	// Set up a custom handler to capture logs
	handler := &recordHandler{}
	logger := &slogLogger{slog: slog.New(handler)}

	// Test Info
	logger.Info("info message", "key", "value")
	require.Len(t, handler.records, 1, "Expected one log record after Info")
	assert.Equal(t, slog.LevelInfo, handler.records[0].Level, "Expected Info level")
	assert.Equal(t, "info message", handler.records[0].Message, "Expected correct message")
	handler.records[0].Attrs(func(attr slog.Attr) bool {
		if attr.Key == "key" {
			assert.Equal(t, "value", attr.Value.Any(), "Expected key-value pair")
		}
		return true
	})

	// Test Warn
	logger.Warn("warn message", "key", "value")
	require.Len(t, handler.records, 2, "Expected two log records after Warn")
	assert.Equal(t, slog.LevelWarn, handler.records[1].Level, "Expected Warn level")
	assert.Equal(t, "warn message", handler.records[1].Message, "Expected correct message")

	// Test Error
	logger.Error("error message", "key", "value")
	require.Len(t, handler.records, 3, "Expected three log records after Error")
	assert.Equal(t, slog.LevelError, handler.records[2].Level, "Expected Error level")
	assert.Equal(t, "error message", handler.records[2].Message, "Expected correct message")

	// Test Debug
	logger.Debug("debug message", "key", "value")
	require.Len(t, handler.records, 4, "Expected four log records after Debug")
	assert.Equal(t, slog.LevelDebug, handler.records[3].Level, "Expected Debug level")
	assert.Equal(t, "debug message", handler.records[3].Message, "Expected correct message")

	// Test With
	newLogger := logger.With("contextKey", "contextValue")
	newLogger.Info("info with context")
	require.Len(t, handler.records, 5, "Expected five log records after With")
	assert.Equal(t, "info with context", handler.records[4].Message, "Expected correct message")
	handler.records[4].Attrs(func(attr slog.Attr) bool {
		if attr.Key == "contextKey" {
			assert.Equal(t, "contextValue", attr.Value.Any(), "Expected context key-value pair")
		}
		return true
	})
}

// TestNewLogger tests the NewLogger function.
func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	assert.NotNil(t, logger, "Expected NewLogger to return a non-nil logger")
	assert.Implements(t, (*Logger)(nil), logger, "Expected logger to implement Logger interface")
}
