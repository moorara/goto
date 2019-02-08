package log

import (
	"os"
	"strings"

	kitLog "github.com/go-kit/kit/log"
	kitLevel "github.com/go-kit/kit/log/level"
)

// Level is the type for log level
type Level int

const (
	// Debug log level
	Debug Level = iota
	// Info log level
	Info
	// Warn log level
	Warn
	// Error log level
	Error
	// None log level
	None
)

// Logger wraps a go-kit Logger
type Logger struct {
	Name   string
	Level  Level
	Logger kitLog.Logger
}

// NewLogger creates a new logger
func NewLogger(logger kitLog.Logger, name, level string) *Logger {
	var lev Level
	logger = kitLog.With(
		logger,
		"logger", name,
		"timestamp", kitLog.DefaultTimestampUTC,
	)

	switch strings.ToLower(level) {
	case "debug":
		lev = Debug
		logger = kitLevel.NewFilter(logger, kitLevel.AllowDebug())
	case "info":
		lev = Info
		logger = kitLevel.NewFilter(logger, kitLevel.AllowInfo())
	case "warn":
		lev = Warn
		logger = kitLevel.NewFilter(logger, kitLevel.AllowWarn())
	case "error":
		lev = Error
		logger = kitLevel.NewFilter(logger, kitLevel.AllowError())
	case "none":
		lev = None
		logger = kitLevel.NewFilter(logger, kitLevel.AllowNone())
	}

	return &Logger{
		Name:   name,
		Level:  lev,
		Logger: logger,
	}
}

// NewNopLogger creates a new logger for testing purposes
func NewNopLogger() *Logger {
	logger := kitLog.NewNopLogger()
	return NewLogger(logger, "nop", "none")
}

// NewJSONLogger creates a new logger logging in JSON
func NewJSONLogger(name, level string) *Logger {
	logger := kitLog.NewJSONLogger(os.Stdout)
	return NewLogger(logger, name, level)
}

// NewFmtLogger creates a new logger logging using fmt format strings
func NewFmtLogger(name, level string) *Logger {
	logger := kitLog.NewLogfmtLogger(os.Stdout)
	return NewLogger(logger, name, level)
}

// With returns a new logger which always logs a set of key-value pairs
func (l *Logger) With(kv ...interface{}) *Logger {
	return &Logger{
		Name:   l.Name,
		Level:  l.Level,
		Logger: kitLog.With(l.Logger, kv...),
	}
}

// SyncLogger returns a new logger which can be used concurrently by goroutines.
// Only one goroutine is allowed to log at a time and other goroutines will block until the logger is available.
func (l *Logger) SyncLogger() *Logger {
	return &Logger{
		Name:   l.Name,
		Level:  l.Level,
		Logger: kitLog.NewSyncLogger(l.Logger),
	}
}

// Debug logs in debug level
func (l *Logger) Debug(kv ...interface{}) error {
	return kitLevel.Debug(l.Logger).Log(kv...)
}

// Info logs in info level
func (l *Logger) Info(kv ...interface{}) error {
	return kitLevel.Info(l.Logger).Log(kv...)
}

// Warn logs in warn level
func (l *Logger) Warn(kv ...interface{}) error {
	return kitLevel.Warn(l.Logger).Log(kv...)
}

// Error logs in error level
func (l *Logger) Error(kv ...interface{}) error {
	return kitLevel.Error(l.Logger).Log(kv...)
}
