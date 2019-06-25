package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	tests := []struct {
		name          string
		request       *http.Request
		statusCode    int
		body          string
		expectedError string
	}{
		{
			"400",
			&http.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/",
				},
			},
			http.StatusBadRequest,
			"Invalid request",
			"GET / 400: Invalid request",
		},
		{
			"500",
			&http.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/",
				},
			},
			http.StatusInternalServerError,
			"Internal error",
			"POST / 500: Internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			br := strings.NewReader(tc.body)
			rc := ioutil.NopCloser(br)

			res := &http.Response{
				Request:    tc.request,
				StatusCode: tc.statusCode,
				Body:       rc,
			}

			err := NewError(res)
			assert.Equal(t, tc.request, err.Request)
			assert.Equal(t, tc.statusCode, err.StatusCode)
			assert.Equal(t, tc.body, err.Message)

			var e error = err
			assert.Equal(t, tc.expectedError, e.Error())
		})
	}
}

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		statusClass string
	}{
		{"200", 101, "1xx"},
		{"200", 200, "2xx"},
		{"201", 201, "2xx"},
		{"201", 202, "2xx"},
		{"300", 300, "3xx"},
		{"400", 400, "4xx"},
		{"400", 403, "4xx"},
		{"404", 404, "4xx"},
		{"400", 409, "4xx"},
		{"500", 500, "5xx"},
		{"500", 501, "5xx"},
		{"502", 502, "5xx"},
		{"500", 503, "5xx"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}

			middleware := func(w http.ResponseWriter, r *http.Request) {
				rw := NewResponseWriter(w)
				handler(rw, r)

				assert.Equal(t, tc.statusCode, rw.StatusCode)
				assert.Equal(t, tc.statusClass, rw.StatusClass)
			}

			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()
			middleware(w, r)
		})
	}
}
