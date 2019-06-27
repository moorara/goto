package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	xhttp "github.com/moorara/goto/http"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Create a logger
	logger := log.NewLogger(log.Options{
		Name:        "client",
		Environment: "dev",
		Region:      "us-east-1",
		Component:   "hello-client",
	})

	// Create a metrics factory
	mf := metrics.NewFactory(metrics.FactoryOptions{})

	// Create a tracer
	tracer, closer, _ := trace.NewTracer(trace.Options{Name: "client"})
	defer closer.Close()

	// Create an http client middleware
	mid := xhttp.NewClientObservabilityMiddleware(logger, mf, tracer)

	// Create an http client
	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: &http.Transport{},
	}

	// Create an http request
	req, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	// Make the request to http server
	res, err := mid.Wrap(context.Background(), req, client.Do)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	logger.Info("message", res.Status, "res.body", string(b), "res.statusCode", res.StatusCode)

	http.Handle("/metrics", promhttp.Handler())
	logger.Info("message", "starting server on localhost:8081 ...")
	panic(http.ListenAndServe(":8081", nil))
}
