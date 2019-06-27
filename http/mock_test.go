package http

import (
	opentracing "github.com/opentracing/opentracing-go"
	opentracingLog "github.com/opentracing/opentracing-go/log"
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
