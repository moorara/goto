package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/prometheus/client_golang/prometheus"
	promModel "github.com/prometheus/client_model/go"
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

func TestNewMiddleware(t *testing.T) {
	logger := log.NewLogger(log.Options{
		Level:       "info",
		Name:        "logger",
		Environment: "test",
	})

	promReg := prometheus.NewRegistry()
	mFac := metrics.NewFactory(metrics.FactoryOptions{Registerer: promReg})

	tracer, closer, _ := trace.NewTracer(trace.Options{})
	defer closer.Close()

	tests := []struct {
		name   string
		logger *log.Logger
		mf     *metrics.Factory
		tracer opentracing.Tracer
	}{
		{
			"Default",
			logger,
			mFac,
			tracer,
		},
		{
			"WithMocks",
			log.NewNopLogger(),
			metrics.NewFactory(metrics.FactoryOptions{}),
			mocktracer.New(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMiddleware(tc.logger, tc.mf, tc.tracer)

			assert.Equal(t, tc.logger, m.logger)
			assert.NotNil(t, m.metrics)
			assert.Equal(t, tc.tracer, m.tracer)
		})
	}
}

func TestMiddlewareWrap(t *testing.T) {
	tests := []struct {
		name                string
		req                 *http.Request
		reqSpan             opentracing.Span
		resDelay            time.Duration
		resStatusCode       int
		expectedProto       string
		expectedMethod      string
		expectedURL         string
		expectedStatusCode  int
		expectedStatusClass string
	}{
		{
			name:                "200",
			req:                 httptest.NewRequest("GET", "/v1/dogs/breeds", nil),
			reqSpan:             nil,
			resDelay:            10 * time.Millisecond,
			resStatusCode:       200,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/dogs/breeds",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
		{
			name:                "301",
			req:                 httptest.NewRequest("GET", "/v1/dogs/breeds/1234", nil),
			reqSpan:             nil,
			resDelay:            10 * time.Millisecond,
			resStatusCode:       301,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/dogs/breeds/1234",
			expectedStatusCode:  301,
			expectedStatusClass: "3xx",
		},
		{
			name:                "404",
			req:                 httptest.NewRequest("POST", "/v1/breeds/dogs", nil),
			reqSpan:             nil,
			resDelay:            10 * time.Millisecond,
			resStatusCode:       404,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "POST",
			expectedURL:         "/v1/breeds/dogs",
			expectedStatusCode:  404,
			expectedStatusClass: "4xx",
		},
		{
			name:                "500",
			req:                 httptest.NewRequest("PUT", "/v1/dogs/breeds/abcd", nil),
			reqSpan:             nil,
			resDelay:            10 * time.Millisecond,
			resStatusCode:       500,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "PUT",
			expectedURL:         "/v1/dogs/breeds/abcd",
			expectedStatusCode:  500,
			expectedStatusClass: "5xx",
		},
		{
			name:                "WithRequestSpan",
			req:                 httptest.NewRequest("DELETE", "/v1/dogs/breeds/1234-abcd", nil),
			reqSpan:             mocktracer.New().StartSpan("parent-span"),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       204,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "DELETE",
			expectedURL:         "/v1/dogs/breeds/1234-abcd",
			expectedStatusCode:  204,
			expectedStatusClass: "2xx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var reqSpan opentracing.Span
			buff := bytes.Buffer{}

			// Create logger
			logger := log.NewLogger(log.Options{
				Writer:      &buff,
				Format:      log.JSON,
				Level:       "info",
				Name:        "logger",
				Environment: "test",
				Region:      "local",
				Component:   "http",
			})

			// Create metrics factory
			promReg := prometheus.NewRegistry()
			metricsFactory := metrics.NewFactory(metrics.FactoryOptions{
				Registerer: promReg,
			})

			// Create tracer
			tracer := mocktracer.New()

			// Inject parent span context
			if tc.reqSpan != nil {
				carrier := opentracing.HTTPHeadersCarrier(tc.req.Header)
				err := tracer.Inject(tc.reqSpan.Context(), opentracing.HTTPHeaders, carrier)
				assert.NoError(t, err)
			}

			// Create http middleware
			mid := NewMiddleware(logger, metricsFactory, tracer)
			assert.NotNil(t, mid)

			// Test http handler
			handler := mid.Wrap(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.resDelay)
				reqSpan = opentracing.SpanFromContext(r.Context())
				w.WriteHeader(tc.resStatusCode)
			})

			// Trigger a mock request
			rec := httptest.NewRecorder()
			handler(rec, tc.req)

			res := rec.Result()
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)

			// Verify logs

			var log map[string]interface{}
			err := json.NewDecoder(&buff).Decode(&log)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedProto, log["req.proto"])
			assert.Equal(t, tc.expectedMethod, log["req.method"])
			assert.Equal(t, tc.expectedURL, log["req.url"])
			assert.Equal(t, float64(tc.expectedStatusCode), log["res.statusCode"])
			assert.Equal(t, tc.expectedStatusClass, log["res.statusClass"])
			assert.NotEmpty(t, log["responseTime"])
			assert.NotEmpty(t, log["message"])

			// Verify metrics

			verifyLabels := func(labels []*promModel.LabelPair) {
				for _, l := range labels {
					switch *l.Name {
					case "method":
						assert.Equal(t, tc.expectedMethod, *l.Value)
					case "url":
						assert.Equal(t, tc.expectedURL, *l.Value)
					case "statusCode":
						assert.Equal(t, strconv.Itoa(tc.expectedStatusCode), *l.Value)
					case "statusClass":
						assert.Equal(t, tc.expectedStatusClass, *l.Value)
					}
				}
			}

			metricFamilies, err := promReg.Gather()
			assert.NoError(t, err)

			for _, metricFamily := range metricFamilies {
				switch *metricFamily.Name {
				case gaugeMetricName:
					assert.Equal(t, promModel.MetricType_GAUGE, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case counterMetricName:
					assert.Equal(t, promModel.MetricType_COUNTER, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case histogramMetricName:
					assert.Equal(t, promModel.MetricType_HISTOGRAM, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case summaryMetricName:
					assert.Equal(t, promModel.MetricType_SUMMARY, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				}
			}

			// Verify traces

			span := tracer.FinishedSpans()[0]
			assert.Equal(t, reqSpan, span)
			assert.Equal(t, "http-request", span.OperationName)
			assert.Equal(t, tc.expectedProto, span.Tag("http.proto"))
			assert.Equal(t, tc.expectedMethod, span.Tag("http.method"))
			assert.Equal(t, tc.expectedURL, span.Tag("http.url"))
			assert.Equal(t, uint16(tc.expectedStatusCode), span.Tag("http.status_code"))

			if tc.reqSpan != nil {
				parentSpan, ok := tc.reqSpan.(*mocktracer.MockSpan)
				assert.True(t, ok)
				assert.Equal(t, parentSpan.SpanContext.SpanID, span.ParentID)
				assert.Equal(t, parentSpan.SpanContext.TraceID, span.SpanContext.TraceID)
			}
		})
	}
}
