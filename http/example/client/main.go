package main

import (
	"context"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

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
		Component:   "hello-client",
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
	mid := xhttp.NewClientObservabilityMiddleware(logger, mf, tracer)

	for {
		// A random delay between 1s to 5s
		d := 1 + rand.Intn(4)
		time.Sleep(time.Duration(d) * time.Second)

		// Create an http client
		client := http.Client{
			Timeout:   10 * time.Second,
			Transport: &http.Transport{},
		}

		// Create an http request
		req, _ := http.NewRequest("GET", serverAddress+"/", nil)

		// Make the request to http server
		res, err := mid.Wrap(context.Background(), req, client.Do)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		logger.Info("message", string(b))
	}
}
