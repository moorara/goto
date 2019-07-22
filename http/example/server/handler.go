package main

import (
	"math/rand"
	"net/http"
	"time"

	xhttp "github.com/moorara/goto/http"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracingLog "github.com/opentracing/opentracing-go/log"
)

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
