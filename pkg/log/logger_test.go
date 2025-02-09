package log

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		level       string
		expected    logrus.Level
		expectError bool
	}{
		{"info", logrus.InfoLevel, false},
		{"debug", logrus.DebugLevel, false},
		{"invalid", logrus.InfoLevel, true}, // Mantiene el nivel anterior si es inválido
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			err := SetLogLevel(tt.level)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, logger.Level)
			}
		})
	}
}

func TestInfoLog(t *testing.T) {
	t.Run("Without Fields", func(t *testing.T) {
		var buf bytes.Buffer
		originalOutput := logger.Out
		defer func() { logger.SetOutput(originalOutput) }()
		logger.SetOutput(&buf)

		SetLogLevel("info")

		Info("Info message")
		assert.Contains(t, buf.String(), "Info message")
	})

	t.Run("With Fields", func(t *testing.T) {
		var buf bytes.Buffer
		originalOutput := logger.Out
		defer func() { logger.SetOutput(originalOutput) }()
		logger.SetOutput(&buf)

		SetLogLevel("info")

		Info("Info with fields", Fields{"myKey": "myValue"})
		assert.Contains(t, buf.String(), "Info with fields")
		assert.Contains(t, buf.String(), "myKey=myValue")
	})
}

func TestDebugLog(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := logger.Out
	defer func() { logger.SetOutput(originalOutput) }()
	logger.SetOutput(&buf)

	SetLogLevel("debug")

	Debug("Debug message")
	assert.Contains(t, buf.String(), "Debug message")
}

func TestLogWithFields(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := logger.Out
	defer func() { logger.SetOutput(originalOutput) }()
	logger.SetOutput(&buf)

	SetLogLevel("info")

	Info("Log message with fields", Fields{"key": "value"})
	assert.Contains(t, buf.String(), "key=value")
	assert.Contains(t, buf.String(), "Log message with fields")
}

func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := logger.Out
	defer func() { logger.SetOutput(originalOutput) }()
	logger.SetOutput(&buf)

	SetLogLevel("warning")

	Info("This info message should not appear")
	Warn("This warning message should appear")

	assert.NotContains(t, buf.String(), "This info message should not appear")
	assert.Contains(t, buf.String(), "This warning message should appear")
}
