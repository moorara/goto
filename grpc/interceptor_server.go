package grpc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var loggerContextKey = contextKey("logger")

// LoggerFromContext returns a logger set by grpc server interceptor on each incoming context
func LoggerFromContext(ctx context.Context) (*log.Logger, bool) {
	val := ctx.Value(loggerContextKey)
	logger, ok := val.(*log.Logger)

	return logger, ok
}

const (
	serverKind                = "server"
	serverSpanName            = "grpc-server-request"
	serverGaugeMetricName     = "grpc_server_requests"
	serverCounterMetricName   = "grpc_server_requests_total"
	serverHistogramMetricName = "grpc_server_request_duration_seconds"
	serverSummaryMetricName   = "grpc_server_request_duration_quantiles_seconds"
)

// ServerInterceptor is a gRPC server interceptor for logging, metrics, and tracing
type ServerInterceptor struct {
	logger  *log.Logger
	metrics *metrics.RequestMetrics
	tracer  opentracing.Tracer
}

// NewServerInterceptor creates a new instance of gRPC server interceptor
func NewServerInterceptor(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ServerInterceptor {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(serverGaugeMetricName, "gauge metric for number of active server-side grpc requests", []string{"package", "service", "method", "stream"}),
		ReqCounter:      mf.Counter(serverCounterMetricName, "counter metric for total number of server-side grpc requests", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationHist: mf.Histogram(serverHistogramMetricName, "histogram metric for duration of server-side grpc requests in seconds", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationSumm: mf.Summary(serverSummaryMetricName, "summary metric for duration of server-side grpc requests in seconds", []string{"package", "service", "method", "stream", "success"}),
	}

	return &ServerInterceptor{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (i *ServerInterceptor) createSpan(ctx context.Context) opentracing.Span {
	var span opentracing.Span
	var parentSpanContext opentracing.SpanContext

	// Get trace information from incoming metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		carrier := &metadataTextMap{md}
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

func (i *ServerInterceptor) getRequestID(ctx context.Context) string {
	var requestID string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		vals := md.Get(requestIDKey)
		if len(vals) > 0 {
			requestID = vals[0]
		}
	}

	if requestID == "" {
		requestID = uuid.New().String()
	}

	return requestID
}

// UnaryInterceptor is the gRPC UnaryServerInterceptor for logging, metrics, and tracing
func (i *ServerInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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

	// Get or generate request id
	requestID := i.getRequestID(ctx)

	// Capture the request id in logs
	logger = logger.With("requestId", requestID)

	// Update request context
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, requestIDContextKey, requestID)
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

	if err != nil {
		pairs = append(pairs, "grpc.error", err.Error())
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

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracingLog.String("grpc.error", err.Error()),
		)
	}

	return res, err
}

// StreamInterceptor is the gRPC StreamServerInterceptor for logging, metrics, and tracing
func (i *ServerInterceptor) StreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()

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
	span := i.createSpan(ctx)
	defer span.Finish()

	// Get or generate request id
	requestID := i.getRequestID(ctx)

	// Capture the request id in logs
	logger = logger.With("requestId", requestID)

	// Update stream context
	ctx = opentracing.ContextWithSpan(ctx, span)
	ctx = context.WithValue(ctx, requestIDContextKey, requestID)
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

	if err != nil {
		pairs = append(pairs, "grpc.error", err.Error())
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

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracingLog.String("grpc.error", err.Error()),
		)
	}

	return err
}
