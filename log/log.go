package log

import (
	"os"
	"strings"

	kitLog "github.com/go-kit/kit/log"
	kitLevel "github.com/go-kit/kit/log/level"
)

type (
	// Level is the type for log level
	Level int

	// Options contains optional options for Logger
	Options struct {
		depth       int
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

const (
	// DebugLevel log
	DebugLevel Level = iota
	// InfoLevel log
	InfoLevel
	// WarnLevel log
	WarnLevel
	// ErrorLevel log
	ErrorLevel
	// NoneLevel log
	NoneLevel
)

var singleton = NewJSONLogger(Options{
	depth: 7,
	Name:  "singleton",
})

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

	// Set default depth
	if opts.depth == 0 {
		opts.depth = 6
	}

	logger = kitLog.NewSyncLogger(logger)
	logger = kitLog.With(logger,
		"caller", kitLog.Caller(opts.depth),
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
		lev = DebugLevel
		logger = kitLevel.NewFilter(logger, kitLevel.AllowDebug())
	case "info":
		lev = InfoLevel
		logger = kitLevel.NewFilter(logger, kitLevel.AllowInfo())
	case "warn":
		lev = WarnLevel
		logger = kitLevel.NewFilter(logger, kitLevel.AllowWarn())
	case "error":
		lev = ErrorLevel
		logger = kitLevel.NewFilter(logger, kitLevel.AllowError())
	case "none":
		lev = NoneLevel
		logger = kitLevel.NewFilter(logger, kitLevel.AllowNone())
	default:
		lev = InfoLevel
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

// Debug logs in debug level
func Debug(kv ...interface{}) error {
	return singleton.Debug(kv...)
}

// Info logs in info level
func Info(kv ...interface{}) error {
	return singleton.Info(kv...)
}

// Warn logs in warn level
func Warn(kv ...interface{}) error {
	return singleton.Warn(kv...)
}

// Error logs in error level
func Error(kv ...interface{}) error {
	return singleton.Error(kv...)
}
