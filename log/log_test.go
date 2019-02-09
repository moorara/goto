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

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel Level
	}{
		{"app", "debug", Debug},
		{"app", "info", Info},
		{"app", "warn", Warn},
		{"app", "error", Error},
		{"app", "none", None},
	}

	for _, tc := range tests {
		l := new(mockLogger)
		logger := NewLogger(l, tc.name, tc.level)
		assert.NotNil(t, logger)
		assert.Equal(t, logger.Name, tc.name)
		assert.Equal(t, logger.Level, tc.expectedLevel)
	}
}

func TestNewNopLogger(t *testing.T) {
	logger := NewNopLogger()
	assert.NotNil(t, logger)
}

func TestNewJSONLogger(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel Level
	}{
		{"app", "debug", Debug},
		{"app", "info", Info},
		{"app", "warn", Warn},
		{"app", "error", Error},
		{"app", "none", None},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger := NewJSONLogger(tc.name, tc.level)
			assert.NotNil(t, logger)
			assert.Equal(t, logger.Name, tc.name)
			assert.Equal(t, logger.Level, tc.expectedLevel)
		})
	}
}

func TestNewFmtLogger(t *testing.T) {
	tests := []struct {
		name          string
		level         string
		expectedLevel Level
	}{
		{"app", "debug", Debug},
		{"app", "info", Info},
		{"app", "warn", Warn},
		{"app", "error", Error},
		{"app", "none", None},
	}

	for _, tc := range tests {
		logger := NewFmtLogger(tc.name, tc.level)
		assert.NotNil(t, logger)
		assert.Equal(t, logger.Name, tc.name)
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

func TestSyncLogger(t *testing.T) {
	tests := []struct {
		mockLogger mockLogger
	}{
		{
			mockLogger{},
		},
	}

	for _, tc := range tests {
		logger := &Logger{Logger: &tc.mockLogger}
		logger = logger.SyncLogger()
		assert.NotNil(t, logger)
	}
}

func TestLogger(t *testing.T) {
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
