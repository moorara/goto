package report

import (
	"errors"
	"net/http"
	"testing"

	rollbar "github.com/rollbar/rollbar-go"
	"github.com/stretchr/testify/assert"
)

type mockRollbarClient struct {
	ErrorWithStackSkipInLevel string
	ErrorWithStackSkipInError error
	ErrorWithStackSkipInSkip  int

	ErrorWithStackSkipWithExtrasInLevel  string
	ErrorWithStackSkipWithExtrasInError  error
	ErrorWithStackSkipWithExtrasInSkip   int
	ErrorWithStackSkipWithExtrasInExtras map[string]interface{}

	RequestErrorWithStackSkipInLevel   string
	RequestErrorWithStackSkipInRequest *http.Request
	RequestErrorWithStackSkipInError   error
	RequestErrorWithStackSkipInSkip    int

	RequestErrorWithStackSkipWithExtrasInLevel   string
	RequestErrorWithStackSkipWithExtrasInRequest *http.Request
	RequestErrorWithStackSkipWithExtrasInError   error
	RequestErrorWithStackSkipWithExtrasInSkip    int
	RequestErrorWithStackSkipWithExtrasInExtras  map[string]interface{}

	WaitCalled bool
}

func (m *mockRollbarClient) ErrorWithStackSkip(level string, err error, skip int) {
	m.ErrorWithStackSkipInLevel = level
	m.ErrorWithStackSkipInError = err
	m.ErrorWithStackSkipInSkip = skip
}

func (m *mockRollbarClient) ErrorWithStackSkipWithExtras(level string, err error, skip int, extras map[string]interface{}) {
	m.ErrorWithStackSkipWithExtrasInLevel = level
	m.ErrorWithStackSkipWithExtrasInError = err
	m.ErrorWithStackSkipWithExtrasInSkip = skip
	m.ErrorWithStackSkipWithExtrasInExtras = extras
}

func (m *mockRollbarClient) RequestErrorWithStackSkip(level string, r *http.Request, err error, skip int) {
	m.RequestErrorWithStackSkipInLevel = level
	m.RequestErrorWithStackSkipInRequest = r
	m.RequestErrorWithStackSkipInError = err
	m.RequestErrorWithStackSkipInSkip = skip
}

func (m *mockRollbarClient) RequestErrorWithStackSkipWithExtras(level string, r *http.Request, err error, skip int, extras map[string]interface{}) {
	m.RequestErrorWithStackSkipWithExtrasInLevel = level
	m.RequestErrorWithStackSkipWithExtrasInRequest = r
	m.RequestErrorWithStackSkipWithExtrasInError = err
	m.RequestErrorWithStackSkipWithExtrasInSkip = skip
	m.RequestErrorWithStackSkipWithExtrasInExtras = extras
}

func (m *mockRollbarClient) Wait() {
	m.WaitCalled = true
}

func TestNewRollbarReporter(t *testing.T) {
	tests := []struct {
		opts              RollbarOptions
		expectedSkipDepth int
	}{
		{
			RollbarOptions{},
			defaultSkipDepth,
		},
	}

	for _, tc := range tests {
		reporter := NewRollbarReporter(tc.opts)

		assert.NotNil(t, reporter)
		assert.NotNil(t, reporter.client)
		assert.Equal(t, tc.expectedSkipDepth, reporter.skipDepth)
	}
}

func TestReporterSetOptions(t *testing.T) {
	tests := []struct {
		opts              RollbarOptions
		expectedSkipDepth int
	}{
		{
			RollbarOptions{},
			defaultSkipDepth,
		},
		{
			RollbarOptions{
				Token:       "rollbar-token",
				Environment: "test",
				CodeVersion: "0.1.0",
				ProjectURL:  "github.com/owner/repo",
				skipDepth:   6,
			},
			6,
		},
	}

	for _, tc := range tests {
		reporter := &RollbarReporter{}
		reporter.setOptions(tc.opts)

		assert.NotNil(t, reporter)
		assert.NotNil(t, reporter.client)
		assert.Equal(t, tc.expectedSkipDepth, reporter.skipDepth)
	}
}

func TestReporterError(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		err       error
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			err:       errors.New("error occurred"),
		},
	}

	for _, tc := range tests {
		reporter := &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		reporter.Error(tc.err)

		assert.Equal(t, rollbar.ERR, tc.client.ErrorWithStackSkipInLevel)
		assert.Equal(t, tc.err, tc.client.ErrorWithStackSkipInError)
		assert.Equal(t, tc.skipDepth, tc.client.ErrorWithStackSkipInSkip)
	}
}

