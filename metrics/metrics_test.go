package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	model "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	tests := []struct {
		namespace         string
		buckets           []float64
		quantiles         map[float64]float64
		expectedNamespace string
		expectedBuckets   []float64
		expectedQuantiles map[float64]float64
	}{
		{
			namespace:         "service_name",
			buckets:           nil,
			quantiles:         nil,
			expectedNamespace: "service_name",
			expectedBuckets:   defaultBuckets,
			expectedQuantiles: defaultQuantiles,
		},
		{
			namespace:         "service-name",
			buckets:           []float64{0.01, 0.10, 0.50, 1.00, 5.00},
			quantiles:         nil,
			expectedNamespace: "service_name",
			expectedBuckets:   []float64{0.01, 0.10, 0.50, 1.00, 5.00},
			expectedQuantiles: defaultQuantiles,
		},
		{
			namespace: "service name",
			buckets:   []float64{0.01, 0.10, 0.50, 1.00, 5.00},
			quantiles: map[float64]float64{
				0.1:  0.1,
				0.95: 0.01,
				0.99: 0.001,
			},
			expectedNamespace: "service_name",
			expectedBuckets:   []float64{0.01, 0.10, 0.50, 1.00, 5.00},
			expectedQuantiles: map[float64]float64{
				0.1:  0.1,
				0.95: 0.01,
				0.99: 0.001,
			},
		},
	}

	for _, tc := range tests {
		mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)

		assert.Equal(t, tc.expectedNamespace, mf.namespace)
		assert.Equal(t, tc.expectedBuckets, mf.buckets)
		assert.Equal(t, tc.expectedQuantiles, mf.quantiles)
	}
}

func TestCounter(t *testing.T) {
	tests := []struct {
		namespace    string
		buckets      []float64
		quantiles    map[float64]float64
		name         string
		description  string
		labels       []string
		labelValues  []string
		addValue     float64
		expectedName string
	}{
		{
			namespace:    "service_name",
			buckets:      nil,
			quantiles:    nil,
			name:         "counter_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			addValue:     2,
			expectedName: "service_name_counter_metric_name",
		},
	}

	for _, tc := range tests {
		mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
		counter := mf.Counter(tc.name, tc.description, tc.labels)

		reg := prometheus.NewRegistry()
		reg.MustRegister(counter)
		counter.WithLabelValues(tc.labelValues...).Inc()
		counter.WithLabelValues(tc.labelValues...).Add(tc.addValue)

		metricFamilies, err := reg.Gather()
		assert.NoError(t, err)
		for _, metricFamily := range metricFamilies {
			assert.Equal(t, tc.expectedName, *metricFamily.Name)
			assert.Equal(t, tc.description, *metricFamily.Help)
			assert.Equal(t, model.MetricType_COUNTER, *metricFamily.Type)
		}
	}
}

func TestGauge(t *testing.T) {
	tests := []struct {
		namespace    string
		buckets      []float64
		quantiles    map[float64]float64
		name         string
		description  string
		labels       []string
		labelValues  []string
		addValue     float64
		subValue     float64
		expectedName string
	}{
		{
			namespace:    "service_name",
			buckets:      nil,
			quantiles:    nil,
			name:         "gauge_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			addValue:     2,
			subValue:     2,
			expectedName: "service_name_gauge_metric_name",
		},
	}

	for _, tc := range tests {
		mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
		gauge := mf.Gauge(tc.name, tc.description, tc.labels)

		reg := prometheus.NewRegistry()
		reg.MustRegister(gauge)
		gauge.WithLabelValues(tc.labelValues...).Inc()
		gauge.WithLabelValues(tc.labelValues...).Add(tc.addValue)
		gauge.WithLabelValues(tc.labelValues...).Add(tc.subValue)

		metricFamilies, err := reg.Gather()
		assert.NoError(t, err)
		for _, metricFamily := range metricFamilies {
			assert.Equal(t, tc.expectedName, *metricFamily.Name)
			assert.Equal(t, tc.description, *metricFamily.Help)
			assert.Equal(t, model.MetricType_GAUGE, *metricFamily.Type)
		}
	}
}

