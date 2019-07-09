package grpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/prometheus/client_golang/prometheus"
	promModel "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func extractSpanContext(ctx context.Context, tracer opentracing.Tracer) opentracing.SpanContext {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		carrier := &metadataTextMap{md}
		parentSpanContext, err := tracer.Extract(opentracing.TextMap, carrier)
		if err == nil {
			return parentSpanContext
		}
	}

	return nil
}

func TestNewClientObservabilityInterceptor(t *testing.T) {
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
			coi := NewClientObservabilityInterceptor(tc.logger, tc.mf, tc.tracer)

			assert.Equal(t, tc.logger, coi.logger)
			assert.NotNil(t, coi.metrics)
			assert.Equal(t, tc.tracer, coi.tracer)
		})
	}
}

func TestInjectSpan(t *testing.T) {
	tracer := mocktracer.New()
	ctx := context.Background()
	ctxWithMD := metadata.NewOutgoingContext(ctx, metadata.Pairs("key", "value"))

	tests := []struct {
		name     string
		tracer   opentracing.Tracer
		ctx      context.Context
		span     opentracing.Span
		expected bool
	}{
		{
			name:     "WithoutMetadata",
			tracer:   tracer,
			ctx:      ctx,
			span:     tracer.StartSpan("test-span"),
			expected: true,
		},
		{
			name:     "WithMetadata",
			tracer:   tracer,
			ctx:      ctxWithMD,
			span:     tracer.StartSpan("test-span"),
			expected: true,
		},
		{
			name:     "InjectFails",
			tracer:   tracer,
			ctx:      ctx,
			span:     &mockSpan{},
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			i := &ClientObservabilityInterceptor{
				tracer: tc.tracer,
			}

			ctx := i.injectSpan(tc.ctx, tc.span)

			injectedSpanContext := extractSpanContext(ctx, tc.tracer)
			assert.Equal(t, tc.expected, injectedSpanContext != nil)
		})
	}
}

