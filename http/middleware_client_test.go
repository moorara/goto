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
	"github.com/moorara/goto/trace"
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

func TestNewClientObservabilityMiddleware(t *testing.T) {
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
			m := NewClientObservabilityMiddleware(tc.logger, tc.mf, tc.tracer)

			assert.Equal(t, tc.logger, m.logger)
			assert.NotNil(t, m.metrics)
			assert.Equal(t, tc.tracer, m.tracer)
		})
	}
}

func TestClientObservabilityMiddlewareInjectSpan(t *testing.T) {
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
			m := &ClientObservabilityMiddleware{
				tracer: tc.tracer,
			}

			m.injectSpan(tc.req, tc.span)

			injectedSpanContext := extractSpanContext(tc.req, tc.tracer)
			assert.Equal(t, tc.expected, injectedSpanContext != nil)
		})
	}
}

func TestClientObservabilityMiddlewareWrap(t *testing.T) {
	tests := []struct {
		name                string
		parentSpan          opentracing.Span
		ctx                 context.Context
		req                 *http.Request
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
			parentSpan:          nil,
			ctx:                 context.Background(),
			req:                 httptest.NewRequest("GET", "/v1/dogs", nil),
			resDelay:            10 * time.Millisecond,
			resError:            errors.New("reachability error"),
			resStatusCode:       0,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/dogs",
			expectedStatusCode:  -1,
			expectedStatusClass: "",
		},
		{
			name:                "200",
			parentSpan:          nil,
			ctx:                 context.Background(),
			req:                 httptest.NewRequest("GET", "/v1/dogs/breeds", nil),
			resDelay:            10 * time.Millisecond,
			resError:            nil,
			resStatusCode:       200,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/dogs/breeds",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
		{
			name:                "301",
			parentSpan:          nil,
			ctx:                 context.Background(),
			req:                 httptest.NewRequest("GET", "/v1/dogs/breeds/1234", nil),
			resDelay:            10 * time.Millisecond,
			resError:            nil,
			resStatusCode:       301,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "GET",
			expectedURL:         "/v1/dogs/breeds/1234",
			expectedStatusCode:  301,
			expectedStatusClass: "3xx",
		},
		{
			name:                "404",
			parentSpan:          nil,
			ctx:                 context.Background(),
			req:                 httptest.NewRequest("POST", "/v1/breeds/dogs", nil),
			resDelay:            10 * time.Millisecond,
			resError:            nil,
			resStatusCode:       404,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "POST",
			expectedURL:         "/v1/breeds/dogs",
			expectedStatusCode:  404,
			expectedStatusClass: "4xx",
		},
		{
			name:                "500",
			parentSpan:          nil,
			ctx:                 context.Background(),
			req:                 httptest.NewRequest("PUT", "/v1/dogs/breeds/abcd", nil),
			resDelay:            10 * time.Millisecond,
			resError:            nil,
			resStatusCode:       500,
			expectedProto:       "HTTP/1.1",
			expectedMethod:      "PUT",
			expectedURL:         "/v1/dogs/breeds/abcd",
			expectedStatusCode:  500,
			expectedStatusClass: "5xx",
		},
		{
			name:                "WithParentSpan",
			parentSpan:          mocktracer.New().StartSpan("parent-span"),
			ctx:                 context.Background(),
			req:                 httptest.NewRequest("DELETE", "/v1/dogs/breeds/1234-abcd", nil),
			resDelay:            10 * time.Millisecond,
			resError:            nil,
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
			buff := &bytes.Buffer{}
			var injectedSpanContext opentracing.SpanContext

			logger := log.NewLogger(log.Options{Writer: buff})
			promReg := prometheus.NewRegistry()
			metricsFactory := metrics.NewFactory(metrics.FactoryOptions{Registerer: promReg})
			tracer := mocktracer.New()

			// Create http client middleware
			mid := NewClientObservabilityMiddleware(logger, metricsFactory, tracer)
			assert.NotNil(t, mid)

			// Insert the parent span if any
			if tc.parentSpan != nil {
				tc.ctx = opentracing.ContextWithSpan(tc.ctx, tc.parentSpan)
			}

			// Test http doer
			doer := func(req *http.Request) (*http.Response, error) {
				time.Sleep(tc.resDelay)
				injectedSpanContext = extractSpanContext(req, tracer)
				if tc.resError != nil {
					return nil, tc.resError
				}
				return &http.Response{StatusCode: tc.resStatusCode}, nil
			}

			// Wrap and make the request
			res, err := mid.Wrap(tc.ctx, tc.req, doer)

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
		})
	}
}
