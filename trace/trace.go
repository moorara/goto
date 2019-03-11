package trace

import (
	"fmt"
	"io"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	opentracing "github.com/opentracing/opentracing-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	jmetrics "github.com/uber/jaeger-lib/metrics"
	jprometheus "github.com/uber/jaeger-lib/metrics/prometheus"
)

// jaegerLogger implements jaeger.Logger
type jaegerLogger struct {
	logger log.Logger
}

func (l *jaegerLogger) Error(msg string) {
	if l.logger != nil {
		level.Error(l.logger).Log("message", msg)
	}
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	if l.logger != nil {
		level.Info(l.logger).Log("message", fmt.Sprintf(msg, args...))
	}
}

// NewConstSampler creates a constant Jaeger sampler
//   enabled true will report all traces
//   enabled false will skip all traces
func NewConstSampler(enabled bool) *jconfig.SamplerConfig {
	var param float64
	if enabled {
		param = 1
	}

	return &jconfig.SamplerConfig{
		Type:  "const",
		Param: param,
	}
}

// NewProbabilisticSampler creates a probabilistic Jaeger sampler
//   probability is between 0 and 1
func NewProbabilisticSampler(probability float64) *jconfig.SamplerConfig {
	return &jconfig.SamplerConfig{
		Type:  "probabilistic",
		Param: probability,
	}
}

// NewRateLimitingSampler creates a rate limited Jaeger sampler
//   rate is the number of spans per second
func NewRateLimitingSampler(rate float64) *jconfig.SamplerConfig {
	return &jconfig.SamplerConfig{
		Type:  "rateLimiting",
		Param: rate,
	}
}

// NewRemoteSampler creates a Jaeger sampler pulling remote sampling strategies
//   probability is the initial probability between 0 and 1 before a remote sampling strategy is recieved
//   serverURL is the address of sampling server
//   interval specifies the rate of polling remote sampling strategies
func NewRemoteSampler(probability float64, serverURL string, interval time.Duration) *jconfig.SamplerConfig {
	return &jconfig.SamplerConfig{
		Type:                    "remote",
		Param:                   probability,
		SamplingServerURL:       serverURL,
		SamplingRefreshInterval: interval,
	}
}

// NewAgentReporter creates a Jaeger reporter reporting to jaeger-agent
//   agentAddr is the address of Jaeger agent
//   logSpans true will log all spans
func NewAgentReporter(agentAddr string, logSpans bool) *jconfig.ReporterConfig {
	return &jconfig.ReporterConfig{
		LocalAgentHostPort: agentAddr,
		LogSpans:           logSpans,
	}
}

// NewCollectorReporter creates a Jaeger reporter reporting to jaeger-collector
//   collectorAddr is the address of Jaeger collector
//   logSpans true will log all spans
func NewCollectorReporter(collectorAddr string, logSpans bool) *jconfig.ReporterConfig {
	return &jconfig.ReporterConfig{
		CollectorEndpoint: collectorAddr,
		LogSpans:          logSpans,
	}
}

// NewTracer creates a new tracer
func NewTracer(name string, sampler *jconfig.SamplerConfig, reporter *jconfig.ReporterConfig, logger log.Logger, reg prometheus.Registerer) (opentracing.Tracer, io.Closer, error) {
	jgConfig := &jconfig.Configuration{
		ServiceName: name,
		Sampler:     sampler,
		Reporter:    reporter,
	}

	jlogger := &jaegerLogger{logger}
	loggerOpt := jconfig.Logger(jlogger)

	regOpt := jprometheus.WithRegisterer(reg)
	factory := jprometheus.New(regOpt).Namespace(jmetrics.NSOptions{Name: name})
	metricsOpt := jconfig.Metrics(factory)

	return jgConfig.NewTracer(loggerOpt, metricsOpt)
}
