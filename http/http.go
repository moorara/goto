// Package http provides utility for http servers and clients.
//
// Deprecated: this package has been frozen and deprecated in favor of github.com/moorara/observe/xhttp
package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// contextKey is the type for the keys added to context
type contextKey string

const (
	requestIDHeader     = "Request-Id"
	requestIDContextKey = contextKey("RequestID")
)

// Error is an http error
type Error struct {
	Request    *http.Request
	StatusCode int
	Message    string
}

// NewError creates a new instance of Error
func NewError(res *http.Response) *Error {
	err := &Error{
		Request:    res.Request,
		StatusCode: res.StatusCode,
	}

	if res.Body != nil {
		if data, e := ioutil.ReadAll(res.Body); e == nil {
			err.Message = string(data)
		}
	}

	return err
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %s %d: %s", e.Request.Method, e.Request.URL.Path, e.StatusCode, e.Message)
}

// ResponseWriter extends the functionality of standard http.ResponseWriter
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	StatusClass string
}

// NewResponseWriter creates a new response writer
func NewResponseWriter(rw http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
	}
}

// WriteHeader overrides the default implementation of http.WriteHeader
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)

	// Only capture the first value
	if r.StatusCode == 0 {
		r.StatusCode = statusCode
		r.StatusClass = fmt.Sprintf("%dxx", statusCode/100)
	}
}
