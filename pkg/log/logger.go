package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

type Fields map[string]interface{}
type LogConfig struct {
	LogToFile bool
	FilePath  string
}

func DefaultLoggerConfig() LogConfig {
	return LogConfig{
		LogToFile: false,
		FilePath:  "",
	}
}

func ConfigureLogger(config LogConfig) error {
	if config.LogToFile {
		if config.FilePath == "" {
			return fmt.Errorf("file path must be provided when logging to file")
		}
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
		logger.SetOutput(file)
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetOutput(os.Stdout) // default to stdout
	}
	return nil
}
func SetLogLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}
	logger.Level = lvl
	return nil
}
func Debug(message interface{}, fields ...Fields) {
	logWithFields(logrus.DebugLevel, message, fields...)
}

func Info(message interface{}, fields ...Fields) {
	logWithFields(logrus.InfoLevel, message, fields...)
}

func Warn(message interface{}, fields ...Fields) {
	logWithFields(logrus.WarnLevel, message, fields...)
}

func Error(message interface{}, fields ...Fields) {
	logWithFields(logrus.ErrorLevel, message, fields...)
}

func Fatal(message interface{}, fields ...Fields) {
	logWithFields(logrus.FatalLevel, message, fields...)
}

func Panic(message interface{}, fields ...Fields) {
	logWithFields(logrus.PanicLevel, message, fields...)
}

func logWithFields(level logrus.Level, message interface{}, fields ...Fields) {
	inputFields := Fields{}
	if len(fields) > 0 && fields[0] != nil {
		inputFields = fields[0]
	}

	if logger.Level >= level {
		entry := logger.WithFields(logrus.Fields(inputFields))
		entry.Data["file"] = fileInfo(3)
		switch level {
		case logrus.DebugLevel:
			entry.Debug(message)
		case logrus.InfoLevel:
			entry.Info(message)
		case logrus.WarnLevel:
			entry.Warn(message)
		case logrus.ErrorLevel:
			entry.Error(message)
		case logrus.FatalLevel:
			entry.Fatal(message)
		case logrus.PanicLevel:
			entry.Panic(message)
		}
	}
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		file = file[strings.LastIndex(file, "/")+1:]
	}
	return fmt.Sprintf("%s:%d", file, line)
}
