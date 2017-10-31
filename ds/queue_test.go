package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * value is defined in ds_test.go
 * newValueArray is defined in queue_test.go
 */

func TestQueue(t *testing.T) {
	tests := []struct {
		nodeSize             int
		expectedSize         int
		expectedIsEmpty      bool
		expectedPeek         value
		enqueueItems         []value
		expectedContains     []value
		expectedDequeueItems []value
	}{
		{
			2, 0, true,
			value{},
			[]value{},
			[]value{},
			[]value{},
		},
		{
			2, 2, false,
			value{"a"},
			newValueArray("a", "b"),
			newValueArray("a", "b"),
			newValueArray("a", "b"),
		},
		{
			2, 3, false,
			value{"a"},
			newValueArray("a", "b", "c"),
			newValueArray("a", "b", "c"),
			newValueArray("a", "b", "c"),
		},
		{
			2, 7, false,
			value{"a"},
			newValueArray("a", "b", "c", "d", "e", "f", "g"),
			newValueArray("a", "b", "c", "d", "e", "f", "g"),
			newValueArray("a", "b", "c", "d", "e", "f", "g"),
		},
	}

	for _, test := range tests {
		queue := NewQueue(test.nodeSize)

		// Queue initially should be empty
		assert.Zero(t, queue.Size())
		assert.True(t, queue.IsEmpty())
		assert.Nil(t, queue.Dequeue())
		assert.Nil(t, queue.Peek())
		queue.Contains(nil)

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
		assert.Nil(t, queue.Dequeue())
		assert.Nil(t, queue.Peek())
	}
}
