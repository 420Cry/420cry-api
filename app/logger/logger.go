// Package logger provides structured logging functionality for the application.
package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with application-specific configuration
type Logger struct {
	*logrus.Logger
}

// LogLevel represents the logging level
type LogLevel string

const (
	// DebugLevel represents debug log level
	DebugLevel LogLevel = "debug"
	// InfoLevel represents info log level
	InfoLevel LogLevel = "info"
	// WarnLevel represents warn log level
	WarnLevel LogLevel = "warn"
	// ErrorLevel represents error log level
	ErrorLevel LogLevel = "error"
)

// NewLogger creates a new structured logger
func NewLogger() *Logger {
	logger := logrus.New()

	// Set default formatter to JSON for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set output to stdout
	logger.SetOutput(os.Stdout)

	// Set log level based on environment
	level := getLogLevel()
	logger.SetLevel(level)

	return &Logger{logger}
}

// getLogLevel determines the log level from environment variables
func getLogLevel() logrus.Level {
	envLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch envLevel {
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		// Default to info level
		return logrus.InfoLevel
	}
}

// WithField adds a field to the logger context
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger context
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithError adds an error to the logger context
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// LogRequest logs HTTP request information
func (l *Logger) LogRequest(method, path, clientIP string, statusCode int, duration time.Duration) {
	l.WithFields(logrus.Fields{
		"type":        "http_request",
		"method":      method,
		"path":        path,
		"client_ip":   clientIP,
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}).Info("HTTP request processed")
}

// LogDatabaseOperation logs database operation information
func (l *Logger) LogDatabaseOperation(operation, table string, duration time.Duration, err error) {
	fields := logrus.Fields{
		"type":        "database_operation",
		"operation":   operation,
		"table":       table,
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		l.WithFields(fields).WithError(err).Error("Database operation failed")
	} else {
		l.WithFields(fields).Info("Database operation completed")
	}
}

// LogBusinessOperation logs business logic operation information
func (l *Logger) LogBusinessOperation(operation string, userID interface{}, duration time.Duration, err error) {
	fields := logrus.Fields{
		"type":        "business_operation",
		"operation":   operation,
		"user_id":     userID,
		"duration_ms": duration.Milliseconds(),
	}

	if err != nil {
		l.WithFields(fields).WithError(err).Error("Business operation failed")
	} else {
		l.WithFields(fields).Info("Business operation completed")
	}
}

// LogSecurityEvent logs security-related events
func (l *Logger) LogSecurityEvent(event string, clientIP string, userID interface{}, details map[string]interface{}) {
	fields := logrus.Fields{
		"type":      "security_event",
		"event":     event,
		"client_ip": clientIP,
		"user_id":   userID,
	}

	// Add additional details if provided
	for key, value := range details {
		fields[key] = value
	}

	l.WithFields(fields).Warn("Security event occurred")
}

// SetLogLevel changes the log level at runtime
func (l *Logger) SetLogLevel(level LogLevel) {
	switch level {
	case DebugLevel:
		l.SetLevel(logrus.DebugLevel)
	case InfoLevel:
		l.SetLevel(logrus.InfoLevel)
	case WarnLevel:
		l.SetLevel(logrus.WarnLevel)
	case ErrorLevel:
		l.SetLevel(logrus.ErrorLevel)
	}
}

// Global logger instance
var defaultLogger *Logger

// GetLogger returns the global logger instance
func GetLogger() *Logger {
	if defaultLogger == nil {
		defaultLogger = NewLogger()
	}
	return defaultLogger
}

// SetDefaultLogger sets the global logger instance
func SetDefaultLogger(logger *Logger) {
	defaultLogger = logger
}
