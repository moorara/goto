package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// contextKey is the type for the keys added to context
type contextKey string

const (
	loggerContextKey = contextKey("logger")

	serverKind                = "server"
	serverSpanName            = "http-server-request"
	serverGaugeMetricName     = "http_server_requests"
	serverCounterMetricName   = "http_server_requests_total"
	serverHistogramMetricName = "http_server_request_duration_seconds"
	serverSummaryMetricName   = "http_server_request_duration_quantiles_seconds"
)

// LoggerForRequest returns a logger set by http middleware on each request context
func LoggerForRequest(r *http.Request) (*log.Logger, bool) {
	ctx := r.Context()
	val := ctx.Value(loggerContextKey)
	logger, ok := val.(*log.Logger)

	return logger, ok
}

// ServerObservabilityMiddleware is an http server middleware for logging, metrics, and tracing
type ServerObservabilityMiddleware struct {
	logger  *log.Logger
	metrics *metrics.RequestMetrics
	tracer  opentracing.Tracer
}

// NewServerObservabilityMiddleware creates a new instance of http server middleware for observability
func NewServerObservabilityMiddleware(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ServerObservabilityMiddleware {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(serverGaugeMetricName, "gauge metric for number of active server-side http requests", []string{"method", "url"}),
		ReqCounter:      mf.Counter(serverCounterMetricName, "counter metric for total number of server-side http requests", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationHist: mf.Histogram(serverHistogramMetricName, "histogram metric for duration of server-side http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationSumm: mf.Summary(serverSummaryMetricName, "summary metric for duration of server-side http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
	}

	return &ServerObservabilityMiddleware{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (m *ServerObservabilityMiddleware) createSpan(r *http.Request) opentracing.Span {
	var span opentracing.Span

	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	parentSpanContext, err := m.tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err != nil {
		span = m.tracer.StartSpan(serverSpanName)
	} else {
		span = m.tracer.StartSpan(serverSpanName, opentracing.ChildOf(parentSpanContext))
	}

	return span
}

// Wrap accepts an http handler and return a new http handler that takes care of logging, metrics, and tracing
func (m *ServerObservabilityMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proto := r.Proto
		method := r.Method
		url := r.URL.Path

		// Increment guage metric
		m.metrics.ReqGauge.WithLabelValues(method, url).Inc()

		// Create a new logger that logs the context of current request
		logger := m.logger.With(
			"http.kind", serverKind,
			"req.proto", proto,
			"req.method", method,
			"req.url", url,
		)

		// Create a new span
		span := m.createSpan(r)
		defer span.Finish()

		// Update request context
		ctx := r.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)
		ctx = context.WithValue(ctx, loggerContextKey, logger)

		// Call next http handler
		start := time.Now()
		rw := NewResponseWriter(w)
		req := r.WithContext(ctx)
		next(rw, req)
		statusCode := rw.StatusCode
		statusClass := rw.StatusClass
		duration := time.Since(start).Seconds()

		pairs := []interface{}{
			"res.statusCode", statusCode,
			"res.statusClass", statusClass,
			"responseTime", duration,
			"message", fmt.Sprintf("%s %s %d %f", method, url, statusCode, duration),
		}

		// Logging
		switch {
		case statusCode >= 500:
			logger.Error(pairs...)
		case statusCode >= 400:
			logger.Warn(pairs...)
		case statusCode >= 100:
			fallthrough
		default:
			logger.Info(pairs...)
		}

		// Metrics
		statusText := strconv.Itoa(statusCode)
		m.metrics.ReqGauge.WithLabelValues(method, url).Dec()
		m.metrics.ReqCounter.WithLabelValues(method, url, statusText, statusClass).Inc()
		m.metrics.ReqDurationHist.WithLabelValues(method, url, statusText, statusClass).Observe(duration)
		m.metrics.ReqDurationSumm.WithLabelValues(method, url, statusText, statusClass).Observe(duration)

		// Tracing
		// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
		span.SetTag("http.proto", proto)
		ext.HTTPMethod.Set(span, method)
		ext.HTTPUrl.Set(span, url)
		ext.HTTPStatusCode.Set(span, uint16(statusCode))
		/* span.LogFields(
			opentracingLog.String("key", value),
		) */
	}
}
