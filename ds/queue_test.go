package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * stringComparator is defined in ds_test.go
 * stringBitStringer is defined in queue_test.go
 */

func TestQueue(t *testing.T) {
	tests := []struct {
		nodeSize             int
		comparator           Comparator
		enqueueItems         []string
		expectedSize         int
		expectedIsEmpty      bool
		expectedPeek         string
		expectedContains     []string
		expectedDequeueItems []string
	}{
		{
			2,
			&StringComparator{},
			[]string{},
			0, true,
			"",
			[]string{},
			[]string{},
		},
		{
			2,
			&StringComparator{},
			[]string{"a", "b"},
			2, false,
			"a",
			[]string{"a", "b"},
			[]string{"a", "b"},
		},
		{
			2,
			&StringComparator{},
			[]string{"a", "b", "c"},
			3, false,
			"a",
			[]string{"a", "b", "c"},
			[]string{"a", "b", "c"},
		},
		{
			2,
			&StringComparator{},
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			7, false,
			"a",
			[]string{"a", "b", "c", "d", "e", "f", "g"},
			[]string{"a", "b", "c", "d", "e", "f", "g"},
		},
	}

	for _, test := range tests {
		queue := NewQueue(test.nodeSize, test.comparator)

		// Queue initially should be empty
		assert.Zero(t, queue.Size())
		assert.True(t, queue.IsEmpty())
		assert.Nil(t, queue.Peek())
		queue.Contains(nil)
		assert.Nil(t, queue.Dequeue())

		for _, item := range test.enqueueItems {
			queue.Enqueue(item)
		}

		assert.Equal(t, test.expectedSize, queue.Size())
		assert.Equal(t, test.expectedIsEmpty, queue.IsEmpty())

		if test.expectedSize == 0 {
			assert.Nil(t, queue.Peek())
		} else {
			assert.Equal(t, test.expectedPeek, queue.Peek())
		}

		for _, item := range test.expectedContains {
			assert.True(t, queue.Contains(item))
		}

		for _, item := range test.expectedDequeueItems {
			assert.Equal(t, item, queue.Dequeue())
		}

		// Queue should be empty at the end
		assert.Zero(t, queue.Size())
		assert.True(t, queue.IsEmpty())
		assert.Nil(t, queue.Peek())
		queue.Contains(nil)
		assert.Nil(t, queue.Dequeue())
	}
}

func BenchmarkQueue(b *testing.B) {
	queue := NewQueue(1024, &IntComparator{})

	for n := 0; n < b.N; n++ {
		queue.Enqueue(n)
	}

	for n := 0; n < b.N; n++ {
		queue.Dequeue()
	}
}
