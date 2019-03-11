# metrics

This is a *helper* package for creating consistent [**Prometheus**](https://prometheus.io) metrics.

## Quick Start

For creating new metrics using default *registry*, **buckets**, and **quantiles**:

```go
package main

import (
	"log"
  "net/http"

  "github.com/moorara/goto/metrics"
	"github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
  mf := metrics.NewFactory("hello_service", nil, nil)

  // Create a histogram metric
  histogram := mf.Histogram("histogram_metric_name", "metric description", []string{"environment", "region"})
  prometheus.MustRegister(histogram)
  histogram.WithLabelValues("prodcution", "us-east-1").Observe(0.1234)

  // Expose metrics via /metrics endpoint and an HTTP server
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

For creating new metrics using a new *registry* and custom **buckets** and **quantiles**:

```go
package main

import (
	"log"
  "net/http"

  "github.com/moorara/goto/metrics"
	"github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
  registry := prometheus.NewRegistry()
  mf := metrics.NewFactory("hello_service", []float64{0.01, 0.10, 0.50, 1.00, 5.00}, map[float64]float64{
    0.1:  0.1,
    0.95: 0.01,
    0.99: 0.001,
  })

  // Create default system metrics
  sys := mf.SystemMetrics()
  registry.MustRegister(sys.Go)
  registry.MustRegister(sys.Process)

  // Create a histogram metric
  histogram := mf.Histogram("histogram_metric_name", "metric description", []string{"environment", "region"})
  registry.MustRegister(histogram)
  histogram.WithLabelValues("prodcution", "us-east-1").Observe(0.1234)

  // Expose metrics via /metrics endpoint and an HTTP server
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
