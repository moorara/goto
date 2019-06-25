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

func injectSpan(ctx context.Context, tracer opentracing.Tracer, span opentracing.Span) context.Context {
	// Get any metadata if set
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		md = md.Copy()
	} else {
		md = metadata.New(nil)
	}

	carrier := &MetadataTextMap{md}
	err := tracer.Inject(span.Context(), opentracing.TextMap, carrier)
	if err != nil {
		return ctx
	}

	return metadata.NewIncomingContext(ctx, md)
}

func TestNewServerObservabilityInterceptor(t *testing.T) {
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
			soi := NewServerObservabilityInterceptor(tc.logger, tc.mf, tc.tracer)

			assert.Equal(t, tc.logger, soi.logger)
			assert.NotNil(t, soi.metrics)
			assert.Equal(t, tc.tracer, soi.tracer)
		})
	}
}

func TestUnaryServerInterceptor(t *testing.T) {
	tests := []struct {
		name            string
		parentSpan      opentracing.Span
		ctx             context.Context
		req             interface{}
		info            *grpc.UnaryServerInfo
		mockDelay       time.Duration
		mockRespError   error
		mockRespRes     interface{}
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
			req:             nil,
			info:            &grpc.UnaryServerInfo{FullMethod: ""},
			mockDelay:       0,
			mockRespError:   nil,
			mockRespRes:     nil,
			verify:          false,
			expectedPackage: "",
			expectedService: "",
			expectedMethod:  "",
			expectedStream:  "",
			expectedSuccess: false,
		},
		{
			name:            "HandlerFails",
			parentSpan:      nil,
			ctx:             context.Background(),
			req:             nil,
			info:            &grpc.UnaryServerInfo{FullMethod: "/package.service/method"},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   errors.New("error on grpc method"),
			mockRespRes:     nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "false",
			expectedSuccess: false,
		},
		{
			name:            "HandlerSucceeds",
			parentSpan:      nil,
			ctx:             context.Background(),
			req:             nil,
			info:            &grpc.UnaryServerInfo{FullMethod: "/package.service/method"},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			mockRespRes:     nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "false",
			expectedSuccess: true,
		},
		{
			name:            "HandlerSucceedsWithParentSpan",
			parentSpan:      mocktracer.New().StartSpan("parent-span"),
			ctx:             context.Background(),
			req:             nil,
			info:            &grpc.UnaryServerInfo{FullMethod: "/package.service/method"},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			mockRespRes:     nil,
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
			var insertedSpan opentracing.Span

			logger := log.NewLogger(log.Options{Writer: buff})
			promReg := prometheus.NewRegistry()
			mf := metrics.NewFactory(metrics.FactoryOptions{Registerer: promReg})
			tracer := mocktracer.New()
			i := NewServerObservabilityInterceptor(logger, mf, tracer)

			// Inject the parent span context if any
			if tc.parentSpan != nil {
				tc.ctx = injectSpan(tc.ctx, tracer, tc.parentSpan)
			}

			handler := func(ctx context.Context, req interface{}) (interface{}, error) {
				time.Sleep(tc.mockDelay)
				insertedSpan = opentracing.SpanFromContext(ctx)
				return tc.mockRespRes, tc.mockRespError
			}

			res, err := i.UnaryInterceptor(tc.ctx, tc.req, tc.info, handler)
			assert.Equal(t, tc.mockRespError, err)
			assert.Equal(t, tc.mockRespRes, res)

			if tc.verify {
				// Verify logs

				var log map[string]interface{}
				err := json.NewDecoder(buff).Decode(&log)
				assert.NoError(t, err)
				assert.Equal(t, serverKind, log["grpc.kind"])
				assert.Equal(t, tc.expectedPackage, log["grpc.package"])
				assert.Equal(t, tc.expectedService, log["grpc.service"])
				assert.Equal(t, tc.expectedMethod, log["grpc.method"])
				assert.Equal(t, tc.expectedStream, log["grpc.stream"])
				assert.Equal(t, tc.expectedSuccess, log["grpc.success"])
				assert.NotEmpty(t, log["responseTime"])
				assert.NotEmpty(t, log["message"])

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
					case serverGaugeMetricName:
						assert.Equal(t, promModel.MetricType_GAUGE, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					case serverCounterMetricName:
						assert.Equal(t, promModel.MetricType_COUNTER, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					case serverHistogramMetricName:
						assert.Equal(t, promModel.MetricType_HISTOGRAM, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					case serverSummaryMetricName:
						assert.Equal(t, promModel.MetricType_SUMMARY, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					}
				}

				// Verify traces

				span := tracer.FinishedSpans()[0]
				assert.Equal(t, insertedSpan, span)
				assert.Equal(t, serverSpanName, span.OperationName)
				assert.Equal(t, ext.SpanKindEnum("server"), span.Tag("span.kind"))
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
			}
		})
	}
}

