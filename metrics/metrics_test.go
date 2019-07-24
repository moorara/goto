package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	model "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	registry := prometheus.NewRegistry()

	tests := []struct {
		name               string
		opts               FactoryOptions
		expectedPrefix     string
		expectedBuckets    []float64
		expectedQuantiles  map[float64]float64
		expectedRegisterer prometheus.Registerer
	}{
		{
			name:               "Defaults",
			opts:               FactoryOptions{},
			expectedPrefix:     "",
			expectedBuckets:    defaultBuckets,
			expectedQuantiles:  defaultQuantiles,
			expectedRegisterer: prometheus.DefaultRegisterer,
		},
		{
			name: "WithPrefix",
			opts: FactoryOptions{
				Prefix: "service_name",
			},
			expectedPrefix:     "service_name",
			expectedBuckets:    defaultBuckets,
			expectedQuantiles:  defaultQuantiles,
			expectedRegisterer: prometheus.DefaultRegisterer,
		},
		{
			name: "WithBuckets",
			opts: FactoryOptions{
				Buckets: []float64{0.01, 0.10, 0.50, 1.00, 5.00},
			},
			expectedPrefix:     "",
			expectedBuckets:    []float64{0.01, 0.10, 0.50, 1.00, 5.00},
			expectedQuantiles:  defaultQuantiles,
			expectedRegisterer: prometheus.DefaultRegisterer,
		},
		{
			name: "WithQuantiles",
			opts: FactoryOptions{
				Quantiles: map[float64]float64{
					0.1:  0.1,
					0.95: 0.01,
					0.99: 0.001,
				},
				Registerer: nil,
			},
			expectedPrefix:  "",
			expectedBuckets: defaultBuckets,
			expectedQuantiles: map[float64]float64{
				0.1:  0.1,
				0.95: 0.01,
				0.99: 0.001,
			},
			expectedRegisterer: prometheus.DefaultRegisterer,
		},
		{
			name: "WithRegistry",
			opts: FactoryOptions{
				Registerer: registry,
			},
			expectedPrefix:     "",
			expectedBuckets:    defaultBuckets,
			expectedQuantiles:  defaultQuantiles,
			expectedRegisterer: registry,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.opts)

			assert.Equal(t, tc.expectedPrefix, mf.prefix)
			assert.Equal(t, tc.expectedBuckets, mf.buckets)
			assert.Equal(t, tc.expectedQuantiles, mf.quantiles)
			assert.Equal(t, tc.expectedRegisterer, mf.registerer)
		})
	}
}

func TestCounter(t *testing.T) {
	tests := []struct {
		name         string
		opts         FactoryOptions
		metricName   string
		description  string
		labels       []string
		labelValues  []string
		addValue     float64
		expectedName string
	}{
		{
			name:         "Defaults",
			opts:         FactoryOptions{},
			metricName:   "counter_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			addValue:     2,
			expectedName: "counter_metric_name",
		},
		{
			name: "WithPrefix",
			opts: FactoryOptions{
				Prefix: "service-name",
			},
			metricName:   "counter_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			addValue:     2,
			expectedName: "service_name_counter_metric_name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.opts)
			counter := mf.Counter(tc.metricName, tc.description, tc.labels)

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
		})
	}
}

func TestGauge(t *testing.T) {
	tests := []struct {
		name         string
		opts         FactoryOptions
		metricName   string
		description  string
		labels       []string
		labelValues  []string
		addValue     float64
		subValue     float64
		expectedName string
	}{
		{
			name:         "Defaults",
			opts:         FactoryOptions{},
			metricName:   "gauge_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			addValue:     2,
			subValue:     2,
			expectedName: "gauge_metric_name",
		},
		{
			name: "WithPrefix",
			opts: FactoryOptions{
				Prefix: "service-name",
			},
			metricName:   "gauge_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			addValue:     2,
			subValue:     2,
			expectedName: "service_name_gauge_metric_name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.opts)
			gauge := mf.Gauge(tc.metricName, tc.description, tc.labels)

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
		})
	}
}

func TestHistogram(t *testing.T) {
	tests := []struct {
		name         string
		opts         FactoryOptions
		metricName   string
		description  string
		labels       []string
		labelValues  []string
		value        float64
		expectedName string
	}{
		{
			name:         "Defaults",
			opts:         FactoryOptions{},
			metricName:   "histogram_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			value:        0.1234,
			expectedName: "histogram_metric_name",
		},
		{
			name: "WithPrefix",
			opts: FactoryOptions{
				Prefix: "service-name",
			},
			metricName:   "histogram_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			value:        0.1234,
			expectedName: "service_name_histogram_metric_name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.opts)
			histogram := mf.Histogram(tc.metricName, tc.description, tc.labels)

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
		})
	}
}

func TestSummary(t *testing.T) {
	tests := []struct {
		name         string
		opts         FactoryOptions
		metricName   string
		description  string
		labels       []string
		labelValues  []string
		value        float64
		expectedName string
	}{
		{
			name:         "Defaults",
			opts:         FactoryOptions{},
			metricName:   "summary_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			value:        0.1234,
			expectedName: "summary_metric_name",
		},
		{
			name: "WithPrefix",
			opts: FactoryOptions{
				Prefix: "service-name",
			},
			metricName:   "summary_metric_name",
			description:  "metric description",
			labels:       []string{"environment", "region"},
			labelValues:  []string{"prodcution", "us-east-1"},
			value:        0.1234,
			expectedName: "service_name_summary_metric_name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mf := NewFactory(tc.opts)
			summary := mf.Summary(tc.metricName, tc.description, tc.labels)

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
		})
	}
}
