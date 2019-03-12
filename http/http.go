package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type (
	// Error is an http request error
	Error struct {
		Request    *http.Request
		StatusCode int
		Body       string
	}

	// ResponseWriter extends the functionality of standard http.ResponseWriter
	ResponseWriter struct {
		http.ResponseWriter
		StatusCode  int
		StatusClass string
	}
)

// NewError creates a new instance of Error
func NewError(res *http.Response) *Error {
	err := &Error{
		Request:    res.Request,
		StatusCode: res.StatusCode,
	}

	if res.Body != nil {
		if data, e := ioutil.ReadAll(res.Body); e == nil {
			err.Body = string(data)
		}
	}

	return err
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %s %d: %s", e.Request.Method, e.Request.URL.Path, e.StatusCode, e.Body)
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

	if r.StatusCode == 0 {
		r.StatusCode = statusCode
		r.StatusClass = fmt.Sprintf("%dxx", statusCode/100)
	}
}
