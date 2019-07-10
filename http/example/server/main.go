package main

import (
	"math/rand"
	"net/http"
	"time"

	xhttp "github.com/moorara/goto/http"
	"github.com/moorara/goto/log"
	"github.com/moorara/goto/metrics"
	"github.com/moorara/goto/trace"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const port = ":10080"

type server struct {
	tracer opentracing.Tracer
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	// A random delay between 5ms to 50ms
	d := 5 + rand.Intn(45)
	time.Sleep(time.Duration(d) * time.Millisecond)

	logger, _ := xhttp.LoggerForRequest(r)
	logger.Info("message", "handled the request successfully!")

	// Create a new span
	parentSpan := opentracing.SpanFromContext(r.Context())
	span := s.tracer.StartSpan("send-greeting", opentracing.ChildOf(parentSpan.Context()))
	ext.DBType.Set(span, "sql")
	ext.DBStatement.Set(span, "SELECT * FROM messages")
	span.LogFields(opentracingLog.String("message", "sending the greeting message"))
	span.Finish()

	w.Write([]byte("Hello, World!"))
}

func main() {
	// Create a logger
	logger := log.NewLogger(log.Options{
		Name:        "server",
		Environment: "dev",
		Region:      "us-east-1",
		Component:   "hello-server",
	})

	// Create a metrics factory
	mf := metrics.NewFactory(metrics.FactoryOptions{})

	// Create a tracer
	tracer, closer, _ := trace.NewTracer(trace.Options{Name: "server"})
	defer closer.Close()

	// Create an http server middleware
	mid := xhttp.NewServerObservabilityMiddleware(logger, mf, tracer)

	// Wrap the http handler
	s := &server{tracer: tracer}
	h := mid.Wrap(s.handler)

	http.Handle("/", h)
	http.Handle("/metrics", promhttp.Handler())
	logger.Info("message", "starting http server ...", "port", port)
	panic(http.ListenAndServe(port, nil))
}
