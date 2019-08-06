// Package trace can be used for creating Jaeger tracers
// Deprecated: this package has been frozen and deprecated in favor of https://github.com/moorara/observe/trace
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
	level.Error(l.logger).Log("message", msg)
}

func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
	level.Info(l.logger).Log("message", fmt.Sprintf(msg, args...))
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

// Options contains optional options for Tracer
type Options struct {
	Name     string
	Sampler  *jconfig.SamplerConfig
	Reporter *jconfig.ReporterConfig
	Logger   log.Logger
	PromReg  prometheus.Registerer
}

// NewTracer creates a new tracer
func NewTracer(opts Options) (opentracing.Tracer, io.Closer, error) {
	if opts.Name == "" {
		opts.Name = "tracer"
	}

	if opts.Sampler == nil {
		opts.Sampler = NewConstSampler(true)
	}

	if opts.Reporter == nil {
		opts.Reporter = NewAgentReporter("localhost:6831", false)
	}

	jgOpts := []jconfig.Option{}
	jgConfig := &jconfig.Configuration{
		ServiceName: opts.Name,
		Sampler:     opts.Sampler,
		Reporter:    opts.Reporter,
	}

	if opts.Logger != nil {
		jlogger := &jaegerLogger{opts.Logger}
		loggerOpt := jconfig.Logger(jlogger)
		jgOpts = append(jgOpts, loggerOpt)
	}

	if opts.PromReg != nil {
		regOpt := jprometheus.WithRegisterer(opts.PromReg)
		factory := jprometheus.New(regOpt).Namespace(jmetrics.NSOptions{Name: opts.Name})
		metricsOpt := jconfig.Metrics(factory)
		jgOpts = append(jgOpts, metricsOpt)
	}

	return jgConfig.NewTracer(jgOpts...)
}
