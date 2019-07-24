package main

import (
	"net/http"

	xhttp "github.com/moorara/goto/http"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const port = ":10081"
const serverAddress = "http://localhost:10080"

func main() {
	// Create a logger
	logger := log.NewLogger(log.Options{
		Name:        "client",
		Environment: "dev",
		Region:      "us-east-1",
		Component:   "http-client",
	})

	// Create a metrics factory
	mf := metrics.NewFactory(metrics.FactoryOptions{})

	// Create a tracer
	tracer, closer, _ := trace.NewTracer(trace.Options{Name: "client"})
	defer closer.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		logger.Info("message", "starting http server ...", "port", port)
		panic(http.ListenAndServe(port, nil))
	}()

	// Create an http client middleware
	mid := xhttp.NewClientMiddleware(logger, mf, tracer)

	c := &client{
		logger: logger,
		mid:    mid,
	}

	for {
		c.call()
	}
}
