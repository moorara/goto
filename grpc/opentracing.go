package grpc

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

// MetadataTextMap implements opentracing.TextMapReader and opentracing.TextMapWriter
type MetadataTextMap struct {
	metadata.MD
}

// Set normalizes the key and appends a value to it
// https://godoc.org/github.com/opentracing/opentracing-go#TextMapWriter
func (m *MetadataTextMap) Set(key, val string) {
	key = strings.ToLower(key)
	m.MD[key] = append(m.MD[key], val)
}

// ForeachKey is an iterator for all key-values pairs
// https://godoc.org/github.com/opentracing/opentracing-go#TextMapReader
func (m *MetadataTextMap) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range m.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
