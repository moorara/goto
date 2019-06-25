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

// ClientObservabilityInterceptor is a gRPC client interceptor for logging, metrics, and tracing
type ClientObservabilityInterceptor struct {
	logger  *log.Logger
	metrics *metrics.RequestMetrics
	tracer  opentracing.Tracer
}

// NewClientObservabilityInterceptor creates a new instance of gRPC server interceptor for observability
func NewClientObservabilityInterceptor(logger *log.Logger, mf *metrics.Factory, tracer opentracing.Tracer) *ClientObservabilityInterceptor {
	metrics := &metrics.RequestMetrics{
		ReqGauge:        mf.Gauge(clientGaugeMetricName, "gauge metric for number of active grpc client requests", []string{"package", "service", "method", "stream"}),
		ReqCounter:      mf.Counter(clientCounterMetricName, "counter metric for total number of grpc client requests", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationHist: mf.Histogram(clientHistogramMetricName, "histogram metric for duration of grpc client requests in seconds", []string{"package", "service", "method", "stream", "success"}),
		ReqDurationSumm: mf.Summary(clientSummaryMetricName, "summary metric for duration of grpc client requests in seconds", []string{"package", "service", "method", "stream", "success"}),
	}

	return &ClientObservabilityInterceptor{
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
	}
}

func (i *ClientObservabilityInterceptor) createSpan(ctx context.Context) opentracing.Span {
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

func (i *ClientObservabilityInterceptor) injectTrace(ctx context.Context, span opentracing.Span) context.Context {
	// Get any metadata if set
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}

	carrier := &MetadataTextMap{md}
	err := i.tracer.Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		span.LogFields(
			opentracingLog.Error(err),
			opentracingLog.String("message", "Tracer.Inject() failed"),
		)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// UnaryInterceptor is the gRPC UnaryClientInterceptor for logging, metrics, and tracing
func (i *ClientObservabilityInterceptor) UnaryInterceptor(ctx context.Context, fullMethod string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
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
	ctx = i.injectTrace(ctx, span)

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
	/* span.LogFields(
		opentracingLog.String("key", value),
	) */

	return err
}

// StreamInterceptor is the gRPC StreamClientInterceptor for logging, metrics, and tracing
func (i *ClientObservabilityInterceptor) StreamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, fullMethod string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
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
	ctx = i.injectTrace(ctx, span)

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
	/* span.LogFields(
		opentracingLog.String("key", value),
	) */

	return cs, err
}