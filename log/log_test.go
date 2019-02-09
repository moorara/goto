package log

import (
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
			Options{
				Level:       "",
				Name:        "app",
				Environment: "test",
				Region:      "local",
			},
			Info,
		},
		{
			Options{
				Level:       "debug",
				Name:        "app",
				Environment: "dev",
				Region:      "us-east-1",
			},
			Debug,
		},
		{
			Options{
				Level:       "info",
				Name:        "app",
				Environment: "stage",
				Region:      "us-east-1",
			},
			Info,
		},
		{
			Options{
				Level:       "warn",
				Name:        "app",
				Environment: "prod",
				Region:      "us-east-1",
			},
			Warn,
		},
		{
			Options{
				Level:       "error",
				Name:        "app",
				Environment: "prod",
				Region:      "us-east-1",
			},
			Error,
		},
		{
			Options{
				Level:       "none",
				Name:        "app",
				Environment: "test",
				Region:      "local",
			},
			None,
		},
	}

	for _, tc := range tests {
		l := new(mockLogger)
		logger := NewLogger(l, tc.opts)
		assert.NotNil(t, logger)
		assert.Equal(t, logger.Level, tc.expectedLevel)
	}
}

func TestNewJSONLogger(t *testing.T) {
	tests := []struct {
		opts          Options
		expectedLevel Level
	}{
		{
			Options{
				Level:       "",
				Name:        "app",
				Environment: "test",
				Region:      "local",
			},
			Info,
		},
		{
			Options{
				Level:       "debug",
				Name:        "app",
				Environment: "dev",
				Region:      "us-east-1",
			},
			Debug,
		},
		{
			Options{
				Level:       "info",
				Name:        "app",
				Environment: "stage",
				Region:      "us-east-1",
			},
			Info,
		},
		{
			Options{
				Level:       "warn",
				Name:        "app",
				Environment: "prod",
				Region:      "us-east-1",
			},
			Warn,
		},
		{
			Options{
				Level:       "error",
				Name:        "app",
				Environment: "prod",
				Region:      "us-east-1",
			},
			Error,
		},
		{
			Options{
				Level:       "none",
				Name:        "app",
				Environment: "test",
				Region:      "local",
			},
			None,
		},
	}

	for _, tc := range tests {
		logger := NewJSONLogger(tc.opts)
		assert.NotNil(t, logger)
		assert.Equal(t, logger.Level, tc.expectedLevel)
	}
}

func TestNewFmtLogger(t *testing.T) {
	tests := []struct {
		opts          Options
		expectedLevel Level
	}{
		{
			Options{
				Level:       "",
				Name:        "app",
				Environment: "test",
				Region:      "local",
			},
			Info,
		},
		{
			Options{
				Level:       "debug",
				Name:        "app",
				Environment: "dev",
				Region:      "us-east-1",
			},
			Debug,
		},
		{
			Options{
				Level:       "info",
				Name:        "app",
				Environment: "stage",
				Region:      "us-east-1",
			},
			Info,
		},
		{
			Options{
				Level:       "warn",
				Name:        "app",
				Environment: "prod",
				Region:      "us-east-1",
			},
			Warn,
		},
		{
			Options{
				Level:       "error",
				Name:        "app",
				Environment: "prod",
				Region:      "us-east-1",
			},
			Error,
		},
		{
			Options{
				Level:       "none",
				Name:        "app",
				Environment: "test",
				Region:      "local",
			},
			None,
		},
	}

	for _, tc := range tests {
		logger := NewFmtLogger(tc.opts)
		assert.NotNil(t, logger)
		assert.Equal(t, logger.Level, tc.expectedLevel)
	}
}

func TestWith(t *testing.T) {
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

func TestLog(t *testing.T) {
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
			t.Run("DebugLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Debug(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("InfoLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Info(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("WarnLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Warn(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})

			t.Run("ErrorLevel", func(t *testing.T) {
				logger := &Logger{Logger: &tc.mockLogger}
				err := logger.Error(tc.kv...)
				assert.Equal(t, tc.expectedError, err)
				for _, val := range tc.expectedKV {
					assert.Contains(t, tc.mockLogger.LogInKV, val)
				}
			})
		})
	}
}