func TestStreamServerInterceptor(t *testing.T) {
	tests := []struct {
		name            string
		parentSpan      opentracing.Span
		srv             interface{}
		ss              *mockServerStream
		info            *grpc.StreamServerInfo
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
			srv:             nil,
			ss:              &mockServerStream{ContextOutContext: context.Background()},
			info:            &grpc.StreamServerInfo{FullMethod: ""},
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
			name:            "HandlerFails",
			parentSpan:      nil,
			srv:             nil,
			ss:              &mockServerStream{ContextOutContext: context.Background()},
			info:            &grpc.StreamServerInfo{FullMethod: "/package.service/method"},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   errors.New("error on grpc method"),
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "true",
			expectedSuccess: false,
		},
		{
			name:            "HandlerSucceeds",
			parentSpan:      nil,
			srv:             nil,
			ss:              &mockServerStream{ContextOutContext: context.Background()},
			info:            &grpc.StreamServerInfo{FullMethod: "/package.service/method"},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
			verify:          true,
			expectedPackage: "package",
			expectedService: "service",
			expectedMethod:  "method",
			expectedStream:  "true",
			expectedSuccess: true,
		},
		{
			name:            "HandlerSucceedsWithParentSpan",
			parentSpan:      mocktracer.New().StartSpan("parent-span"),
			srv:             nil,
			ss:              &mockServerStream{ContextOutContext: context.Background()},
			info:            &grpc.StreamServerInfo{FullMethod: "/package.service/method"},
			mockDelay:       10 * time.Millisecond,
			mockRespError:   nil,
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
			var insertedSpan opentracing.Span

			logger := log.NewLogger(log.Options{Writer: buff})
			promReg := prometheus.NewRegistry()
			mf := metrics.NewFactory(metrics.FactoryOptions{Registerer: promReg})
			tracer := mocktracer.New()
			i := NewServerObservabilityInterceptor(logger, mf, tracer)

			// Inject the parent span context if any
			if tc.parentSpan != nil {
				tc.ss.ContextOutContext = injectSpan(tc.ss.ContextOutContext, tracer, tc.parentSpan)
			}

			handler := func(srv interface{}, stream grpc.ServerStream) error {
				time.Sleep(tc.mockDelay)
				insertedSpan = opentracing.SpanFromContext(stream.Context())
				return tc.mockRespError
			}

			err := i.StreamInterceptor(tc.srv, tc.ss, tc.info, handler)
			assert.Equal(t, tc.mockRespError, err)

			if tc.verify {
				// Verify logs

				var log map[string]interface{}
				err := json.NewDecoder(buff).Decode(&log)
				assert.NoError(t, err)
				assert.Equal(t, serverKind, log["grpc.kind"])
				assert.Equal(t, tc.expectedPackage, log["grpc.package"])
				assert.Equal(t, tc.expectedService, log["grpc.service"])
				assert.Equal(t, tc.expectedMethod, log["grpc.method"])
				assert.Equal(t, tc.expectedStream, log["grpc.stream"])
				assert.Equal(t, tc.expectedSuccess, log["grpc.success"])
				assert.NotEmpty(t, log["responseTime"])
				assert.NotEmpty(t, log["message"])

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
					case serverGaugeMetricName:
						assert.Equal(t, promModel.MetricType_GAUGE, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					case serverCounterMetricName:
						assert.Equal(t, promModel.MetricType_COUNTER, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					case serverHistogramMetricName:
						assert.Equal(t, promModel.MetricType_HISTOGRAM, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					case serverSummaryMetricName:
						assert.Equal(t, promModel.MetricType_SUMMARY, *metricFamily.Type)
						verifyLabels(metricFamily.Metric[0].Label)
					}
				}

				// Verify traces

				span := tracer.FinishedSpans()[0]
				assert.Equal(t, insertedSpan, span)
				assert.Equal(t, serverSpanName, span.OperationName)
				assert.Equal(t, ext.SpanKindEnum("server"), span.Tag("span.kind"))
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
			}
		})
	}
}
