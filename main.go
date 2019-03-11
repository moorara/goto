package main

import (
	"github.com/moorara/goto/trace"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	sampler := trace.NewConstSampler(true)
	reporter := trace.NewAgentReporter("localhost:6831", false)
	tracer, closer, _ := trace.NewTracer("hello_service", sampler, reporter, nil, prometheus.DefaultRegisterer)
	defer closer.Close()

	span := tracer.StartSpan("hello-world")
	defer span.Finish()
	span.LogFields(
		log.String("environment", "prodcution"),
		log.String("region", "us-east-1"),
	)
	ext.HTTPMethod.Set(span, "GET")
	ext.HTTPStatusCode.Set(span, 200)
}