func TestUnaryClientInterceptor(t *testing.T) {
	tests := []struct {
		name            string
		parentSpan      opentracing.Span
		ctx             context.Context
		method          string
		req             interface{}
		res             interface{}
		cc              *grpc.ClientConn
		opts            []grpc.CallOption
		mockDelay       time.Duration
		mockRespError   error
		verify          bool
		expectedPackage string
		expectedService string
		expectedMethod  string
		expectedStream  string
		expectedSuccess bool
	}{
		{
			name:            "InvalidMethod",
			parentSpan:      nil,
			ctx:             context.Background(),
			method:          "",
			req:             nil,
			res:             nil,
			cc:              nil,
			opts:            nil,
			mockDelay:       0,
			mockRespError:   nil,
			verify:          false,
			expectedPackage: "",
			expectedService: "",
			expectedMethod:  "",
			expectedStream:  "",
			expectedSuccess: false,
		},
		{
			name:            "InvokerFails",
			parentSpan:      nil,
			ctx:             context.Background(),
			method:          "/package.service/method",
			req:             nil,
			res:             nil,
			cc:              &grpc.ClientConn{},
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   errors.New("error on grpc method"),
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "false",
			expectedSuccess: false,
		},
		{
			name:            "InvokerSucceeds",
			parentSpan:      nil,
			ctx:             context.Background(),
			method:          "/package.service/method",
			req:             nil,
			res:             nil,
			cc:              &grpc.ClientConn{},
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "false",
			expectedSuccess: true,
		},
		{
			name:            "InvokerSucceedsWithMetadata",
			parentSpan:      nil,
			ctx:             metadata.NewOutgoingContext(context.Background(), metadata.Pairs("key", "value")),
			method:          "/package.service/method",
			req:             nil,
			res:             nil,
			cc:              &grpc.ClientConn{},
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "false",
			expectedSuccess: true,
		},
		{
			name:            "InvokerSucceedsWithParentSpan",
			parentSpan:      mocktracer.New().StartSpan("parent-span"),
			ctx:             context.Background(),
			method:          "/package.service/method",
			req:             nil,
			res:             nil,
			cc:              &grpc.ClientConn{},
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "false",
			expectedSuccess: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buff := &bytes.Buffer{}
			var injectedSpanContext opentracing.SpanContext

			logger := log.NewLogger(log.Options{Writer: buff})
			promReg := prometheus.NewRegistry()
			mf := metrics.NewFactory(metrics.FactoryOptions{Registerer: promReg})
			tracer := mocktracer.New()

			// Create the interceptor
			i := NewClientObservabilityInterceptor(logger, mf, tracer)
			assert.NotNil(t, i)

			// Insert the parent span if any
			if tc.parentSpan != nil {
				tc.ctx = opentracing.ContextWithSpan(tc.ctx, tc.parentSpan)
			}

			invoker := func(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				time.Sleep(tc.mockDelay)
				injectedSpanContext = extractSpanContext(ctx, tracer)
				return tc.mockRespError
			}

			err := i.UnaryInterceptor(tc.ctx, tc.method, tc.req, tc.res, tc.cc, invoker, tc.opts...)
			assert.Equal(t, tc.mockRespError, err)

			if tc.verify {
				// Verify logs

				var log map[string]interface{}
				err := json.NewDecoder(buff).Decode(&log)
				assert.NoError(t, err)
				assert.Equal(t, clientKind, log["grpc.kind"])
				assert.Equal(t, tc.expectedPackage, log["grpc.package"])
				assert.Equal(t, tc.expectedService, log["grpc.service"])
				assert.Equal(t, tc.expectedMethod, log["grpc.method"])
				assert.Equal(t, tc.expectedStream, log["grpc.stream"])
				assert.Equal(t, tc.expectedSuccess, log["grpc.success"])
				assert.NotEmpty(t, log["responseTime"])
				assert.NotEmpty(t, log["message"])

				if tc.mockRespError != nil {
					assert.Equal(t, tc.mockRespError.Error(), log["grpc.error"])
				}

				// Verify metrics

				verifyLabels := func(labels []*promModel.LabelPair) {
					for _, l := range labels {
						switch *l.Name {
						case "package":
							assert.Equal(t, tc.expectedPackage, *l.Value)
						case "service":
							assert.Equal(t, tc.expectedService, *l.Value)
						case "method":
							assert.Equal(t, tc.expectedMethod, *l.Value)
						case "stream":
							assert.Equal(t, tc.expectedStream, *l.Value)
						case "success":
							assert.Equal(t, strconv.FormatBool(tc.expectedSuccess), *l.Value)
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
				assert.Equal(t, ext.SpanKindEnum("client"), span.Tag("span.kind"))
				assert.Equal(t, tc.expectedPackage, span.Tag("grpc.package"))
				assert.Equal(t, tc.expectedService, span.Tag("grpc.service"))
				assert.Equal(t, tc.expectedMethod, span.Tag("grpc.method"))
				assert.Equal(t, tc.expectedStream, span.Tag("grpc.stream"))
				assert.Equal(t, tc.expectedSuccess, span.Tag("grpc.success"))

				if tc.parentSpan != nil {
					parentSpan, ok := tc.parentSpan.(*mocktracer.MockSpan)
					assert.True(t, ok)
					assert.Equal(t, parentSpan.SpanContext.SpanID, span.ParentID)
					assert.Equal(t, parentSpan.SpanContext.TraceID, span.SpanContext.TraceID)
				}

				spanLogs := span.Logs()
				if tc.mockRespError != nil {
					lf := spanLogs[0].Fields[0]
					assert.Equal(t, "grpc.error", lf.Key)
					assert.Equal(t, tc.mockRespError.Error(), lf.ValueString)
				}
			}
		})
	}
}

func TestStreamClientInterceptor(t *testing.T) {
	tests := []struct {
		name            string
		parentSpan      opentracing.Span
		ctx             context.Context
		desc            *grpc.StreamDesc
		cc              *grpc.ClientConn
		method          string
		opts            []grpc.CallOption
		mockDelay       time.Duration
		mockRespError   error
		mockRespCS      grpc.ClientStream
		verify          bool
		expectedPackage string
		expectedService string
		expectedMethod  string
		expectedStream  string
		expectedSuccess bool
	}{
		{
			name:            "InvalidMethod",
			parentSpan:      nil,
			ctx:             context.Background(),
			desc:            nil,
			cc:              nil,
			method:          "",
			opts:            nil,
			mockDelay:       0,
			mockRespError:   nil,
			mockRespCS:      nil,
			verify:          false,
			expectedPackage: "",
			expectedService: "",
			expectedMethod:  "",
			expectedStream:  "",
			expectedSuccess: false,
		},
		{
			name:            "StreamerFails",
			parentSpan:      nil,
			ctx:             context.Background(),
			desc:            &grpc.StreamDesc{},
			cc:              &grpc.ClientConn{},
			method:          "/package.service/method",
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   errors.New("error on grpc method"),
			mockRespCS:      nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "true",
			expectedSuccess: false,
		},
		{
			name:            "StreamerSucceeds",
			parentSpan:      nil,
			ctx:             context.Background(),
			desc:            &grpc.StreamDesc{},
			cc:              &grpc.ClientConn{},
			method:          "/package.service/method",
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			mockRespCS:      nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "true",
			expectedSuccess: true,
		},
		{
			name:            "StreamerSucceedsWithMetadata",
			parentSpan:      nil,
			ctx:             metadata.NewOutgoingContext(context.Background(), metadata.Pairs("key", "value")),
			desc:            &grpc.StreamDesc{},
			cc:              &grpc.ClientConn{},
			method:          "/package.service/method",
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			mockRespCS:      nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "true",
			expectedSuccess: true,
		},
		{
			name:            "StreamerSucceedsWithParentSpan",
			parentSpan:      mocktracer.New().StartSpan("parent-span"),
			ctx:             context.Background(),
			desc:            &grpc.StreamDesc{},
			cc:              &grpc.ClientConn{},
			method:          "/package.service/method",
			opts:            []grpc.CallOption{},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			mockRespCS:      nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "true",
			expectedSuccess: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			buff := &bytes.Buffer{}
			var injectedSpanContext opentracing.SpanContext

			logger := log.NewLogger(log.Options{Writer: buff})
			promReg := prometheus.NewRegistry()
			mf := metrics.NewFactory(metrics.FactoryOptions{Registerer: promReg})
			tracer := mocktracer.New()

			// Create the interceptor
			i := NewClientObservabilityInterceptor(logger, mf, tracer)
			assert.NotNil(t, i)

			// Insert the parent span if any
			if tc.parentSpan != nil {
				tc.ctx = opentracing.ContextWithSpan(tc.ctx, tc.parentSpan)
			}

			streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				time.Sleep(tc.mockDelay)
				injectedSpanContext = extractSpanContext(ctx, tracer)
				return tc.mockRespCS, tc.mockRespError
			}

			cs, err := i.StreamInterceptor(tc.ctx, tc.desc, tc.cc, tc.method, streamer, tc.opts...)
			assert.Equal(t, tc.mockRespError, err)
			assert.Equal(t, tc.mockRespCS, cs)

			if tc.verify {
				// Verify logs

				var log map[string]interface{}
				err := json.NewDecoder(buff).Decode(&log)
				assert.NoError(t, err)
				assert.Equal(t, clientKind, log["grpc.kind"])
				assert.Equal(t, tc.expectedPackage, log["grpc.package"])
				assert.Equal(t, tc.expectedService, log["grpc.service"])
				assert.Equal(t, tc.expectedMethod, log["grpc.method"])
				assert.Equal(t, tc.expectedStream, log["grpc.stream"])
				assert.Equal(t, tc.expectedSuccess, log["grpc.success"])
				assert.NotEmpty(t, log["responseTime"])
				assert.NotEmpty(t, log["message"])

				if tc.mockRespError != nil {
					assert.Equal(t, tc.mockRespError.Error(), log["grpc.error"])
				}

				// Verify metrics

				verifyLabels := func(labels []*promModel.LabelPair) {
					for _, l := range labels {
						switch *l.Name {
						case "package":
							assert.Equal(t, tc.expectedPackage, *l.Value)
						case "service":
							assert.Equal(t, tc.expectedService, *l.Value)
						case "method":
							assert.Equal(t, tc.expectedMethod, *l.Value)
						case "stream":
							assert.Equal(t, tc.expectedStream, *l.Value)
						case "success":
							assert.Equal(t, strconv.FormatBool(tc.expectedSuccess), *l.Value)
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
				assert.Equal(t, ext.SpanKindEnum("client"), span.Tag("span.kind"))
				assert.Equal(t, tc.expectedPackage, span.Tag("grpc.package"))
				assert.Equal(t, tc.expectedService, span.Tag("grpc.service"))
				assert.Equal(t, tc.expectedMethod, span.Tag("grpc.method"))
				assert.Equal(t, tc.expectedStream, span.Tag("grpc.stream"))
				assert.Equal(t, tc.expectedSuccess, span.Tag("grpc.success"))

				if tc.parentSpan != nil {
					parentSpan, ok := tc.parentSpan.(*mocktracer.MockSpan)
					assert.True(t, ok)
					assert.Equal(t, parentSpan.SpanContext.SpanID, span.ParentID)
					assert.Equal(t, parentSpan.SpanContext.TraceID, span.SpanContext.TraceID)
				}

				spanLogs := span.Logs()
				if tc.mockRespError != nil {
					lf := spanLogs[0].Fields[0]
					assert.Equal(t, "grpc.error", lf.Key)
					assert.Equal(t, tc.mockRespError.Error(), lf.ValueString)
				}
			}
		})
	}
}
