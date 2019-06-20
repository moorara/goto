package trace

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/uber/jaeger-client-go/config"
)

func TestJaegerLogger(t *testing.T) {
	tests := []struct {
		errorMsg         string
		infoMsg          string
		infoArgs         []interface{}
		expectedErrorMsg string
		expectedInfoMsg  string
	}{
		{
			errorMsg:         "test error message",
			infoMsg:          "test %s %s",
			infoArgs:         []interface{}{"info", "message"},
			expectedErrorMsg: "test error message",
			expectedInfoMsg:  "test info message",
		},
	}

	for _, tc := range tests {
		// Logger with pipe to read from
		rd, wr, _ := os.Pipe()
		dec := json.NewDecoder(rd)
		logger := log.NewJSONLogger(wr)

		jlogger := &jaegerLogger{logger}
		jlogger.Error(tc.errorMsg)
		jlogger.Infof(tc.infoMsg, tc.infoArgs...)

		var log map[string]interface{}

		// Verify Error
		err := dec.Decode(&log)
		assert.NoError(t, err)
		assert.Equal(t, "error", log["level"])
		assert.Equal(t, tc.expectedErrorMsg, log["message"])

		// Verify Infof
		err = dec.Decode(&log)
		assert.NoError(t, err)
		assert.Equal(t, "info", log["level"])
		assert.Equal(t, tc.expectedInfoMsg, log["message"])
	}
}

func TestNewConstSampler(t *testing.T) {
	tests := []struct {
		enabled bool
	}{
		{true},
		{false},
	}

	for _, tc := range tests {
		sampler := NewConstSampler(tc.enabled)
		assert.NotNil(t, sampler)
	}
}

func TestNewProbabilisticSampler(t *testing.T) {
	tests := []struct {
		probability float64
	}{
		{0.0},
		{0.5},
		{1.0},
	}

	for _, tc := range tests {
		sampler := NewProbabilisticSampler(tc.probability)
		assert.NotNil(t, sampler)
	}
}

func TestNewRateLimitingSampler(t *testing.T) {
	tests := []struct {
		rate float64
	}{
		{0.1},
		{1.0},
		{10},
	}

	for _, tc := range tests {
		sampler := NewRateLimitingSampler(tc.rate)
		assert.NotNil(t, sampler)
	}
}

func TestNewRemoteSampler(t *testing.T) {
	tests := []struct {
		initial   float64
		serverURL string
		interval  time.Duration
	}{
		{
			0.5,
			"http://jaeger-agent:5778/sampling",
			time.Minute,
		},
	}

	for _, tc := range tests {
		sampler := NewRemoteSampler(tc.initial, tc.serverURL, tc.interval)
		assert.NotNil(t, sampler)
	}
}

func TestNewAgentReporter(t *testing.T) {
	tests := []struct {
		agentAddr string
		logSpans  bool
	}{
		{
			"jaeger-agent:6831",
			false,
		},
		{
			"jaeger-agent:6831",
			true,
		},
	}

	for _, tc := range tests {
		reporter := NewAgentReporter(tc.agentAddr, tc.logSpans)
		assert.NotNil(t, reporter)
	}
}

func TestNewCollectorReporter(t *testing.T) {
	tests := []struct {
		collectorAddr string
		logSpans      bool
	}{
		{
			"http://jaeger-collector:5678/api/traces",
			false,
		},
		{
			"http://jaeger-collector:5678/api/traces",
			true,
		},
	}

	for _, tc := range tests {
		reporter := NewCollectorReporter(tc.collectorAddr, tc.logSpans)
		assert.NotNil(t, reporter)
	}
}

func TestNewTracer(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		sampler     *config.SamplerConfig
		reporter    *config.ReporterConfig
		logger      log.Logger
		reg         prometheus.Registerer
	}{
		{
			"WithLoggingAndMetrics",
			"service_name",
			&config.SamplerConfig{},
			&config.ReporterConfig{},
			nil,
			prometheus.NewRegistry(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracer, closer, err := NewTracer(tc.serviceName, tc.sampler, tc.reporter, tc.logger, tc.reg)
			assert.NoError(t, err)

			defer closer.Close()

			assert.NotNil(t, tracer)
			assert.NotNil(t, closer)
		})
	}
}
