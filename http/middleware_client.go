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
	opentracingLog "github.com/opentracing/opentracing-go/log"
)

const (
	clientKind                = "client"
	clientSpanName            = "http-client-request"
	clientGaugeMetricName     = "http_client_requests"
	clientCounterMetricName   = "http_client_requests_total"
	clientHistogramMetricName = "http_client_request_duration_seconds"
	clientSummaryMetricName   = "http_client_request_duration_quantiles_seconds"
)

type (
	// Doer is the interface for standard http.Client Do method
	Doer func(*http.Request) (*http.Response, error)

	// ClientObservabilityMiddleware is an http client middleware for logging, metrics, and tracing
	ClientObservabilityMiddleware struct {
		logger  *log.Logger
		metrics *metrics.RequestMetrics
		tracer  opentracing.Tracer
	}
)

// NewClientObservabilityMiddleware creates a new instance of http client middleware for observability
func NewClientObservabilityMiddleware(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ClientObservabilityMiddleware {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(clientGaugeMetricName, "gauge metric for number of active client-side http requests", []string{"method", "url"}),
		ReqCounter:      mf.Counter(clientCounterMetricName, "counter metric for total number of client-side http requests", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationHist: mf.Histogram(clientHistogramMetricName, "histogram metric for duration of client-side http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
		ReqDurationSumm: mf.Summary(clientSummaryMetricName, "summary metric for duration of client-side http requests in seconds", []string{"method", "url", "statusCode", "statusClass"}),
	}

	return &ClientObservabilityMiddleware{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (m *ClientObservabilityMiddleware) createSpan(ctx context.Context) opentracing.Span {
	var span opentracing.Span

	// Get trace information from the context if passed
	parentSpan := opentracing.SpanFromContext(ctx)

	if parentSpan == nil {
		span = m.tracer.StartSpan(clientSpanName)
	} else {
		span = m.tracer.StartSpan(clientSpanName, opentracing.ChildOf(parentSpan.Context()))
	}

	return span
}

func (m *ClientObservabilityMiddleware) injectSpan(req *http.Request, span opentracing.Span) {
	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err := m.tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		span.LogFields(
			opentracingLog.Error(err),
			opentracingLog.String("message", "Tracer.Inject() failed"),
		)
	}
}

// Wrap takes care of logging, metrics, and tracing for client-side http requests
func (m *ClientObservabilityMiddleware) Wrap(ctx context.Context, req *http.Request, doer Doer) (*http.Response, error) {
	proto := req.Proto
	method := req.Method
	url := req.URL.Path

	// Increment guage metric
	m.metrics.ReqGauge.WithLabelValues(method, url).Inc()

	// Create a new span
	span := m.createSpan(ctx)
	defer span.Finish()

	// Propagate the current trace
	m.injectSpan(req, span)

	start := time.Now()
	res, err := doer(req)
	duration := time.Since(start).Seconds()

	var statusCode int
	var statusClass string

	if err != nil {
		statusCode = -1
		statusClass = ""
	} else {
		statusCode = res.StatusCode
		statusClass = fmt.Sprintf("%dxx", statusCode/100)
	}

	pairs := []interface{}{
		"req.proto", proto,
		"req.method", method,
		"req.url", url,
		"res.statusCode", statusCode,
		"res.statusClass", statusClass,
		"responseTime", duration,
		"message", fmt.Sprintf("%s %s %d %f", method, url, statusCode, duration),
	}

	// Logging
	switch {
	case statusCode >= 500:
		m.logger.Error(pairs...)
	case statusCode >= 400:
		m.logger.Warn(pairs...)
	case statusCode >= 100:
		fallthrough
	default:
		m.logger.Info(pairs...)
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

	return res, err
}
