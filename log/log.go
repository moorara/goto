// Package log can be used for structured (JSON) logging
//
// Deprecated: this package has been frozen and deprecated in favor of github.com/moorara/observe/log
package log

import (
	"io"
	"os"
	"strings"

	kitLog "github.com/go-kit/kit/log"
	kitLevel "github.com/go-kit/kit/log/level"
)

type (
	// Format is the type for output format
	Format int

	// Level is the type for logging level
	Level int

	// Options contains optional options for Logger
	Options struct {
		depth       int
		Writer      io.Writer
		Format      Format
		Level       string
		Name        string
		Environment string
		Region      string
		Component   string
	}

	// Logger wraps a go-kit Logger
	Logger struct {
		Level  Level
		Logger kitLog.Logger
	}
)

const (
	// JSON represents a json logger
	JSON Format = iota
	// Logfmt represents logfmt logger
	Logfmt
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

var singleton = NewLogger(Options{
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
func NewLogger(opts Options) *Logger {
	logger := &Logger{}
	logger.setOptions(opts)
	return logger
}

func (l *Logger) setOptions(opts Options) {
	var lev Level
	var logger kitLog.Logger

	if opts.depth == 0 {
		opts.depth = 6
	}

	if opts.Writer == nil {
		opts.Writer = os.Stdout
	}

	switch opts.Format {
	case Logfmt:
		logger = kitLog.NewLogfmtLogger(opts.Writer)
	case JSON:
		fallthrough
	default:
		logger = kitLog.NewJSONLogger(opts.Writer)
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

	if opts.Component != "" {
		logger = kitLog.With(logger, "component", opts.Component)
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

	l.Level = lev
	l.Logger = logger
}

// With returns a new logger which always logs a set of key-value pairs
func (l *Logger) With(kv ...interface{}) *Logger {
	return &Logger{
		Level:  l.Level,
		Logger: kitLog.With(l.Logger, kv...),
	}
}

// Debug logs a debug-level event
func (l *Logger) Debug(kv ...interface{}) error {
	return kitLevel.Debug(l.Logger).Log(kv...)
}

// Info logs an info-level event
func (l *Logger) Info(kv ...interface{}) error {
	return kitLevel.Info(l.Logger).Log(kv...)
}

// Warn logs a warn-level event
func (l *Logger) Warn(kv ...interface{}) error {
	return kitLevel.Warn(l.Logger).Log(kv...)
}

// Error logs an error-level event
func (l *Logger) Error(kv ...interface{}) error {
	return kitLevel.Error(l.Logger).Log(kv...)
}

// SetOptions set optional options for singleton logger
func SetOptions(opts Options) {
	opts.depth = 7
	singleton.setOptions(opts)
}

// Debug logs a debug-level event using singleton logger
func Debug(kv ...interface{}) error {
	return singleton.Debug(kv...)
}

// Info logs an info-level event using singleton logger
func Info(kv ...interface{}) error {
	return singleton.Info(kv...)
}

// Warn logs a warn-level event using singleton logger
func Warn(kv ...interface{}) error {
	return singleton.Warn(kv...)
}

// Error logs an error-level event using singleton logger
func Error(kv ...interface{}) error {
	return singleton.Error(kv...)
}
