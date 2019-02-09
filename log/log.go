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

type (
	// Options contains optional options for Logger
	Options struct {
		Level       string
		Name        string
		Environment string
		Region      string
	}

	// Logger wraps a go-kit Logger
	Logger struct {
		Level  Level
		Logger kitLog.Logger
	}
)

// NewNopLogger creates a new logger for testing purposes
func NewNopLogger() *Logger {
	logger := kitLog.NewNopLogger()
	return &Logger{
		Logger: logger,
	}
}

// NewLogger creates a new logger
func NewLogger(logger kitLog.Logger, opts Options) *Logger {
	var lev Level

	logger = kitLog.NewSyncLogger(logger)
	logger = kitLog.With(logger,
		"caller", kitLog.Caller(6), // 6 is the caller depth
		"timestamp", kitLog.DefaultTimestampUTC,
	)

	if opts.Name != "" {
		logger = kitLog.With(logger, "logger", opts.Name)
	}

	if opts.Environment != "" {
		logger = kitLog.With(logger, "environment", opts.Environment)
	}

	if opts.Region != "" {
		logger = kitLog.With(logger, "region", opts.Region)
	}

	switch strings.ToLower(opts.Level) {
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
	default:
		lev = Info
		logger = kitLevel.NewFilter(logger, kitLevel.AllowInfo())
	}

	return &Logger{
		Level:  lev,
		Logger: logger,
	}
}

// NewJSONLogger creates a new logger logging in JSON
func NewJSONLogger(opts Options) *Logger {
	logger := kitLog.NewJSONLogger(os.Stdout)
	return NewLogger(logger, opts)
}

// NewFmtLogger creates a new logger logging using fmt format strings
func NewFmtLogger(opts Options) *Logger {
	logger := kitLog.NewLogfmtLogger(os.Stdout)
	return NewLogger(logger, opts)
}

// With returns a new logger which always logs a set of key-value pairs
func (l *Logger) With(kv ...interface{}) *Logger {
	return &Logger{
		Level:  l.Level,
		Logger: kitLog.With(l.Logger, kv...),
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
