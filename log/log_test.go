package log

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockLogger struct {
	LogInKV     []interface{}
	LogOutError error
}

func (m *mockLogger) Log(kv ...interface{}) error {
	m.LogInKV = kv
	return m.LogOutError
}

func TestNewNopLogger(t *testing.T) {
	logger := NewNopLogger()
	assert.NotNil(t, logger)
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		opts          Options
		expectedLevel Level
	}{
		{
			Options{},
			InfoLevel,
		},
	}

	for _, tc := range tests {
		logger := NewLogger(tc.opts)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.Logger)
		assert.Equal(t, logger.Level, tc.expectedLevel)
	}
}

func TestLoggerSetOptions(t *testing.T) {
	tests := []struct {
		name          string
		opts          Options
		expectedLevel Level
	}{
		{
			"NoLevel",
			Options{
				Name:        "instance",
				Environment: "test",
				Region:      "local",
				Component:   "app",
			},
			InfoLevel,
		},
		{
			"DebugLevel",
			Options{
				Format:      Logfmt,
				Level:       "debug",
				Name:        "instance",
				Environment: "dev",
				Region:      "us-east-1",
				Component:   "app",
			},
			DebugLevel,
		},
		{
			"InfoLevel",
			Options{
				Format:      JSON,
				Level:       "info",
				Name:        "instance",
				Environment: "stage",
				Region:      "us-east-1",
				Component:   "app",
			},
			InfoLevel,
		},
		{
			"WarnLevel",
			Options{
				Format:      JSON,
				Level:       "warn",
				Name:        "instance",
				Environment: "prod",
				Region:      "us-east-1",
				Component:   "app",
			},
			WarnLevel,
		},
		{
			"ErrorLevel",
			Options{
				Format:      JSON,
				Level:       "error",
				Name:        "instance",
				Environment: "prod",
				Region:      "us-east-1",
				Component:   "app",
			},
			ErrorLevel,
		},
		{
			"NoneLevel",
			Options{
				Level:       "none",
				Name:        "instance",
				Environment: "test",
				Region:      "local",
				Component:   "app",
			},
			NoneLevel,
		},
		{
			"CustomWriter",
			Options{
				Writer:      &bytes.Buffer{},
				Name:        "instance",
				Environment: "test",
				Region:      "local",
				Component:   "app",
			},
			InfoLevel,
		},
	}

	for _, tc := range tests {
		logger := &Logger{}
		logger.setOptions(tc.opts)

		assert.NotNil(t, logger)
		assert.NotNil(t, logger.Logger)
		assert.Equal(t, logger.Level, tc.expectedLevel)
	}
}

func TestLoggerWith(t *testing.T) {
	tests := []struct {
		mockLogger mockLogger
		kv         []interface{}
	}{
		{
			mockLogger{},
			[]interface{}{"version", "0.1.0", "revision", "1234567", "context", "test"},
		},
	}

	for _, tc := range tests {
		logger := &Logger{Logger: &tc.mockLogger}
		logger = logger.With(tc.kv...)
		assert.NotNil(t, logger)
	}
}

func TestLoggerLog(t *testing.T) {
	tests := []struct {
		name          string
		mockLogger    mockLogger
		kv            []interface{}
		expectedError error
		expectedKV    []interface{}
	}{
		{
			"Error",
			mockLogger{
				LogOutError: errors.New("log error"),
			},
			[]interface{}{"message", "operation failed", "reason", "no capacity"},
			errors.New("log error"),
			[]interface{}{"message", "operation failed", "reason", "no capacity"},
		},
		{
			"Success",
			mockLogger{},
			[]interface{}{"message", "operation succeeded", "region", "home"},
			nil,
			[]interface{}{"message", "operation succeeded", "region", "home"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := &Logger{Logger: &tc.mockLogger}

			t.Run("DebugLevel", func(t *testing.T) {
				err := logger.Debug(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("InfoLevel", func(t *testing.T) {
				err := logger.Info(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("WarnLevel", func(t *testing.T) {
				err := logger.Warn(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("ErrorLevel", func(t *testing.T) {
				err := logger.Error(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})
		})
	}
}

func TestSingletonSetOptions(t *testing.T) {
	tests := []struct {
		opts          Options
		expectedLevel Level
	}{
		{
			Options{
				Level:       "",
				Name:        "instance",
				Environment: "test",
				Region:      "local",
				Component:   "app",
			},
			InfoLevel,
		},
		{
			Options{
				Format:      Logfmt,
				Level:       "debug",
				Name:        "instance",
				Environment: "dev",
				Region:      "us-east-1",
				Component:   "app",
			},
			DebugLevel,
		},
		{
			Options{
				Format:      JSON,
				Level:       "info",
				Name:        "instance",
				Environment: "stage",
				Region:      "us-east-1",
				Component:   "app",
			},
			InfoLevel,
		},
		{
			Options{
				Format:      JSON,
				Level:       "warn",
				Name:        "instance",
				Environment: "prod",
				Region:      "us-east-1",
				Component:   "app",
			},
			WarnLevel,
		},
		{
			Options{
				Format:      JSON,
				Level:       "error",
				Name:        "instance",
				Environment: "prod",
				Region:      "us-east-1",
				Component:   "app",
			},
			ErrorLevel,
		},
		{
			Options{
				Level:       "none",
				Name:        "instance",
				Environment: "test",
				Region:      "local",
				Component:   "app",
			},
			NoneLevel,
		},
	}

	for _, tc := range tests {
		SetOptions(tc.opts)

		assert.NotNil(t, singleton.Logger)
		assert.Equal(t, singleton.Level, tc.expectedLevel)
	}
}

func TestSingletonLog(t *testing.T) {
	tests := []struct {
		name          string
		mockLogger    mockLogger
		kv            []interface{}
		expectedError error
		expectedKV    []interface{}
	}{
		{
			"Error",
			mockLogger{
				LogOutError: errors.New("log error"),
			},
			[]interface{}{"message", "operation failed", "reason", "no capacity"},
			errors.New("log error"),
			[]interface{}{"message", "operation failed", "reason", "no capacity"},
		},
		{
			"Success",
			mockLogger{},
			[]interface{}{"message", "operation succeeded", "region", "home"},
			nil,
			[]interface{}{"message", "operation succeeded", "region", "home"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			singleton = &Logger{Logger: &tc.mockLogger}

			t.Run("DebugLevel", func(t *testing.T) {
				err := Debug(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("InfoLevel", func(t *testing.T) {
				err := Info(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("WarnLevel", func(t *testing.T) {
				err := Warn(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("ErrorLevel", func(t *testing.T) {
				err := Error(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})
		})
	}
}
