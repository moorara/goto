package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/prometheus/client_golang/prometheus"
	promModel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestLoggerForRequest(t *testing.T) {
	tests := []struct {
		name       string
		logger     *log.Logger
		expectedOK bool
	}{
		{
			name:       "WithoutLogger",
			logger:     nil,
			expectedOK: false,
		},
		{
			name:       "WithLogger",
			logger:     log.NewNopLogger(),
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			assert.NoError(t, err)

			if tc.logger != nil {
				ctx := context.WithValue(req.Context(), loggerContextKey, tc.logger)
				req = req.WithContext(ctx)
			}

			logger, ok := LoggerForRequest(req)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.logger, logger)
		})
	}
}

func TestNewServerMiddleware(t *testing.T) {
	tests := []struct {
		name   string
		logger *log.Logger
		mf     *metrics.Factory
		tracer opentracing.Tracer
	}{
		{
			name:   "WithMocks",
			logger: log.NewNopLogger(),
			mf:     metrics.NewFactory(metrics.FactoryOptions{}),
			tracer: mocktracer.New(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewServerMiddleware(tc.logger, tc.mf, tc.tracer)

			assert.Equal(t, tc.logger, m.logger)
			assert.NotNil(t, m.metrics)
			assert.Equal(t, tc.tracer, m.tracer)
		})
	}
}

func TestServerMiddlewareRequestID(t *testing.T) {
	tests := []struct {
		name          string
		req           *http.Request
		requestID     string
		resStatusCode int
	}{
		{
			name:          "WithoutRequestID",
			req:           httptest.NewRequest("GET", "/v1/items", nil),
			requestID:     "",
			resStatusCode: 200,
		},
		{
			name:          "WithRequestID",
			req:           httptest.NewRequest("GET", "/v1/items", nil),
			requestID:     "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			resStatusCode: 200,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var requestIDFromHeader string
			var requestIDFromContext string
			var requestIDFromResponse string

			mid := &ServerMiddleware{}

			if tc.requestID != "" {
				tc.req.Header.Set(requestIDHeader, tc.requestID)
			}

			// Test http handler
			handler := mid.RequestID(func(w http.ResponseWriter, r *http.Request) {
				requestIDFromHeader = r.Header.Get(requestIDHeader)
				requestIDFromContext, _ = r.Context().Value(requestIDContextKey).(string)
				w.WriteHeader(tc.resStatusCode)
			})

			// Handle the mock request
			rec := httptest.NewRecorder()
			handler(rec, tc.req)

			res := rec.Result()
			assert.Equal(t, tc.resStatusCode, res.StatusCode)

			// Verify request id
			requestIDFromResponse = res.Header.Get(requestIDHeader)
			if tc.requestID == "" {
				assert.NotEmpty(t, requestIDFromHeader)
				assert.NotEmpty(t, requestIDFromContext)
				assert.NotEmpty(t, requestIDFromResponse)
			} else {
				assert.Equal(t, tc.requestID, requestIDFromHeader)
				assert.Equal(t, tc.requestID, requestIDFromContext)
				assert.Equal(t, tc.requestID, requestIDFromResponse)
			}
		})
	}
}

