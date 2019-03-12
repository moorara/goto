# metrics

This is a helper package for creating consistent [**Prometheus**](https://prometheus.io) metrics.

## Quick Start

For creating new metrics using default **buckets** and **quantiles** and register them with default *registry*:

```go
package main

import (
  "log"
  "net/http"

  "github.com/moorara/goto/metrics"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
  mf := metrics.NewFactory(metrics.FactoryOptions{})

  // Create a histogram metric
  histogram := mf.Histogram("auth_service_session_duration_seconds", "duration of user sessions", []string{"environment", "region"})
  histogram.WithLabelValues("prodcution", "us-east-1").Observe(60)

  // Expose metrics via /metrics endpoint and an HTTP server
  http.Handle("/metrics", promhttp.Handler())
  log.Fatal(http.ListenAndServe(":8080", nil))
}
```

For creating new metrics using custom **buckets** and **quantiles** and register them with a new *registry*:

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
  mf := metrics.NewFactory(metrics.FactoryOptions{
    Registerer: registry,
    Buckets:    []float64{0.01, 0.10, 0.50, 1.00, 5.00},
    Quantiles:  map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
  })

  // Create a histogram metric
  histogram := mf.Histogram("auth_service_session_duration_seconds", "duration of user sessions", []string{"environment", "region"})
  histogram.WithLabelValues("prodcution", "us-east-1").Observe(60)

  // Expose metrics via /metrics endpoint and an HTTP server
  http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
  log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Defaults

**Default buckets:**

```go
[]float64{0.01, 0.10, 0.50, 1.00, 5.00}
```

**Default quantiles:**

```go
map[float64]float64{
  0.1:  0.1,
  0.5:  0.05,
  0.95: 0.01,
  0.99: 0.001,
}
```
