package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArrayNodeWithSizeAndNext(t *testing.T) {
	tests := []struct {
		size int
		next *arrayNode
	}{
		{64, nil},
		{256, nil},
		{1024, &arrayNode{}},
		{4096, &arrayNode{}},
	}

	for _, tc := range tests {
		n := newArrayNode(tc.size, tc.next)

		assert.Equal(t, tc.next, n.next)
		assert.Equal(t, tc.size, len(n.block))
	}
}
