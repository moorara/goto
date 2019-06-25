package grpc

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/metadata"
)

type mockSpan struct {
	FinishCalled bool

	FinishWithOptionsInOpts opentracing.FinishOptions

	ContextOutSpanContext opentracing.SpanContext

	SetOperationNameInOpName string
	SetOperationNameOutSpan  opentracing.Span

	SetTagInKey   string
	SetTagInValue interface{}
	SetTagOutSpan opentracing.Span

	LogFieldsInFields []opentracingLog.Field

	LogKVInAltKeyValues []interface{}

	SetBaggageItemInRestrictedKey string
	SetBaggageItemInValue         string
	SetBaggageItemOutSpan         opentracing.Span

	BaggageItemInRestrictedKey string
	BaggageItemOutResult       string

	TracerOutTracer opentracing.Tracer

	LogEventInEvent string

	LogEventWithPayloadInEvent   string
	LogEventWithPayloadInPayload interface{}

	LogInData opentracing.LogData
}

func (m *mockSpan) Finish() {
	m.FinishCalled = true
}

func (m *mockSpan) FinishWithOptions(opts opentracing.FinishOptions) {
	m.FinishWithOptionsInOpts = opts
}

func (m *mockSpan) Context() opentracing.SpanContext {
	return m.ContextOutSpanContext
}

func (m *mockSpan) SetOperationName(opName string) opentracing.Span {
	m.SetOperationNameInOpName = opName
	return m.SetOperationNameOutSpan
}

func (m *mockSpan) SetTag(key string, value interface{}) opentracing.Span {
	m.SetTagInKey = key
	m.SetTagInValue = value
	return m.SetTagOutSpan
}

func (m *mockSpan) LogFields(fields ...opentracingLog.Field) {
	m.LogFieldsInFields = fields
}

func (m *mockSpan) LogKV(altKeyValues ...interface{}) {
	m.LogKVInAltKeyValues = altKeyValues
}

func (m *mockSpan) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	m.SetBaggageItemInRestrictedKey = restrictedKey
	m.SetBaggageItemInValue = value
	return m.SetBaggageItemOutSpan
}

func (m *mockSpan) BaggageItem(restrictedKey string) string {
	m.BaggageItemInRestrictedKey = restrictedKey
	return m.BaggageItemOutResult
}

func (m *mockSpan) Tracer() opentracing.Tracer {
	return m.TracerOutTracer
}

func (m *mockSpan) LogEvent(event string) {
	m.LogEventInEvent = event
}

func (m *mockSpan) LogEventWithPayload(event string, payload interface{}) {
	m.LogEventWithPayloadInEvent = event
	m.LogEventWithPayloadInPayload = payload
}

func (m *mockSpan) Log(data opentracing.LogData) {
	m.LogInData = data
}

type mockServerStream struct {
	SetHeaderInMD     metadata.MD
	SetHeaderOutError error

	SendHeaderInMD     metadata.MD
	SendHeaderOutError error

	SetTrailerInMD metadata.MD

	ContextOutContext context.Context

	SendMsgInMsg    interface{}
	SendMsgOutError error

	RecvMsgInMsg    interface{}
	RecvMsgOutError error
}

func (m *mockServerStream) SetHeader(md metadata.MD) error {
	m.SetHeaderInMD = md
	return m.SetHeaderOutError
}

func (m *mockServerStream) SendHeader(md metadata.MD) error {
	m.SendHeaderInMD = md
	return m.SendHeaderOutError
}

func (m *mockServerStream) SetTrailer(md metadata.MD) {
	m.SetTrailerInMD = md
}

func (m *mockServerStream) Context() context.Context {
	return m.ContextOutContext
}

func (m *mockServerStream) SendMsg(msg interface{}) error {
	m.SendMsgInMsg = msg
	return m.SendMsgOutError
}

func (m *mockServerStream) RecvMsg(msg interface{}) error {
	m.RecvMsgInMsg = msg
	return m.RecvMsgOutError
}
