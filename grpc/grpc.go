package grpc

import (
	"context"
	"regexp"

	"github.com/moorara/goto/log"
	"google.golang.org/grpc"
)

// contextKey is the type for the keys added to context
type contextKey string

var (
	loggerContextKey = contextKey("logger")
	methodRegex      = regexp.MustCompile(`(/|\.)`)
)

// LoggerFromContext returns a logger set by grpc server interceptor on each incoming context
func LoggerFromContext(ctx context.Context) (*log.Logger, bool) {
	val := ctx.Value(loggerContextKey)
	logger, ok := val.(*log.Logger)

	return logger, ok
}

func parseMethod(fullMethod string) (string, string, string, bool) {
	// fullMethod should have the form /package.service/method
	subs := methodRegex.Split(fullMethod, 4)
	if len(subs) != 4 {
		return "", "", "", false
	}

	return subs[1], subs[2], subs[3], true
}

type xServerStream struct {
	grpc.ServerStream
	context context.Context
}

func (s *xServerStream) Context() context.Context {
	if s.context == nil {
		return s.ServerStream.Context()
	}

	return s.context
}

// ServerStreamWithContext return new grpc.ServerStream with a new context
func ServerStreamWithContext(stream grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	if ss, ok := stream.(*xServerStream); ok {
		return ss
	}

	return &xServerStream{
		ServerStream: stream,
		context:      ctx,
	}
}
