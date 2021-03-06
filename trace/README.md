# trace

This is a helper package for creating a [Jaeger](https://www.jaegertracing.io) tracer
that reports traces in [OpenTracing](https://opentracing.io) format.

## Quick Start

For creating a tracer with a *constant sampler* and an *agent reporter*:

```go
package main

import (
  "github.com/moorara/goto/trace"
  "github.com/opentracing/opentracing-go/ext"
  "github.com/opentracing/opentracing-go/log"
  "github.com/prometheus/client_golang/prometheus"
)

func main() {
  tracer, closer, _ := trace.NewTracer(trace.Options{Name: "hello_service"})
  defer closer.Close()

  span := tracer.StartSpan("hello-world")
  defer span.Finish()
  ext.HTTPMethod.Set(span, "GET")
  ext.HTTPStatusCode.Set(span, 200)
  span.LogFields(
    log.String("environment", "prodcution"),
    log.String("region", "us-east-1"),
  )
}
```