func TestReporterErrorWithMetadata(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		err       error
		metadata  map[string]interface{}
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			err:       errors.New("error occurred"),
			metadata: map[string]interface{}{
				"code": 7,
			},
		},
	}

	for _, tc := range tests {
		reporter := &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		reporter.ErrorWithMetadata(tc.err, tc.metadata)

		assert.Equal(t, rollbar.ERR, tc.client.ErrorWithStackSkipWithExtrasInLevel)
		assert.Equal(t, tc.err, tc.client.ErrorWithStackSkipWithExtrasInError)
		assert.Equal(t, tc.skipDepth, tc.client.ErrorWithStackSkipWithExtrasInSkip)
		assert.Equal(t, tc.metadata, tc.client.ErrorWithStackSkipWithExtrasInExtras)
	}
}

func TestReporterHTTPError(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		req       *http.Request
		err       error
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			req:       &http.Request{},
			err:       errors.New("error occurred"),
		},
	}

	for _, tc := range tests {
		reporter := &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		reporter.HTTPError(tc.req, tc.err)

		assert.Equal(t, rollbar.ERR, tc.client.RequestErrorWithStackSkipInLevel)
		assert.Equal(t, tc.req, tc.client.RequestErrorWithStackSkipInRequest)
		assert.Equal(t, tc.err, tc.client.RequestErrorWithStackSkipInError)
		assert.Equal(t, tc.skipDepth, tc.client.RequestErrorWithStackSkipInSkip)
	}
}

func TestReporterHTTPErrorWithMetadata(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		req       *http.Request
		err       error
		metadata  map[string]interface{}
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			req:       &http.Request{},
			err:       errors.New("error occurred"),
			metadata: map[string]interface{}{
				"code": 7,
			},
		},
	}

	for _, tc := range tests {
		reporter := &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		reporter.HTTPErrorWithMetadata(tc.req, tc.err, tc.metadata)

		assert.Equal(t, rollbar.ERR, tc.client.RequestErrorWithStackSkipWithExtrasInLevel)
		assert.Equal(t, tc.req, tc.client.RequestErrorWithStackSkipWithExtrasInRequest)
		assert.Equal(t, tc.err, tc.client.RequestErrorWithStackSkipWithExtrasInError)
		assert.Equal(t, tc.skipDepth, tc.client.RequestErrorWithStackSkipWithExtrasInSkip)
		assert.Equal(t, tc.metadata, tc.client.RequestErrorWithStackSkipWithExtrasInExtras)
	}
}

func TestReporterOnPanic(t *testing.T) {
	tests := []struct {
		name              string
		client            *mockRollbarClient
		skipDepth         int
		panicValue        interface{}
		expectedError     error
		expectedSkipDepth int
	}{
		{
			name:       "Panic",
			client:     &mockRollbarClient{},
			skipDepth:  4,
			panicValue: nil,
		},
		{
			name:              "PanicWithError",
			client:            &mockRollbarClient{},
			skipDepth:         4,
			panicValue:        errors.New("error"),
			expectedError:     errors.New("a panic occurred: error"),
			expectedSkipDepth: 6,
		},
		{
			name:              "PanicWithInt",
			client:            &mockRollbarClient{},
			skipDepth:         4,
			panicValue:        1,
			expectedError:     errors.New("a panic occurred: 1"),
			expectedSkipDepth: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reporter := &RollbarReporter{
				client:    tc.client,
				skipDepth: tc.skipDepth,
			}

			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, rollbar.CRIT, tc.client.ErrorWithStackSkipInLevel)
					assert.Equal(t, tc.expectedError, tc.client.ErrorWithStackSkipInError)
					assert.Equal(t, tc.expectedSkipDepth, tc.client.ErrorWithStackSkipInSkip)
					assert.True(t, tc.client.WaitCalled)
				}
			}()

			defer reporter.OnPanic()
			panic(tc.panicValue)
		})
	}
}

func TestSingletonSetOptions(t *testing.T) {
	tests := []struct {
		opts              RollbarOptions
		expectedSkipDepth int
	}{
		{
			RollbarOptions{},
			defaultSkipDepth,
		},
		{
			RollbarOptions{
				Token:       "rollbar-token",
				Environment: "test",
				CodeVersion: "0.1.0",
				ProjectURL:  "github.com/owner/repo",
				skipDepth:   6,
			},
			6,
		},
	}

	for _, tc := range tests {
		SetOptions(tc.opts)

		assert.NotNil(t, singleton.client)
		assert.Equal(t, tc.expectedSkipDepth, singleton.skipDepth)
	}
}

