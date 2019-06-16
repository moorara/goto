package report

import (
	"fmt"
	"net/http"
	"os"

	rollbar "github.com/rollbar/rollbar-go"
)

const (
	// TODO: figure out the right value!
	defaultSkipDepth = 0
)

type (
	rollbarClient interface {
		ErrorWithStackSkip(level string, err error, skip int)
		ErrorWithStackSkipWithExtras(level string, err error, skip int, extras map[string]interface{})
		RequestErrorWithStackSkip(level string, r *http.Request, err error, skip int)
		RequestErrorWithStackSkipWithExtras(level string, r *http.Request, err error, skip int, extras map[string]interface{})
		Wait()
	}

	// RollbarOptions contains optional options for RollbarReporter
	RollbarOptions struct {
		Token       string
		Environment string
		CodeVersion string
		ProjectURL  string
		skipDepth   int
	}

	// RollbarReporter simplifies rollbar client
	RollbarReporter struct {
		client    rollbarClient
		skipDepth int
	}
)

var singleton = &RollbarReporter{}

// NewRollbarReporter creates a new instance of reporter
func NewRollbarReporter(opts RollbarOptions) *RollbarReporter {
	rollbarReporter := &RollbarReporter{}
	rollbarReporter.setOptions(opts)
	return rollbarReporter
}

func (r *RollbarReporter) setOptions(opts RollbarOptions) {
	hostname, _ := os.Hostname()
	client := rollbar.NewAsync(opts.Token, opts.Environment, opts.CodeVersion, hostname, opts.ProjectURL)

	if opts.skipDepth == 0 {
		opts.skipDepth = defaultSkipDepth
	}

	r.client = client
	r.skipDepth = opts.skipDepth
}

// OnPanic reports a panic and should be used with defer
func (r *RollbarReporter) OnPanic() {
	if e := recover(); e != nil {
		err := fmt.Errorf("panic occurred: %v", e)
		r.client.ErrorWithStackSkip(rollbar.CRIT, err, r.skipDepth+2)
		r.client.Wait()
		panic(e)
	}

	r.client.Wait()
}

// Error reports an error
func (r *RollbarReporter) Error(err error) {
	r.client.ErrorWithStackSkip(rollbar.ERR, err, r.skipDepth)
}

// ErrorWithMetadata reports an error with extra metadata
func (r *RollbarReporter) ErrorWithMetadata(err error, metadata map[string]interface{}) {
	r.client.ErrorWithStackSkipWithExtras(rollbar.ERR, err, r.skipDepth, metadata)
}

// HTTPError reports an error for an http request
func (r *RollbarReporter) HTTPError(req *http.Request, err error) {
	r.client.RequestErrorWithStackSkip(rollbar.ERR, req, err, r.skipDepth)
}

// HTTPErrorWithMetadata reports an error for an http request with extra metdata
func (r *RollbarReporter) HTTPErrorWithMetadata(req *http.Request, err error, metadata map[string]interface{}) {
	r.client.RequestErrorWithStackSkipWithExtras(rollbar.ERR, req, err, r.skipDepth, metadata)
}

// Wait blocks until all errors are reported
func (r *RollbarReporter) Wait() {
	r.client.Wait()
}

// SetOptions sets options for singleton reporter
func SetOptions(opts RollbarOptions) {
	singleton.setOptions(opts)
}

// OnPanic reports a panic and should be used with defer
func OnPanic() {
	if singleton.client != nil {
		if e := recover(); e != nil {
			err := fmt.Errorf("panic occurred: %v", e)
			singleton.client.ErrorWithStackSkip(rollbar.CRIT, err, singleton.skipDepth+2)
			singleton.client.Wait()
			panic(e)
		}
	}

	singleton.client.Wait()
}

// Error reports an error
func Error(err error) {
	if singleton.client != nil {
		singleton.client.ErrorWithStackSkip(rollbar.ERR, err, singleton.skipDepth)
	}
}

// ErrorWithMetadata reports an error with extra metadata
func ErrorWithMetadata(err error, metadata map[string]interface{}) {
	if singleton.client != nil {
		singleton.client.ErrorWithStackSkipWithExtras(rollbar.ERR, err, singleton.skipDepth, metadata)
	}
}

// HTTPError reports an error for an http request
func HTTPError(req *http.Request, err error) {
	if singleton.client != nil {
		singleton.client.RequestErrorWithStackSkip(rollbar.ERR, req, err, singleton.skipDepth)
	}
}

// HTTPErrorWithMetadata reports an error for an http request with extra metdata
func HTTPErrorWithMetadata(req *http.Request, err error, metadata map[string]interface{}) {
	if singleton.client != nil {
		singleton.client.RequestErrorWithStackSkipWithExtras(rollbar.ERR, req, err, singleton.skipDepth, metadata)
	}
}

// Wait blocks until all errors are reported
func Wait() {
	if singleton.client != nil {
		singleton.client.Wait()
	}
}
