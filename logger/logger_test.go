package logger

import (
	"log/slog"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Error("Expected logger to be created, but got nil")
	}

	_, ok := logger.Handler().(*slog.TextHandler)
	if !ok {
		t.Error("Expected logger handler to be of type *slog.TextHandler, but got a different type")
	}
	logger.Info("Logger created successfully")
}