func TestSingletonError(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		err       error
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			err:       errors.New("error occurred"),
		},
	}

	for _, tc := range tests {
		singleton = &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		Error(tc.err)

		assert.Equal(t, rollbar.ERR, tc.client.ErrorWithStackSkipInLevel)
		assert.Equal(t, tc.err, tc.client.ErrorWithStackSkipInError)
		assert.Equal(t, tc.skipDepth, tc.client.ErrorWithStackSkipInSkip)
	}
}

func TestSingletonErrorWithMetadata(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		err       error
		metadata  map[string]interface{}
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			err:       errors.New("error occurred"),
			metadata: map[string]interface{}{
				"code": 7,
			},
		},
	}

	for _, tc := range tests {
		singleton = &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		ErrorWithMetadata(tc.err, tc.metadata)

		assert.Equal(t, rollbar.ERR, tc.client.ErrorWithStackSkipWithExtrasInLevel)
		assert.Equal(t, tc.err, tc.client.ErrorWithStackSkipWithExtrasInError)
		assert.Equal(t, tc.skipDepth, tc.client.ErrorWithStackSkipWithExtrasInSkip)
		assert.Equal(t, tc.metadata, tc.client.ErrorWithStackSkipWithExtrasInExtras)
	}
}

func TestSingletonHTTPError(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		req       *http.Request
		err       error
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			req:       &http.Request{},
			err:       errors.New("error occurred"),
		},
	}

	for _, tc := range tests {
		singleton = &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		HTTPError(tc.req, tc.err)

		assert.Equal(t, rollbar.ERR, tc.client.RequestErrorWithStackSkipInLevel)
		assert.Equal(t, tc.req, tc.client.RequestErrorWithStackSkipInRequest)
		assert.Equal(t, tc.err, tc.client.RequestErrorWithStackSkipInError)
		assert.Equal(t, tc.skipDepth, tc.client.RequestErrorWithStackSkipInSkip)
	}
}

func TestSingletonHTTPErrorWithMetadata(t *testing.T) {
	tests := []struct {
		client    *mockRollbarClient
		skipDepth int
		req       *http.Request
		err       error
		metadata  map[string]interface{}
	}{
		{
			client:    &mockRollbarClient{},
			skipDepth: 4,
			req:       &http.Request{},
			err:       errors.New("error occurred"),
			metadata: map[string]interface{}{
				"code": 7,
			},
		},
	}

	for _, tc := range tests {
		singleton = &RollbarReporter{
			client:    tc.client,
			skipDepth: tc.skipDepth,
		}

		HTTPErrorWithMetadata(tc.req, tc.err, tc.metadata)

		assert.Equal(t, rollbar.ERR, tc.client.RequestErrorWithStackSkipWithExtrasInLevel)
		assert.Equal(t, tc.req, tc.client.RequestErrorWithStackSkipWithExtrasInRequest)
		assert.Equal(t, tc.err, tc.client.RequestErrorWithStackSkipWithExtrasInError)
		assert.Equal(t, tc.skipDepth, tc.client.RequestErrorWithStackSkipWithExtrasInSkip)
		assert.Equal(t, tc.metadata, tc.client.RequestErrorWithStackSkipWithExtrasInExtras)
	}
}

func TestSingletonOnPanic(t *testing.T) {
	tests := []struct {
		name              string
		client            *mockRollbarClient
		skipDepth         int
		panicValue        interface{}
		expectedError     error
		expectedSkipDepth int
	}{
		{
			name:       "Panic",
			client:     &mockRollbarClient{},
			skipDepth:  4,
			panicValue: nil,
		},
		{
			name:              "PanicWithError",
			client:            &mockRollbarClient{},
			skipDepth:         4,
			panicValue:        errors.New("error"),
			expectedError:     errors.New("a panic occurred: error"),
			expectedSkipDepth: 6,
		},
		{
			name:              "PanicWithInt",
			client:            &mockRollbarClient{},
			skipDepth:         4,
			panicValue:        1,
			expectedError:     errors.New("a panic occurred: 1"),
			expectedSkipDepth: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			singleton = &RollbarReporter{
				client:    tc.client,
				skipDepth: tc.skipDepth,
			}

			defer func() {
				if r := recover(); r != nil {
					assert.Equal(t, rollbar.CRIT, tc.client.ErrorWithStackSkipInLevel)
					assert.Equal(t, tc.expectedError, tc.client.ErrorWithStackSkipInError)
					assert.Equal(t, tc.expectedSkipDepth, tc.client.ErrorWithStackSkipInSkip)
					assert.True(t, tc.client.WaitCalled)
				}
			}()

			defer OnPanic()
			panic(tc.panicValue)
		})
	}
}
