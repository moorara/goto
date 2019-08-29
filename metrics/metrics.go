// Package metrics can be used for creating Prometheus metrics
//
// Deprecated: this package has been frozen and deprecated in favor of github.com/moorara/observe/metrics
package metrics

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultBuckets   = []float64{0.01, 0.10, 0.50, 1.00, 5.00}
	defaultQuantiles = map[float64]float64{
		0.1:  0.1,
		0.5:  0.05,
		0.95: 0.01,
		0.99: 0.001,
	}
)

type (
	// FactoryOptions contains optional options for creating a Factory
	FactoryOptions struct {
		Prefix     string
		Buckets    []float64
		Quantiles  map[float64]float64
		Registerer prometheus.Registerer
	}

	// Factory is used for creating new metrics with consistent settings
	Factory struct {
		prefix     string
		buckets    []float64
		quantiles  map[float64]float64
		registerer prometheus.Registerer
	}

	// OpMetrics includes metrics for internal operations
	OpMetrics struct {
		OpLatencyHist *prometheus.HistogramVec
		OpLatencySumm *prometheus.SummaryVec
	}

	// RequestMetrics includes metrics for service requests
	RequestMetrics struct {
		ReqCounter      *prometheus.CounterVec
		ReqGauge        *prometheus.GaugeVec
		ReqDurationHist *prometheus.HistogramVec
		ReqDurationSumm *prometheus.SummaryVec
	}
)

// NewFactory creates a new instance of Factory
func NewFactory(opts FactoryOptions) *Factory {
	if opts.Buckets == nil || len(opts.Buckets) == 0 {
		opts.Buckets = defaultBuckets
	}

	if opts.Quantiles == nil || len(opts.Quantiles) == 0 {
		opts.Quantiles = defaultQuantiles
	}

	if opts.Registerer == nil {
		opts.Registerer = prometheus.DefaultRegisterer
	}

	// GoCollector and ProcessCollector are registered with default Prometheus registry by default
	if opts.Registerer != prometheus.DefaultRegisterer {
		opts.Registerer.MustRegister(prometheus.NewGoCollector())
		opts.Registerer.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	}

	return &Factory{
		prefix:     opts.Prefix,
		buckets:    opts.Buckets,
		quantiles:  opts.Quantiles,
		registerer: opts.Registerer,
	}
}

func (f *Factory) getMetricName(name string) string {
	if f.prefix != "" {
		name = fmt.Sprintf("%s_%s", f.prefix, name)
	}

	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, "-", "_", -1)

	return name
}

// Counter creates a new counter metrics
func (f *Factory) Counter(name, description string, labels []string) *prometheus.CounterVec {
	opts := prometheus.CounterOpts{
		Name: f.getMetricName(name),
		Help: description,
	}

	counter := prometheus.NewCounterVec(opts, labels)
	f.registerer.MustRegister(counter)

	return counter
}

// Gauge creates a new gauge metrics
func (f *Factory) Gauge(name, description string, labels []string) *prometheus.GaugeVec {
	opts := prometheus.GaugeOpts{
		Name: f.getMetricName(name),
		Help: description,
	}

	gauge := prometheus.NewGaugeVec(opts, labels)
	f.registerer.MustRegister(gauge)

	return gauge
}

// Histogram creates a new histogram metrics
func (f *Factory) Histogram(name, description string, labels []string) *prometheus.HistogramVec {
	opts := prometheus.HistogramOpts{
		Name:    f.getMetricName(name),
		Help:    description,
		Buckets: defaultBuckets,
	}

	histogram := prometheus.NewHistogramVec(opts, labels)
	f.registerer.MustRegister(histogram)

	return histogram
}

// Summary creates a new summary metrics
func (f *Factory) Summary(name, description string, labels []string) *prometheus.SummaryVec {
	opts := prometheus.SummaryOpts{
		Name:       f.getMetricName(name),
		Help:       description,
		Objectives: defaultQuantiles,
	}

	summary := prometheus.NewSummaryVec(opts, labels)
	f.registerer.MustRegister(summary)

	return summary
}
