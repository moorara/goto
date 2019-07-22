package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

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

// ServerMiddleware is an http server middleware for logging, metrics, tracing, etc.
type ServerMiddleware struct {
	logger  *log.Logger
	metrics *metrics.RequestMetrics
	tracer  opentracing.Tracer
}

// NewServerMiddleware creates a new instance of http server middleware
func NewServerMiddleware(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ServerMiddleware {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(serverGaugeMetricName, "gauge metric for number of active server-side http requests", []string{"method", "url"}),
		ReqCounter:      mf.Counter(serverCounterMetricName, "counter metric for total number of server-side http requests", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationHist: mf.Histogram(serverHistogramMetricName, "histogram metric for duration of server-side http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationSumm: mf.Summary(serverSummaryMetricName, "summary metric for duration of server-side http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
	}

	return &ServerMiddleware{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

// RequestID ensures incoming requests have unique ids
// This middleware ensures the request headers and context have a unique id
// A new request id will be generated if needed
func (m *ServerMiddleware) RequestID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure request id in headers
		requestID := r.Header.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
			r.Header.Set(requestIDHeader, requestID)
		}

		// Add request id to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, requestIDContextKey, requestID)
		req := r.WithContext(ctx)

		// Add request id to response headers
		w.Header().Set(requestIDHeader, requestID)

		// Call the next http handler
		next(w, req)
	}
}

// Logging takes care of logging for incoming http requests
// Request id will be read from reqeust headers if present
func (m *ServerMiddleware) Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proto := r.Proto
		method := r.Method
		url := r.URL.Path

		// Create a new logger that logs the context of current request
		logger := m.logger.With(
			"http.kind", serverKind,
			"req.proto", proto,
			"req.method", method,
			"req.url", url,
		)

		if requestID := r.Header.Get(requestIDHeader); requestID != "" {
			logger = logger.With("requestId", requestID)
		}

		// Update request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, loggerContextKey, logger)
		req := r.WithContext(ctx)

		// Call the next http handler
		start := time.Now()
		rw := NewResponseWriter(w)
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
	}
}

// Metrics takes care of metrics for incoming http requests
func (m *ServerMiddleware) Metrics(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		url := r.URL.Path

		// Increment guage metric
		m.metrics.ReqGauge.WithLabelValues(method, url).Inc()

		// Call the next http handler
		start := time.Now()
		rw := NewResponseWriter(w)
		next(rw, r)
		statusCode := rw.StatusCode
		statusClass := rw.StatusClass
		duration := time.Since(start).Seconds()

		// Metrics
		statusText := strconv.Itoa(statusCode)
		m.metrics.ReqGauge.WithLabelValues(method, url).Dec()
		m.metrics.ReqCounter.WithLabelValues(method, url, statusText, statusClass).Inc()
		m.metrics.ReqDurationHist.WithLabelValues(method, url, statusText, statusClass).Observe(duration)
		m.metrics.ReqDurationSumm.WithLabelValues(method, url, statusText, statusClass).Observe(duration)
	}
}

func (m *ServerMiddleware) createSpan(r *http.Request) opentracing.Span {
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

// Tracing takes care of tracing for incoming http requests
// Trace information will be read from reqeust headers if present
func (m *ServerMiddleware) Tracing(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proto := r.Proto
		method := r.Method
		url := r.URL.Path

		// Create a new span
		span := m.createSpan(r)
		defer span.Finish()

		// Update request context
		ctx := r.Context()
		ctx = opentracing.ContextWithSpan(ctx, span)
		req := r.WithContext(ctx)

		// Call the next http handler
		rw := NewResponseWriter(w)
		next(rw, req)
		statusCode := rw.StatusCode

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
