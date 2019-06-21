# http

This package provides utilities for http servers.

## Quick Start

### Middleware

The `Middleware` enables observability through **logging**, **metrics**, and **tracing** for http handler functions.

```go
package main

import (
  "net/http"

  xhttp "github.com/moorara/goto/http"
  "github.com/moorara/goto/log"
  "github.com/moorara/goto/metrics"
  "github.com/moorara/goto/trace"
  opentracing "github.com/opentracing/opentracing-go"
  "github.com/opentracing/opentracing-go/ext"
  opentracingLog "github.com/opentracing/opentracing-go/log"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
  // Create a logger
  logger := log.NewLogger(log.Options{
    Format:      log.JSON,
    Name:        "handler",
    Environment: "dev",
    Region:      "us-east-1",
    Component:   "auth-service",
  })

  // Create a metrics factory
  mf := metrics.NewFactory(metrics.FactoryOptions{})

  // Create a tracer
  tracer, closer, _ := trace.NewTracer(trace.Options{Name: "auth-service"})
  defer closer.Close()

  // Create the http middleware and wrap a handler
  mid := xhttp.NewMiddleware(logger, mf, tracer)
  handler := mid.Wrap(func(w http.ResponseWriter, r *http.Request) {
    // Get the http logger passed down
    ctx := r.Context()
    val := ctx.Value(xhttp.LoggerContextKey)
    logger, _ := val.(*log.Logger)

    // Create a new span
    parentSpan := opentracing.SpanFromContext(ctx)
    span := tracer.StartSpan("send-greeting", opentracing.ChildOf(parentSpan.Context()))
    defer span.Finish()
    ext.DBType.Set(span, "sql")
    ext.DBStatement.Set(span, "SELECT * FROM messages")
    span.LogFields(opentracingLog.String("message", "sending the greeting message"))

    logger.Info("message", "handled the request successfully!")
    w.Write([]byte("Hello, World!"))
  })

  http.Handle("/", handler)
  http.Handle("/metrics", promhttp.Handler())
  http.ListenAndServe(":8080", nil)
}
```
