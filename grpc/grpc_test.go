package grpc

import (
	"context"
	"testing"

	"github.com/moorara/goto/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestLoggerFromContext(t *testing.T) {
	tests := []struct {
		name       string
		ctx        context.Context
		logger     *log.Logger
		expectedOK bool
	}{
		{
			name:       "WithoutLogger",
			ctx:        context.Background(),
			logger:     nil,
			expectedOK: false,
		},
		{
			name:       "WithLogger",
			ctx:        context.Background(),
			logger:     log.NewNopLogger(),
			expectedOK: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.logger != nil {
				tc.ctx = context.WithValue(tc.ctx, loggerContextKey, tc.logger)
			}

			logger, ok := LoggerFromContext(tc.ctx)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.logger, logger)
		})
	}
}

func TestParseMethod(t *testing.T) {
	tests := []struct {
		name            string
		fullMethod      string
		expectedOK      bool
		expectedPackage string
		expectedService string
		expectedMethod  string
	}{
		{
			name:            "Failure",
			fullMethod:      "GetPlacesInZone",
			expectedOK:      false,
			expectedPackage: "",
			expectedService: "",
			expectedMethod:  "",
		},
		{
			name:            "Success",
			fullMethod:      "/zonePB.ZoneManager/GetPlacesInZone",
			expectedOK:      true,
			expectedPackage: "zonePB",
			expectedService: "ZoneManager",
			expectedMethod:  "GetPlacesInZone",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pkg, service, method, ok := parseMethod(tc.fullMethod)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedPackage, pkg)
			assert.Equal(t, tc.expectedService, service)
			assert.Equal(t, tc.expectedMethod, method)
		})
	}
}

func TestXServerStream(t *testing.T) {
	ctx1 := context.WithValue(context.Background(), contextKey("user-id"), "1111")
	ctx2 := context.WithValue(ctx1, contextKey("trace-id"), "aaaa")

	tests := []struct {
		name        string
		stream      grpc.ServerStream
		ctx         context.Context
		expextedCtx context.Context
	}{
		{
			name:        "WithoutContext",
			stream:      &mockServerStream{ContextOutContext: ctx1},
			ctx:         nil,
			expextedCtx: ctx1,
		},
		{
			name:        "WithContext",
			stream:      &mockServerStream{ContextOutContext: ctx1},
			ctx:         ctx2,
			expextedCtx: ctx2,
		},
		{
			name:        "AlreadyWrapped",
			stream:      &xServerStream{ServerStream: &mockServerStream{ContextOutContext: ctx1}},
			ctx:         nil,
			expextedCtx: ctx1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ss := ServerStreamWithContext(tc.stream, tc.ctx)

			assert.NotNil(t, ss)
			assert.Equal(t, tc.expextedCtx, ss.Context())
		})
	}
}