func TestServerMiddlewareLogging(t *testing.T) {
	tests := []struct {
		name                string
		req                 *http.Request
		requestID           string
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
			req:                 httptest.NewRequest("GET", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       200,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/items",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
		{
			name:                "301",
			req:                 httptest.NewRequest("GET", "/v1/items/1234", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       301,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/items/1234",
			expectedStatusCode:  301,
			expectedStatusClass: "3xx",
		},
		{
			name:                "404",
			req:                 httptest.NewRequest("POST", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       404,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "POST",
			expectedURL:         "/v1/items",
			expectedStatusCode:  404,
			expectedStatusClass: "4xx",
		},
		{
			name:                "500",
			req:                 httptest.NewRequest("PUT", "/v1/items/1234", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       500,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "PUT",
			expectedURL:         "/v1/items/1234",
			expectedStatusCode:  500,
			expectedStatusClass: "5xx",
		},
		{
			name:                "WithRequestID",
			req:                 httptest.NewRequest("GET", "/v1/items", nil),
			requestID:           "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			resDelay:            10 * time.Millisecond,
			resStatusCode:       200,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/items",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buff := &bytes.Buffer{}
			logger := log.NewLogger(log.Options{Writer: buff})
			mid := &ServerMiddleware{logger: logger}

			if tc.requestID != "" {
				tc.req.Header.Set(requestIDHeader, tc.requestID)
			}

			// Test http handler
			handler := mid.Logging(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.resDelay)
				w.WriteHeader(tc.resStatusCode)
			})

			// Handle the mock request
			rec := httptest.NewRecorder()
			handler(rec, tc.req)

			res := rec.Result()
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)

			// Verify logs

			var log map[string]interface{}
			err := json.NewDecoder(buff).Decode(&log)
			assert.NoError(t, err)
			assert.Equal(t, serverKind, log["http.kind"])
			assert.Equal(t, tc.expectedProto, log["req.proto"])
			assert.Equal(t, tc.expectedMethod, log["req.method"])
			assert.Equal(t, tc.expectedURL, log["req.url"])
			assert.Equal(t, float64(tc.expectedStatusCode), log["res.statusCode"])
			assert.Equal(t, tc.expectedStatusClass, log["res.statusClass"])
			assert.NotEmpty(t, log["responseTime"])
			assert.NotEmpty(t, log["message"])

			if tc.requestID != "" {
				assert.Equal(t, tc.requestID, log["requestId"])
			}
		})
	}
}

