package grpc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestMetadataTextMap(t *testing.T) {
	tests := []struct {
		name          string
		md            metadata.MD
		pairs         map[string][]string
		handlerError  error
		expectedPairs map[string][]string
		expectedError error
	}{
		{
			name:          "Empty",
			md:            metadata.Pairs(),
			pairs:         nil,
			handlerError:  nil,
			expectedPairs: nil,
			expectedError: nil,
		},
		{
			name:         "Error",
			md:           metadata.Pairs("user-id", "1111"),
			pairs:        nil,
			handlerError: errors.New("handler error"),
			expectedPairs: map[string][]string{
				"user-id": []string{"1111"},
			},
			expectedError: errors.New("handler error"),
		},
		{
			name: "Success",
			md:   metadata.Pairs("user-id", "1111"),
			pairs: map[string][]string{
				"tenant-id": []string{"aaaa", "bbbb"},
			},
			handlerError: nil,
			expectedPairs: map[string][]string{
				"user-id":   []string{"1111"},
				"tenant-id": []string{"aaaa", "bbbb"},
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := metadataTextMap{tc.md}

			t.Run("Set", func(t *testing.T) {
				for k, vals := range tc.pairs {
					for _, v := range vals {
						m.Set(k, v)
					}
				}
			})

			t.Run("ForeachKey", func(t *testing.T) {
				err := m.ForeachKey(func(key, val string) error {
					assert.Contains(t, tc.expectedPairs[key], val)
					return tc.handlerError
				})

				assert.Equal(t, tc.expectedError, err)
			})
		})
	}
}
