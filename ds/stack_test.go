package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * value is defined in ds_test.go
 * newValueArray is defined in queue_test.go
 */

func TestStack(t *testing.T) {
	tests := []struct {
		nodeSize         int
		expectedSize     int
		expectedIsEmpty  bool
		expectedPeek     value
		pushItems        []value
		expectedContains []value
		expectedPopItems []value
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
			value{"b"},
			newValueArray("a", "b"),
			newValueArray("a", "b"),
			newValueArray("b", "a"),
		},
		{
			2, 3, false,
			value{"c"},
			newValueArray("a", "b", "c"),
			newValueArray("a", "b", "c"),
			newValueArray("c", "b", "a"),
		},
		{
			2, 7, false,
			value{"g"},
			newValueArray("a", "b", "c", "d", "e", "f", "g"),
			newValueArray("a", "b", "c", "d", "e", "f", "g"),
			newValueArray("g", "f", "e", "d", "c", "b", "a"),
		},
	}

	for _, test := range tests {
		stack := NewStack(test.nodeSize)

		// Stack initially should be empty
		assert.Zero(t, stack.Size())
		assert.True(t, stack.IsEmpty())
		assert.Nil(t, stack.Pop())
		assert.Nil(t, stack.Peek())
		assert.False(t, stack.Contains(nil))

		for _, item := range test.pushItems {
			stack.Push(item)
		}

		assert.Equal(t, test.expectedSize, stack.Size())
		assert.Equal(t, test.expectedIsEmpty, stack.IsEmpty())

		if test.expectedSize == 0 {
			assert.Nil(t, stack.Peek())
		} else {
			assert.Equal(t, test.expectedPeek, stack.Peek())
		}

		for _, item := range test.expectedContains {
			assert.True(t, stack.Contains(item))
		}

		for _, item := range test.expectedPopItems {
			assert.Equal(t, item, stack.Pop())
		}

		// Stack should be empty at the end
		assert.Zero(t, stack.Size())
		assert.True(t, stack.IsEmpty())
		assert.Nil(t, stack.Pop())
		assert.Nil(t, stack.Peek())
		assert.False(t, stack.Contains(nil))
	}
}