func TestServerMiddlewareMetrics(t *testing.T) {
	tests := []struct {
		name                string
		req                 *http.Request
		resDelay            time.Duration
		resStatusCode       int
		expectedMethod      string
		expectedURL         string
		expectedStatusCode  int
		expectedStatusClass string
	}{
		{
			name:                "200",
			req:                 httptest.NewRequest("GET", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       200,
			expectedMethod:      "GET",
			expectedURL:         "/v1/items",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
		{
			name:                "301",
			req:                 httptest.NewRequest("GET", "/v1/items/1234", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       301,
			expectedMethod:      "GET",
			expectedURL:         "/v1/items/1234",
			expectedStatusCode:  301,
			expectedStatusClass: "3xx",
		},
		{
			name:                "404",
			req:                 httptest.NewRequest("POST", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       404,
			expectedMethod:      "POST",
			expectedURL:         "/v1/items",
			expectedStatusCode:  404,
			expectedStatusClass: "4xx",
		},
		{
			name:                "500",
			req:                 httptest.NewRequest("PUT", "/v1/items/1234", nil),
			resDelay:            10 * time.Millisecond,
			resStatusCode:       500,
			expectedMethod:      "PUT",
			expectedURL:         "/v1/items/1234",
			expectedStatusCode:  500,
			expectedStatusClass: "5xx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			promReg := prometheus.NewRegistry()
			metricsFactory := metrics.NewFactory(metrics.FactoryOptions{Registerer: promReg})
			mid := NewServerMiddleware(nil, metricsFactory, nil)
			assert.NotNil(t, mid)

			// Test http handler
			handler := mid.Metrics(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.resDelay)
				w.WriteHeader(tc.resStatusCode)
			})

			// Handle the mock request
			rec := httptest.NewRecorder()
			handler(rec, tc.req)

			res := rec.Result()
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)

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
				case serverGaugeMetricName:
					assert.Equal(t, promModel.MetricType_GAUGE, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case serverCounterMetricName:
					assert.Equal(t, promModel.MetricType_COUNTER, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case serverHistogramMetricName:
					assert.Equal(t, promModel.MetricType_HISTOGRAM, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case serverSummaryMetricName:
					assert.Equal(t, promModel.MetricType_SUMMARY, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				}
			}
		})
	}
}

func TestServerMiddlewareTracing(t *testing.T) {
	tests := []struct {
		name               string
		req                *http.Request
		reqSpan            opentracing.Span
		resDelay           time.Duration
		resStatusCode      int
		expectedProto      string
		expectedMethod     string
		expectedURL        string
		expectedStatusCode int
	}{
		{
			name:               "200",
			req:                httptest.NewRequest("GET", "/v1/items", nil),
			reqSpan:            nil,
			resDelay:           10 * time.Millisecond,
			resStatusCode:      200,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "GET",
			expectedURL:        "/v1/items",
			expectedStatusCode: 200,
		},
		{
			name:               "301",
			req:                httptest.NewRequest("GET", "/v1/items/1234", nil),
			reqSpan:            nil,
			resDelay:           10 * time.Millisecond,
			resStatusCode:      301,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "GET",
			expectedURL:        "/v1/items/1234",
			expectedStatusCode: 301,
		},
		{
			name:               "404",
			req:                httptest.NewRequest("POST", "/v1/items", nil),
			reqSpan:            nil,
			resDelay:           10 * time.Millisecond,
			resStatusCode:      404,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "POST",
			expectedURL:        "/v1/items",
			expectedStatusCode: 404,
		},
		{
			name:               "500",
			req:                httptest.NewRequest("PUT", "/v1/items/1234", nil),
			reqSpan:            nil,
			resDelay:           10 * time.Millisecond,
			resStatusCode:      500,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "PUT",
			expectedURL:        "/v1/items/1234",
			expectedStatusCode: 500,
		},
		{
			name:               "WithRequestSpan",
			req:                httptest.NewRequest("DELETE", "/v1/items/1234", nil),
			reqSpan:            mocktracer.New().StartSpan("parent-span"),
			resDelay:           10 * time.Millisecond,
			resStatusCode:      204,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "DELETE",
			expectedURL:        "/v1/items/1234",
			expectedStatusCode: 204,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var insertedSpan opentracing.Span

			tracer := mocktracer.New()
			mid := &ServerMiddleware{tracer: tracer}

			// Inject the parent span context if any
			if tc.reqSpan != nil {
				carrier := opentracing.HTTPHeadersCarrier(tc.req.Header)
				err := tracer.Inject(tc.reqSpan.Context(), opentracing.HTTPHeaders, carrier)
				assert.NoError(t, err)
			}

			// Test http handler
			handler := mid.Tracing(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.resDelay)
				insertedSpan = opentracing.SpanFromContext(r.Context())
				w.WriteHeader(tc.resStatusCode)
			})

			// Handle the mock request
			rec := httptest.NewRecorder()
			handler(rec, tc.req)

			res := rec.Result()
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)

			// Verify traces

			span := tracer.FinishedSpans()[0]
			assert.Equal(t, insertedSpan, span)
			assert.Equal(t, serverSpanName, span.OperationName)
			assert.Equal(t, tc.expectedProto, span.Tag("http.proto"))
			assert.Equal(t, tc.expectedMethod, span.Tag("http.method"))
			assert.Equal(t, tc.expectedURL, span.Tag("http.url"))
			assert.Equal(t, uint16(tc.expectedStatusCode), span.Tag("http.status_code"))

			if tc.reqSpan != nil {
				reqSpan, ok := tc.reqSpan.(*mocktracer.MockSpan)
				assert.True(t, ok)
				assert.Equal(t, reqSpan.SpanContext.SpanID, span.ParentID)
				assert.Equal(t, reqSpan.SpanContext.TraceID, span.SpanContext.TraceID)
			}

			/* spanLogs := span.Logs()
			spanLogFields := spanLogs[0].Fields
			for _, lf := range spanLogFields {
				switch lf.Key {
				}
			} */
		})
	}
}
