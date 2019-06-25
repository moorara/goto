package grpc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	serverKind                = "server"
	serverSpanName            = "grpc-server-request"
	serverGaugeMetricName     = "grpc_server_requests"
	serverCounterMetricName   = "grpc_server_requests_total"
	serverHistogramMetricName = "grpc_server_request_duration_seconds"
	serverSummaryMetricName   = "grpc_server_request_duration_quantiles_seconds"
)

// ServerObservabilityInterceptor is a gRPC server interceptor for logging, metrics, and tracing
type ServerObservabilityInterceptor struct {
	logger  *log.Logger
	metrics *metrics.RequestMetrics
	tracer  opentracing.Tracer
}

// NewServerObservabilityInterceptor creates a new instance of gRPC server interceptor for observability
func NewServerObservabilityInterceptor(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ServerObservabilityInterceptor {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(serverGaugeMetricName, "gauge metric for number of active grpc server requests", []string{"package", "service", "method", "stream"}),
		ReqCounter:      mf.Counter(serverCounterMetricName, "counter metric for total number of grpc server requests", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationHist: mf.Histogram(serverHistogramMetricName, "histogram metric for duration of grpc server requests in seconds", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationSumm: mf.Summary(serverSummaryMetricName, "summary metric for duration of grpc server requests in seconds", []string{"package", "service", "method", "stream", "success"}),
	}

	return &ServerObservabilityInterceptor{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (i *ServerObservabilityInterceptor) createSpan(ctx context.Context) opentracing.Span {
	var span opentracing.Span
	var parentSpanContext opentracing.SpanContext

	// Get trace information from incoming metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		carrier := &MetadataTextMap{md}
		parentSpanContext, _ = i.tracer.Extract(opentracing.TextMap, carrier)
		// In case of error, we just create a new span without parent and start a new trace!
	}

	if parentSpanContext == nil {
		span = i.tracer.StartSpan(serverSpanName)
	} else {
		span = i.tracer.StartSpan(serverSpanName, opentracing.ChildOf(parentSpanContext))
	}

	return span
}

// UnaryInterceptor is the gRPC UnaryServerInterceptor for logging, metrics, and tracing
func (i *ServerObservabilityInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	stream := "false"
	pkg, service, method, ok := parseMethod(info.FullMethod)
	if !ok {
		return handler(ctx, req)
	}

	// Increment guage metric
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Inc()

	// Create a new logger that logs the context
	logger := i.logger.With(
		"grpc.kind", serverKind,
		"grpc.package", pkg,
		"grpc.service", service,
		"grpc.method", method,
		"grpc.stream", stream,
	)

	// Create a new span
	span := i.createSpan(ctx)
	defer span.Finish()

	// Update request context
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, loggerContextKey, logger)

	// Call the gRPC method handler
	start := time.Now()
	res, err := handler(ctx, req)
	success := err == nil
	duration := time.Since(start).Seconds()

	pairs := []interface{}{
		"grpc.success", success,
		"responseTime", duration,
		"message", fmt.Sprintf("%s %s.%s.%s %f", serverKind, pkg, service, method, duration),
	}

	if success {
		logger.Info(pairs...)
	} else {
		logger.Error(pairs...)
	}

	// Metrics
	successText := strconv.FormatBool(success)
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Dec()
	i.metrics.ReqCounter.WithLabelValues(pkg, service, method, stream, successText).Inc()
	i.metrics.ReqDurationHist.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)
	i.metrics.ReqDurationSumm.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)

	// Tracing
	// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
	ext.SpanKind.Set(span, ext.SpanKindRPCServerEnum)
	span.SetTag("grpc.package", pkg).SetTag("grpc.service", service).SetTag("grpc.method", method).SetTag("grpc.stream", stream).SetTag("grpc.success", success)
	/* span.LogFields(
		opentracingLog.String("key", value),
	) */

	return res, err
}

// StreamInterceptor is the gRPC StreamServerInterceptor for logging, metrics, and tracing
func (i *ServerObservabilityInterceptor) StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	stream := "true"
	pkg, service, method, ok := parseMethod(info.FullMethod)
	if !ok {
		return handler(srv, ss)
	}

	// Increment guage metric
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Inc()

	// Create a new logger that logs the context
	logger := i.logger.With(
		"grpc.kind", serverKind,
		"grpc.package", pkg,
		"grpc.service", service,
		"grpc.method", method,
		"grpc.stream", stream,
	)

	// Create a new span
	span := i.createSpan(ss.Context())
	defer span.Finish()

	// Update stream context
	ctx := ss.Context()
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, loggerContextKey, logger)
	ss = ServerStreamWithContext(ss, ctx)

	// Call the gRPC streaming method handler
	start := time.Now()
	err := handler(srv, ss)
	success := err == nil
	duration := time.Since(start).Seconds()

	pairs := []interface{}{
		"grpc.success", success,
		"responseTime", duration,
		"message", fmt.Sprintf("%s %s.%s.%s %f", serverKind, pkg, service, method, duration),
	}

	if success {
		logger.Info(pairs...)
	} else {
		logger.Error(pairs...)
	}

	// Metrics
	successText := strconv.FormatBool(success)
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Dec()
	i.metrics.ReqCounter.WithLabelValues(pkg, service, method, stream, successText).Inc()
	i.metrics.ReqDurationHist.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)
	i.metrics.ReqDurationSumm.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)

	// Tracing
	// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
	ext.SpanKind.Set(span, ext.SpanKindRPCServerEnum)
	span.SetTag("grpc.package", pkg).SetTag("grpc.service", service).SetTag("grpc.method", method).SetTag("grpc.stream", stream).SetTag("grpc.success", success)
	/* span.LogFields(
		opentracingLog.String("key", value),
	) */

	return err
}
