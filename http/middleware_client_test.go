package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

func extractSpanContext(req *http.Request, tracer opentracing.Tracer) opentracing.SpanContext {
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	parentSpanContext, err := tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err == nil {
		return parentSpanContext
	}

	return nil
}

func TestNewClientMiddleware(t *testing.T) {
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
			m := NewClientMiddleware(tc.logger, tc.mf, tc.tracer)

			assert.Equal(t, tc.logger, m.logger)
			assert.NotNil(t, m.metrics)
			assert.Equal(t, tc.tracer, m.tracer)
		})
	}
}

func TestClientMiddlewareRequestID(t *testing.T) {
	tests := []struct {
		name                         string
		req                          *http.Request
		reqHeaders                   map[string][]string
		reqCtx                       context.Context
		resError                     error
		resStatusCode                int
		expectedRequestIDFromHeader  string
		expectedRequestIDFromContext string
	}{
		{
			name:                         "NoRequestID",
			req:                          httptest.NewRequest("GET", "/v1/items", nil),
			reqHeaders:                   map[string][]string{},
			reqCtx:                       context.Background(),
			resError:                     nil,
			resStatusCode:                200,
			expectedRequestIDFromHeader:  "", // expected to be generated
			expectedRequestIDFromContext: "", // expected to be generated
		},
		{
			name:                         "RequestIDInContext",
			req:                          httptest.NewRequest("GET", "/v1/items", nil),
			reqHeaders:                   map[string][]string{},
			reqCtx:                       context.WithValue(context.Background(), requestIDContextKey, "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			resError:                     nil,
			resStatusCode:                200,
			expectedRequestIDFromHeader:  "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			expectedRequestIDFromContext: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		},
		{
			name:                         "RequestIDInHeaders",
			req:                          httptest.NewRequest("GET", "/v1/items", nil),
			reqHeaders:                   map[string][]string{requestIDHeader: []string{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"}},
			reqCtx:                       context.Background(),
			resError:                     nil,
			resStatusCode:                200,
			expectedRequestIDFromHeader:  "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			expectedRequestIDFromContext: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		},
		{
			name:                         "RequestIDInContextAndHeaders",
			req:                          httptest.NewRequest("GET", "/v1/items", nil),
			reqHeaders:                   map[string][]string{requestIDHeader: []string{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"}},
			reqCtx:                       context.WithValue(context.Background(), requestIDContextKey, "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"),
			resError:                     nil,
			resStatusCode:                200,
			expectedRequestIDFromHeader:  "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			expectedRequestIDFromContext: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var requestIDFromHeader string
			var requestIDFromContext string

			mid := &ClientMiddleware{}

			// Add context and headers to request
			tc.req = tc.req.WithContext(tc.reqCtx)
			for k, vals := range tc.reqHeaders {
				for _, v := range vals {
					tc.req.Header.Add(k, v)
				}
			}

			// Test http doer
			doer := mid.RequestID(func(r *http.Request) (*http.Response, error) {
				requestIDFromHeader = r.Header.Get(requestIDHeader)
				requestIDFromContext, _ = r.Context().Value(requestIDContextKey).(string)
				if tc.resError != nil {
					return nil, tc.resError
				}
				return &http.Response{StatusCode: tc.resStatusCode}, nil
			})

			// Make the mock request
			res, err := doer(tc.req)

			if tc.resError != nil {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.resStatusCode, res.StatusCode)
			}

			// Verify request id

			if tc.expectedRequestIDFromHeader == "" {
				assert.NotEmpty(t, requestIDFromHeader)
			} else {
				assert.Equal(t, tc.expectedRequestIDFromHeader, requestIDFromHeader)
			}

			if tc.expectedRequestIDFromContext == "" {
				assert.NotEmpty(t, requestIDFromContext)
			} else {
				assert.Equal(t, tc.expectedRequestIDFromContext, requestIDFromContext)
			}
		})
	}
}

func TestClientMiddlewareLogging(t *testing.T) {
	tests := []struct {
		name                string
		req                 *http.Request
		requestID           string
		resDelay            time.Duration
		resError            error
		resStatusCode       int
		expectedProto       string
		expectedMethod      string
		expectedURL         string
		expectedStatusCode  int
		expectedStatusClass string
	}{
		{
			name:                "Error",
			req:                 httptest.NewRequest("GET", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resError:            errors.New("uknown error"),
			resStatusCode:       0,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/items",
			expectedStatusCode:  -1,
			expectedStatusClass: "",
		},
		{
			name:                "200",
			req:                 httptest.NewRequest("GET", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resError:            nil,
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
			resError:            nil,
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
			resError:            nil,
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
			resError:            nil,
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
			resError:            nil,
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
			mid := &ClientMiddleware{logger: logger}

			if tc.requestID != "" {
				ctx := tc.req.Context()
				ctx = context.WithValue(ctx, requestIDContextKey, tc.requestID)
				tc.req = tc.req.WithContext(ctx)
			}

			// Test http doer
			doer := mid.Logging(func(r *http.Request) (*http.Response, error) {
				time.Sleep(tc.resDelay)
				if tc.resError != nil {
					return nil, tc.resError
				}
				return &http.Response{StatusCode: tc.resStatusCode}, nil
			})

			// Make the mock request
			res, err := doer(tc.req)

			if tc.resError != nil {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.resStatusCode, res.StatusCode)
			}

			// Verify logs

			var log map[string]interface{}
			err = json.NewDecoder(buff).Decode(&log)
			assert.NoError(t, err)
			assert.Equal(t, clientKind, log["http.kind"])
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

func TestClientMiddlewareMetrics(t *testing.T) {
	tests := []struct {
		name                string
		req                 *http.Request
		resDelay            time.Duration
		resError            error
		resStatusCode       int
		expectedMethod      string
		expectedURL         string
		expectedStatusCode  int
		expectedStatusClass string
	}{
		{
			name:                "Error",
			req:                 httptest.NewRequest("GET", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resError:            errors.New("uknown error"),
			resStatusCode:       0,
			expectedMethod:      "GET",
			expectedURL:         "/v1/items",
			expectedStatusCode:  -1,
			expectedStatusClass: "",
		},
		{
			name:                "200",
			req:                 httptest.NewRequest("GET", "/v1/items", nil),
			resDelay:            10 * time.Millisecond,
			resError:            nil,
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
			resError:            nil,
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
			resError:            nil,
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
			resError:            nil,
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
			mid := NewClientMiddleware(nil, metricsFactory, nil)
			assert.NotNil(t, mid)

			// Test http doer
			doer := mid.Metrics(func(r *http.Request) (*http.Response, error) {
				time.Sleep(tc.resDelay)
				if tc.resError != nil {
					return nil, tc.resError
				}
				return &http.Response{StatusCode: tc.resStatusCode}, nil
			})

			// Make the mock request
			res, err := doer(tc.req)

			if tc.resError != nil {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.resStatusCode, res.StatusCode)
			}

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
				case clientGaugeMetricName:
					assert.Equal(t, promModel.MetricType_GAUGE, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case clientCounterMetricName:
					assert.Equal(t, promModel.MetricType_COUNTER, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case clientHistogramMetricName:
					assert.Equal(t, promModel.MetricType_HISTOGRAM, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				case clientSummaryMetricName:
					assert.Equal(t, promModel.MetricType_SUMMARY, *metricFamily.Type)
					verifyLabels(metricFamily.Metric[0].Label)
				}
			}
		})
	}
}

func TestClientMiddlewareInjectSpan(t *testing.T) {
	tracer := mocktracer.New()

	tests := []struct {
		name     string
		tracer   opentracing.Tracer
		req      *http.Request
		span     opentracing.Span
		expected bool
	}{
		{
			name:     "InjectSucceeds",
			tracer:   tracer,
			req:      httptest.NewRequest("GET", "/", nil),
			span:     tracer.StartSpan("test-span"),
			expected: true,
		},
		{
			name:     "InjectFails",
			tracer:   tracer,
			req:      httptest.NewRequest("GET", "/", nil),
			span:     &mockSpan{},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &ClientMiddleware{
				tracer: tc.tracer,
			}

			m.injectSpan(tc.req, tc.span)

			injectedSpanContext := extractSpanContext(tc.req, tc.tracer)
			assert.Equal(t, tc.expected, injectedSpanContext != nil)
		})
	}
}

func TestClientMiddlewareTracing(t *testing.T) {
	tests := []struct {
		name               string
		req                *http.Request
		parentSpan         opentracing.Span
		resDelay           time.Duration
		resError           error
		resStatusCode      int
		expectedProto      string
		expectedMethod     string
		expectedURL        string
		expectedStatusCode int
	}{
		{
			name:               "Error",
			req:                httptest.NewRequest("GET", "/v1/items", nil),
			parentSpan:         nil,
			resDelay:           10 * time.Millisecond,
			resError:           errors.New("uknown error"),
			resStatusCode:      0,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "GET",
			expectedURL:        "/v1/items",
			expectedStatusCode: -1,
		},
		{
			name:               "200",
			req:                httptest.NewRequest("GET", "/v1/items", nil),
			parentSpan:         nil,
			resDelay:           10 * time.Millisecond,
			resError:           nil,
			resStatusCode:      200,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "GET",
			expectedURL:        "/v1/items",
			expectedStatusCode: 200,
		},
		{
			name:               "301",
			req:                httptest.NewRequest("GET", "/v1/items/1234", nil),
			parentSpan:         nil,
			resDelay:           10 * time.Millisecond,
			resError:           nil,
			resStatusCode:      301,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "GET",
			expectedURL:        "/v1/items/1234",
			expectedStatusCode: 301,
		},
		{
			name:               "404",
			req:                httptest.NewRequest("POST", "/v1/items", nil),
			parentSpan:         nil,
			resDelay:           10 * time.Millisecond,
			resError:           nil,
			resStatusCode:      404,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "POST",
			expectedURL:        "/v1/items",
			expectedStatusCode: 404,
		},
		{
			name:               "500",
			req:                httptest.NewRequest("PUT", "/v1/items/1234", nil),
			parentSpan:         nil,
			resDelay:           10 * time.Millisecond,
			resError:           nil,
			resStatusCode:      500,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "PUT",
			expectedURL:        "/v1/items/1234",
			expectedStatusCode: 500,
		},
		{
			name:               "WithParentSpan",
			req:                httptest.NewRequest("DELETE", "/v1/items/1234", nil),
			parentSpan:         mocktracer.New().StartSpan("parent-span"),
			resDelay:           10 * time.Millisecond,
			resError:           nil,
			resStatusCode:      204,
			expectedProto:      "HTTP/1.1",
			expectedMethod:     "DELETE",
			expectedURL:        "/v1/items/1234",
			expectedStatusCode: 204,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var injectedSpanContext opentracing.SpanContext

			tracer := mocktracer.New()
			mid := &ClientMiddleware{tracer: tracer}

			// Insert the parent span if any
			if tc.parentSpan != nil {
				ctx := tc.req.Context()
				ctx = opentracing.ContextWithSpan(ctx, tc.parentSpan)
				tc.req = tc.req.WithContext(ctx)
			}

			// Test http doer
			doer := mid.Tracing(func(r *http.Request) (*http.Response, error) {
				time.Sleep(tc.resDelay)
				injectedSpanContext = extractSpanContext(r, tracer)
				if tc.resError != nil {
					return nil, tc.resError
				}
				return &http.Response{StatusCode: tc.resStatusCode}, nil
			})

			// Make the mock request
			res, err := doer(tc.req)

			if tc.resError != nil {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.resStatusCode, res.StatusCode)
			}

			// Verify traces

			span := tracer.FinishedSpans()[0]
			assert.Equal(t, injectedSpanContext, span.Context())
			assert.Equal(t, clientSpanName, span.OperationName)
			assert.Equal(t, tc.expectedProto, span.Tag("http.proto"))
			assert.Equal(t, tc.expectedMethod, span.Tag("http.method"))
			assert.Equal(t, tc.expectedURL, span.Tag("http.url"))
			assert.Equal(t, uint16(tc.expectedStatusCode), span.Tag("http.status_code"))

			if tc.parentSpan != nil {
				parentSpan, ok := tc.parentSpan.(*mocktracer.MockSpan)
				assert.True(t, ok)
				assert.Equal(t, parentSpan.SpanContext.SpanID, span.ParentID)
				assert.Equal(t, parentSpan.SpanContext.TraceID, span.SpanContext.TraceID)
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
