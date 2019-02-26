package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultBuckets   = []float64{0.01, 0.10, 0.50, 1.00}
	defaultQuantiles = map[float64]float64{
		0.1:  0.1,
		0.5:  0.05,
		0.95: 0.01,
		0.99: 0.001,
	}
)

type (
	// Factory creates new metrics factory
	Factory struct {
		namespace string
		buckets   []float64
		quantiles map[float64]float64
	}

	// SystemMetrics includes system metrics
	SystemMetrics struct {
		Go      prometheus.Collector
		Process prometheus.Collector
	}

	// OpMetrics includes metrics for internal operations
	OpMetrics struct {
		OpLatencyHist *prometheus.HistogramVec
		OpLatencySumm *prometheus.SummaryVec
	}

	// RequestMetrics includes metrics for service requests
	RequestMetrics struct {
		ReqGauge        *prometheus.GaugeVec
		ReqCounter      *prometheus.CounterVec
		ReqDurationHist *prometheus.HistogramVec
		ReqDurationSumm *prometheus.SummaryVec
	}
)

// NewFactory creates a new instance of Factory
func NewFactory(namespace string, buckets []float64, quantiles map[float64]float64) *Factory {
	namespace = strings.Replace(namespace, " ", "_", -1)
	namespace = strings.Replace(namespace, "-", "_", -1)

	if buckets == nil || len(buckets) == 0 {
		buckets = defaultBuckets
	}

	if quantiles == nil || len(quantiles) == 0 {
		quantiles = defaultQuantiles
	}

	return &Factory{
		namespace: namespace,
		buckets:   buckets,
		quantiles: quantiles,
	}
}

// Counter creates a new counter metrics
func (f *Factory) Counter(name, description string, labels []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: f.namespace,
			Name:      name,
			Help:      description,
		},
		labels,
	)
}

// Gauge creates a new gauge metrics
func (f *Factory) Gauge(name, description string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: f.namespace,
			Name:      name,
			Help:      description,
		},
		labels,
	)
}

// Histogram creates a new histogram metrics
func (f *Factory) Histogram(name, description string, labels []string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: f.namespace,
			Name:      name,
			Help:      description,
			Buckets:   defaultBuckets,
		},
		labels,
	)
}

// Summary creates a new summary metrics
func (f *Factory) Summary(name, description string, labels []string) *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  f.namespace,
			Name:       name,
			Help:       description,
			Objectives: defaultQuantiles,
		},
		labels,
	)
}

// SystemMetrics creates system metrics
func (f *Factory) SystemMetrics() *SystemMetrics {
	return &SystemMetrics{
		Go: prometheus.NewGoCollector(),
		Process: prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{
			Namespace: f.namespace,
		}),
	}
}

// OpMetrics creates metrics for internal operations
func (f *Factory) OpMetrics() *OpMetrics {
	return &OpMetrics{
		OpLatencyHist: f.Histogram("operations_latency_seconds", "latency of internal operations", []string{"op", "success"}),
		OpLatencySumm: f.Summary("operations_latency_quantiles_seconds", "latency quantiles of internal operations", []string{"op", "success"}),
	}
}

// HTTPRequestMetrics creates metrics for HTTP requests
func (f *Factory) HTTPRequestMetrics() *RequestMetrics {
	return &RequestMetrics{
		ReqGauge:        f.Gauge("http_requests", "current number of http requests", []string{"method", "endpoint"}),
		ReqCounter:      f.Counter("http_requests_total", "total number of http requests", []string{"method", "endpoint", "success"}),
		ReqDurationHist: f.Histogram("http_request_duration_seconds", "duration of http requests", []string{"method", "endpoint", "statusCode", "statusClass"}),
		ReqDurationSumm: f.Summary("http_request_duration_quantiles_seconds", "duration quantiles of http requests", []string{"method", "endpoint", "statusCode", "statusClass"}),
	}
}

// GRPCRequestMetrics creates metrics for gRPC requests
func (f *Factory) GRPCRequestMetrics() *RequestMetrics {
	return &RequestMetrics{
		ReqGauge:        f.Gauge("grpc_requests", "current number of grpc requests", []string{"method"}),
		ReqCounter:      f.Counter("grpc_requests_total", "total number of grpc requests", []string{"method", "success"}),
		ReqDurationHist: f.Histogram("grpc_request_duration_seconds", "duration of grpc requests", []string{"method", "success"}),
		ReqDurationSumm: f.Summary("grpc_request_duration_quantiles_seconds", "duration quantiles of grpc requests", []string{"method", "success"}),
	}
}
