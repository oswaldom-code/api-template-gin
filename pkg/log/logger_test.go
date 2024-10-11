package log

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	err := SetLogLevel("info")
	assert.NoError(t, err)
	assert.Equal(t, logrus.InfoLevel, logger.Level)

	err = SetLogLevel("debug")
	assert.NoError(t, err)
	assert.Equal(t, logrus.DebugLevel, logger.Level)

	err = SetLogLevel("invalid")
	assert.Error(t, err)
}

func TestInfoLog(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	SetLogLevel("info")

	Info("Info message")
	assert.Contains(t, buf.String(), "Info message")

	buf.Reset()

	Info("Info with fields", Fields{"myKey": "myValue"})
	assert.Contains(t, buf.String(), "Info with fields")
	assert.Contains(t, buf.String(), "myKey=myValue")
}

func TestDebugLog(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	SetLogLevel("debug")

	Debug("Debug message")
	assert.Contains(t, buf.String(), "Debug message")
}

func TestLogWithFields(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	SetLogLevel("info")

	Info("Log message with fields", Fields{"key": "value"})
	assert.Contains(t, buf.String(), "key=value")
	assert.Contains(t, buf.String(), "Log message with fields")
}

func TestLogLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	SetLogLevel("warning")

	Info("This info message should not appear")

	Warn("This warning message should appear")

	assert.NotContains(t, buf.String(), "This info message should not appear")
	assert.Contains(t, buf.String(), "This warning message should appear")
}
