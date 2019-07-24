package main

import (
	"net/http"

	xhttp "github.com/moorara/goto/http"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const port = ":10080"

func main() {
	// Create a logger
	logger := log.NewLogger(log.Options{
		Name:        "server",
		Environment: "dev",
		Region:      "us-east-1",
		Component:   "http-server",
	})

	// Create a metrics factory
	mf := metrics.NewFactory(metrics.FactoryOptions{})

	// Create a tracer
	tracer, closer, _ := trace.NewTracer(trace.Options{Name: "server"})
	defer closer.Close()

	// Create an http server middleware
	mid := xhttp.NewServerMiddleware(logger, mf, tracer)

	s := &server{tracer: tracer}
	h := mid.Metrics(mid.RequestID(mid.Tracing(mid.Logging(s.handler))))

	http.Handle("/", h)
	http.Handle("/metrics", promhttp.Handler())
	logger.Info("message", "starting http server ...", "port", port)
	panic(http.ListenAndServe(port, nil))
}
