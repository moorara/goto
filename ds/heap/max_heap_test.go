package heap

import (
	"strconv"
	"testing"

	. "github.com/moorara/go-box/ds"
	"github.com/stretchr/testify/assert"
)

func TestMaxHeap(t *testing.T) {
	tests := []struct {
		initialSize           int
		keyComparator         Comparator
		valueComparator       Comparator
		insertKeys            []int
		insertValues          []string
		expectedSize          int
		expectedIsEmpty       bool
		expectedPeekKey       int
		expectedPeekValue     string
		expectedContainsKey   []int
		expectedContainsValue []string
		expectedDeleteKeys    []int
		expectedDeleteValues  []string
	}{
		{
			2,
			&IntComparator{}, &StringComparator{},
			[]int{}, []string{},
			0, true,
			0, "",
			[]int{}, []string{},
			[]int{}, []string{},
		},
		{
			2,
			&IntComparator{}, &StringComparator{},
			[]int{10, 30, 20}, []string{"ten", "thirty", "twenty"},
			3, false,
			30, "thirty",
			[]int{30, 20, 10}, []string{"thirty", "twenty", "ten"},
			[]int{30, 20, 10}, []string{"thirty", "twenty", "ten"},
		},
		{
			4,
			&IntComparator{}, &StringComparator{},
			[]int{10, 30, 20, 50, 40}, []string{"ten", "thirty", "twenty", "fifty", "forty"},
			5, false,
			50, "fifty",
			[]int{50, 40, 30, 20, 10}, []string{"fifty", "forty", "thirty", "twenty", "ten"},
			[]int{50, 40, 30, 20, 10}, []string{"fifty", "forty", "thirty", "twenty", "ten"},
		},
		{
			4,
			&IntComparator{}, &StringComparator{},
			[]int{10, 30, 20, 50, 40, 60, 70, 90, 80}, []string{"ten", "thirty", "twenty", "fifty", "forty", "sixty", "seventy", "ninety", "eighty"},
			9, false,
			90, "ninety",
			[]int{90, 80, 70, 60, 50, 40, 30, 20, 10}, []string{"ninety", "eighty", "seventy", "sixty", "fifty", "forty", "thirty", "twenty", "ten"},
			[]int{90, 80, 70, 60, 50, 40, 30, 20, 10}, []string{"ninety", "eighty", "seventy", "sixty", "fifty", "forty", "thirty", "twenty", "ten"},
		},
	}

	for _, test := range tests {
		heap := NewMaxHeap(test.initialSize, test.keyComparator, test.valueComparator)

		// Heap initially should be empty
		peekKey, peekValue := heap.Peek()
		deleteKey, deleteValue := heap.Delete()
		assert.Nil(t, peekKey)
		assert.Nil(t, peekValue)
		assert.Nil(t, deleteKey)
		assert.Nil(t, deleteValue)
		assert.Zero(t, heap.Size())
		assert.True(t, heap.IsEmpty())
		assert.False(t, heap.ContainsKey(nil))
		assert.False(t, heap.ContainsValue(nil))

		for i := 0; i < len(test.insertKeys); i++ {
			heap.Insert(test.insertKeys[i], test.insertValues[i])
		}

		assert.Equal(t, test.expectedSize, heap.Size())
		assert.Equal(t, test.expectedIsEmpty, heap.IsEmpty())

		peekKey, peekValue = heap.Peek()
		if test.expectedSize == 0 {
			assert.Nil(t, peekKey)
			assert.Nil(t, peekValue)
		} else {
			assert.Equal(t, test.expectedPeekKey, peekKey)
			assert.Equal(t, test.expectedPeekValue, peekValue)
		}

		for _, key := range test.expectedContainsKey {
			assert.True(t, heap.ContainsKey(key))
		}

		for _, value := range test.expectedContainsValue {
			assert.True(t, heap.ContainsValue(value))
		}

		for i := 0; i < len(test.expectedDeleteKeys); i++ {
			deleteKey, deleteValue = heap.Delete()
			assert.Equal(t, test.expectedDeleteKeys[i], deleteKey)
			assert.Equal(t, test.expectedDeleteValues[i], deleteValue)
		}

		// Heap should be empty at the end
		peekKey, peekValue = heap.Peek()
		deleteKey, deleteValue = heap.Delete()
		assert.Nil(t, peekKey)
		assert.Nil(t, peekValue)
		assert.Nil(t, deleteKey)
		assert.Nil(t, deleteValue)
		assert.Zero(t, heap.Size())
		assert.True(t, heap.IsEmpty())
		assert.False(t, heap.ContainsKey(nil))
		assert.False(t, heap.ContainsValue(nil))
	}
}

func BenchmarkMaxHeap(b *testing.B) {
	heap := NewMaxHeap(1024, &IntComparator{}, &StringComparator{})

	for n := 0; n < b.N; n++ {
		heap.Insert(n, strconv.Itoa(n))
	}

	for n := 0; n < b.N; n++ {
		heap.Delete()
	}
}