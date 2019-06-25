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

var loggerContextKey = contextKey("logger")

// LoggerForRequest returns a logger set by http middleware on each request context
func LoggerForRequest(r *http.Request) (*log.Logger, bool) {
	ctx := r.Context()
	val := ctx.Value(loggerContextKey)
	logger, ok := val.(*log.Logger)

	return logger, ok
}

const (
	defaultSpanName     = "http-request"
	gaugeMetricName     = "http_requests"
	counterMetricName   = "http_requests_total"
	histogramMetricName = "http_request_duration_seconds"
	summaryMetricName   = "http_request_duration_quantiles_seconds"
)

// ObservabilityMiddleware is an http middleware for logging, metrics, and tracing
type ObservabilityMiddleware struct {
	logger  *log.Logger
	metrics *metrics.RequestMetrics
	tracer  opentracing.Tracer
}

// NewObservabilityMiddleware creates a new instance of http middleware for observability
func NewObservabilityMiddleware(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ObservabilityMiddleware {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(gaugeMetricName, "gauge metric for number of active http requests", []string{"method", "url"}),
		ReqCounter:      mf.Counter(counterMetricName, "counter metric for total number of http requests", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationHist: mf.Histogram(histogramMetricName, "histogram metric for duration of http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationSumm: mf.Summary(summaryMetricName, "summary metric for duration of http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
	}

	return &ObservabilityMiddleware{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (m *ObservabilityMiddleware) createSpan(r *http.Request) opentracing.Span {
	var span opentracing.Span

	carrier := opentracing.HTTPHeadersCarrier(r.Header)
	parentSpanContext, err := m.tracer.Extract(opentracing.HTTPHeaders, carrier)
	if err != nil {
		span = m.tracer.StartSpan(defaultSpanName)
	} else {
		span = m.tracer.StartSpan(defaultSpanName, opentracing.ChildOf(parentSpanContext))
	}

	return span
}

// Wrap accepts an http handler and return a new http handler that takes care of logging, metrics, and tracing
func (m *ObservabilityMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		proto := r.Proto
		method := r.Method
		url := r.URL.Path

		// Increment guage metric
		m.metrics.ReqGauge.WithLabelValues(method, url).Inc()

		// Create a new logger that logs the context of current request
		logger := m.logger.With(
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
