package logger

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPrefixedLogger tests the PrefixedLogger implementation.
func TestPrefixedLogger(t *testing.T) {
	// Set up a custom handler to capture logs
	handler := &recordHandler{}
	slogLogger := slog.New(handler)
	prefix := "[TEST] "
	logger := &PrefixedLogger{prefix: prefix, logger: slogLogger}

	// Test Info
	logger.Info("info message", "key", "value")
	require.Len(t, handler.records, 1, "Expected one log record after Info")
	assert.Equal(t, slog.LevelInfo, handler.records[0].Level, "Expected Info level")
	assert.Equal(t, prefix+"info message", handler.records[0].Message, "Expected prefixed message")
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
	assert.Equal(t, prefix+"warn message", handler.records[1].Message, "Expected prefixed message")

	// Test Error
	logger.Error("error message", "key", "value")
	require.Len(t, handler.records, 3, "Expected three log records after Error")
	assert.Equal(t, slog.LevelError, handler.records[2].Level, "Expected Error level")
	assert.Equal(t, prefix+"error message", handler.records[2].Message, "Expected prefixed message")

	// Test Debug
	logger.Debug("debug message", "key", "value")
	require.Len(t, handler.records, 4, "Expected four log records after Debug")
	assert.Equal(t, slog.LevelDebug, handler.records[3].Level, "Expected Debug level")
	assert.Equal(t, prefix+"debug message", handler.records[3].Message, "Expected prefixed message")

	// Test With
	newLogger := logger.With("contextKey", "contextValue")
	newLogger.Info("info with context")
	require.Len(t, handler.records, 5, "Expected five log records after With")
	assert.Equal(t, prefix+"info with context", handler.records[4].Message, "Expected prefixed message with context")
	handler.records[4].Attrs(func(attr slog.Attr) bool {
		if attr.Key == "contextKey" {
			assert.Equal(t, "contextValue", attr.Value.Any(), "Expected context key-value pair")
		}
		return true
	})
}

// TestNewPrefixedLogger tests the NewPrefixedLogger function.
func TestNewPrefixedLogger(t *testing.T) {
	prefix := "[NEW] "
	logger := NewPrefixedLogger(prefix)
	assert.NotNil(t, logger, "Expected NewPrefixedLogger to return a non-nil logger")
	assert.Implements(t, (*Logger)(nil), logger, "Expected logger to implement Logger interface")

	// Verify the prefix is set correctly
	prefixedLogger, ok := logger.(*PrefixedLogger)
	require.True(t, ok, "Expected logger to be of type *PrefixedLogger")
	assert.Equal(t, prefix, prefixedLogger.prefix, "Expected prefix to be set correctly")
}