func TestHistogram(t *testing.T) {
	tests := []struct {
		namespace    string
		buckets      []float64
		quantiles    map[float64]float64
		name         string
		description  string
		labels       []string
		labelValues  []string
		value        float64
		expectedName string
	}{
		{
			namespace:    "service_name",
			buckets:      nil,
			quantiles:    nil,
			name:         "histogram_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			value:        0.1234,
			expectedName: "service_name_histogram_metric_name",
		},
	}

	for _, tc := range tests {
		mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
		histogram := mf.Histogram(tc.name, tc.description, tc.labels)

		reg := prometheus.NewRegistry()
		reg.MustRegister(histogram)
		histogram.WithLabelValues(tc.labelValues...).Observe(tc.value)

		metricFamilies, err := reg.Gather()
		assert.NoError(t, err)
		for _, metricFamily := range metricFamilies {
			assert.Equal(t, tc.expectedName, *metricFamily.Name)
			assert.Equal(t, tc.description, *metricFamily.Help)
			assert.Equal(t, model.MetricType_HISTOGRAM, *metricFamily.Type)
		}
	}
}

func TestSummary(t *testing.T) {
	tests := []struct {
		namespace    string
		buckets      []float64
		quantiles    map[float64]float64
		name         string
		description  string
		labels       []string
		labelValues  []string
		value        float64
		expectedName string
	}{
		{
			namespace:    "service_name",
			buckets:      nil,
			quantiles:    nil,
			name:         "summary_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			value:        0.1234,
			expectedName: "service_name_summary_metric_name",
		},
	}

	for _, tc := range tests {
		mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
		summary := mf.Summary(tc.name, tc.description, tc.labels)

		reg := prometheus.NewRegistry()
		reg.MustRegister(summary)
		summary.WithLabelValues(tc.labelValues...).Observe(tc.value)

		metricFamilies, err := reg.Gather()
		assert.NoError(t, err)
		for _, metricFamily := range metricFamilies {
			assert.Equal(t, tc.expectedName, *metricFamily.Name)
			assert.Equal(t, tc.description, *metricFamily.Help)
			assert.Equal(t, model.MetricType_SUMMARY, *metricFamily.Type)
		}
	}
}

func TestSystemMetrics(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		buckets   []float64
		quantiles map[float64]float64
	}{
		{
			name:      "Defaults",
			namespace: "service_name",
			buckets:   nil,
			quantiles: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
			metrics := mf.SystemMetrics()

			assert.NotNil(t, metrics.Go)
			assert.NotNil(t, metrics.Process)
		})
	}
}

func TestOpMetrics(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		buckets   []float64
		quantiles map[float64]float64
	}{
		{
			name:      "Defaults",
			namespace: "service_name",
			buckets:   nil,
			quantiles: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
			metrics := mf.OpMetrics()

			assert.IsType(t, &prometheus.HistogramVec{}, metrics.OpLatencyHist)
			assert.IsType(t, &prometheus.SummaryVec{}, metrics.OpLatencySumm)
		})
	}
}

func TestHTTPRequestMetrics(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		buckets   []float64
		quantiles map[float64]float64
	}{
		{
			name:      "Defaults",
			namespace: "service_name",
			buckets:   nil,
			quantiles: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
			metrics := mf.HTTPRequestMetrics()

			assert.IsType(t, &prometheus.GaugeVec{}, metrics.ReqGauge)
			assert.IsType(t, &prometheus.CounterVec{}, metrics.ReqCounter)
			assert.IsType(t, &prometheus.HistogramVec{}, metrics.ReqDurationHist)
			assert.IsType(t, &prometheus.SummaryVec{}, metrics.ReqDurationSumm)
		})
	}
}

func TestGRPCRequestMetrics(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		buckets   []float64
		quantiles map[float64]float64
	}{
		{
			name:      "Defaults",
			namespace: "service_name",
			buckets:   nil,
			quantiles: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.namespace, tc.buckets, tc.quantiles)
			metrics := mf.GRPCRequestMetrics()

			assert.IsType(t, &prometheus.GaugeVec{}, metrics.ReqGauge)
			assert.IsType(t, &prometheus.CounterVec{}, metrics.ReqCounter)
			assert.IsType(t, &prometheus.HistogramVec{}, metrics.ReqDurationHist)
			assert.IsType(t, &prometheus.SummaryVec{}, metrics.ReqDurationSumm)
		})
	}
}
