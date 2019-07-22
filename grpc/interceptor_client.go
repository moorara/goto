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

const (
	clientKind                = "client"
	clientSpanName            = "grpc-client-request"
	clientGaugeMetricName     = "grpc_client_requests"
	clientCounterMetricName   = "grpc_client_requests_total"
	clientHistogramMetricName = "grpc_client_request_duration_seconds"
	clientSummaryMetricName   = "grpc_client_request_duration_quantiles_seconds"
)

// ClientInterceptor is a gRPC client interceptor for logging, metrics, and tracing
type ClientInterceptor struct {
	logger  *log.Logger
	metrics *metrics.RequestMetrics
	tracer  opentracing.Tracer
}

// NewClientInterceptor creates a new instance of gRPC server interceptor
func NewClientInterceptor(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ClientInterceptor {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(clientGaugeMetricName, "gauge metric for number of active client-side grpc requests", []string{"package", "service", "method", "stream"}),
		ReqCounter:      mf.Counter(clientCounterMetricName, "counter metric for total number of client-side grpc requests", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationHist: mf.Histogram(clientHistogramMetricName, "histogram metric for duration of client-side grpc requests in seconds", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationSumm: mf.Summary(clientSummaryMetricName, "summary metric for duration of client-side grpc requests in seconds", []string{"package", "service", "method", "stream", "success"}),
	}

	return &ClientInterceptor{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (i *ClientInterceptor) createSpan(ctx context.Context) opentracing.Span {
	var span opentracing.Span

	// Get trace information from the context if passed
	parentSpan := opentracing.SpanFromContext(ctx)

	if parentSpan == nil {
		span = i.tracer.StartSpan(clientSpanName)
	} else {
		span = i.tracer.StartSpan(clientSpanName, opentracing.ChildOf(parentSpan.Context()))
	}

	return span
}

func (i *ClientInterceptor) injectSpan(ctx context.Context, span opentracing.Span) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}

	carrier := &metadataTextMap{md}
	err := i.tracer.Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		span.LogFields(
			opentracingLog.Error(err),
			opentracingLog.String("message", "Tracer.Inject() failed"),
		)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

func (i *ClientInterceptor) injectRequestID(ctx context.Context, requestID string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}

	md.Set(requestIDKey, requestID)

	return metadata.NewOutgoingContext(ctx, md)
}

// UnaryInterceptor is the gRPC UnaryClientInterceptor for logging, metrics, and tracing
func (i *ClientInterceptor) UnaryInterceptor(ctx context.Context, fullMethod string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	stream := "false"
	pkg, service, method, ok := parseMethod(fullMethod)
	if !ok {
		return invoker(ctx, method, req, res, cc, opts...)
	}

	// Increment guage metric
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Inc()

	// Create a new span
	span := i.createSpan(ctx)
	defer span.Finish()

	// Propagate the current trace
	ctx = i.injectSpan(ctx, span)

	// Get request id from context
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	if !ok || requestID == "" {
		requestID = uuid.New().String()
	}

	// Propagate the request id
	ctx = i.injectRequestID(ctx, requestID)

	// Invoke the gRPC method
	start := time.Now()
	err := invoker(ctx, fullMethod, req, res, cc, opts...)
	success := err == nil
	duration := time.Since(start).Seconds()

	pairs := []interface{}{
		"grpc.kind", clientKind,
		"grpc.package", pkg,
		"grpc.service", service,
		"grpc.method", method,
		"grpc.stream", stream,
		"grpc.success", success,
		"responseTime", duration,
		"message", fmt.Sprintf("%s %s.%s.%s %f", clientKind, pkg, service, method, duration),
	}

	if err != nil {
		pairs = append(pairs, "grpc.error", err.Error())
	}

	// requestID is not empty at this point
	pairs = append(pairs, "requestId", requestID)

	if success {
		i.logger.Info(pairs...)
	} else {
		i.logger.Error(pairs...)
	}

	// Metrics
	successText := strconv.FormatBool(success)
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Dec()
	i.metrics.ReqCounter.WithLabelValues(pkg, service, method, stream, successText).Inc()
	i.metrics.ReqDurationHist.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)
	i.metrics.ReqDurationSumm.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)

	// Tracing
	// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	span.SetTag("grpc.package", pkg).SetTag("grpc.service", service).SetTag("grpc.method", method).SetTag("grpc.stream", stream).SetTag("grpc.success", success)

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracingLog.String("grpc.error", err.Error()),
		)
	}

	return err
}

// StreamInterceptor is the gRPC StreamClientInterceptor for logging, metrics, and tracing
func (i *ClientInterceptor) StreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	stream := "true"
	pkg, service, method, ok := parseMethod(fullMethod)
	if !ok {
		return streamer(ctx, desc, cc, method, opts...)
	}

	// Increment guage metric
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Inc()

	// Create a new span
	span := i.createSpan(ctx)
	defer span.Finish()

	// Propagate the current trace
	ctx = i.injectSpan(ctx, span)

	// Get request id from context
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	if !ok || requestID == "" {
		requestID = uuid.New().String()
	}

	// Propagate the request id
	ctx = i.injectRequestID(ctx, requestID)

	// Invoke the gRPC streaming method
	start := time.Now()
	cs, err := streamer(ctx, desc, cc, fullMethod, opts...)
	success := err == nil
	duration := time.Since(start).Seconds()

	pairs := []interface{}{
		"grpc.kind", clientKind,
		"grpc.package", pkg,
		"grpc.service", service,
		"grpc.method", method,
		"grpc.stream", stream,
		"grpc.success", success,
		"responseTime", duration,
		"message", fmt.Sprintf("%s %s.%s.%s %f", clientKind, pkg, service, method, duration),
	}

	if err != nil {
		pairs = append(pairs, "grpc.error", err.Error())
	}

	// requestID is not empty at this point
	pairs = append(pairs, "requestId", requestID)

	if success {
		i.logger.Info(pairs...)
	} else {
		i.logger.Error(pairs...)
	}

	// Metrics
	successText := strconv.FormatBool(success)
	i.metrics.ReqGauge.WithLabelValues(pkg, service, method, stream).Dec()
	i.metrics.ReqCounter.WithLabelValues(pkg, service, method, stream, successText).Inc()
	i.metrics.ReqDurationHist.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)
	i.metrics.ReqDurationSumm.WithLabelValues(pkg, service, method, stream, successText).Observe(duration)

	// Tracing
	// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
	ext.SpanKind.Set(span, ext.SpanKindRPCClientEnum)
	span.SetTag("grpc.package", pkg).SetTag("grpc.service", service).SetTag("grpc.method", method).SetTag("grpc.stream", stream).SetTag("grpc.success", success)

	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(
			opentracingLog.String("grpc.error", err.Error()),
		)
	}

	return cs, err
}
