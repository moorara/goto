package heap

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxHeap(t *testing.T) {
	tests := []struct {
		name                  string
		initialSize           int
		compareKey            func(a, b interface{}) int
		compareValue          func(a, b interface{}) int
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
			"Empty",
			2,
			compareInt, compareString,
			[]int{}, []string{},
			0, true,
			0, "",
			[]int{}, []string{},
			[]int{}, []string{},
		},
		{
			"FewPairs",
			2,
			compareInt, compareString,
			[]int{10, 30, 20}, []string{"ten", "thirty", "twenty"},
			3, false,
			30, "thirty",
			[]int{30, 20, 10}, []string{"thirty", "twenty", "ten"},
			[]int{30, 20, 10}, []string{"thirty", "twenty", "ten"},
		},
		{
			"SomePairs",
			4,
			compareInt, compareString,
			[]int{10, 30, 20, 50, 40}, []string{"ten", "thirty", "twenty", "fifty", "forty"},
			5, false,
			50, "fifty",
			[]int{50, 40, 30, 20, 10}, []string{"fifty", "forty", "thirty", "twenty", "ten"},
			[]int{50, 40, 30, 20, 10}, []string{"fifty", "forty", "thirty", "twenty", "ten"},
		},
		{
			"MorePairs",
			4,
			compareInt, compareString,
			[]int{10, 30, 20, 50, 40, 60, 70, 90, 80}, []string{"ten", "thirty", "twenty", "fifty", "forty", "sixty", "seventy", "ninety", "eighty"},
			9, false,
			90, "ninety",
			[]int{90, 80, 70, 60, 50, 40, 30, 20, 10}, []string{"ninety", "eighty", "seventy", "sixty", "fifty", "forty", "thirty", "twenty", "ten"},
			[]int{90, 80, 70, 60, 50, 40, 30, 20, 10}, []string{"ninety", "eighty", "seventy", "sixty", "fifty", "forty", "thirty", "twenty", "ten"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMaxHeap(tc.initialSize, tc.compareKey, tc.compareValue)

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

			for i := 0; i < len(tc.insertKeys); i++ {
				heap.Insert(tc.insertKeys[i], tc.insertValues[i])
			}

			assert.Equal(t, tc.expectedSize, heap.Size())
			assert.Equal(t, tc.expectedIsEmpty, heap.IsEmpty())

			peekKey, peekValue = heap.Peek()
			if tc.expectedSize == 0 {
				assert.Nil(t, peekKey)
				assert.Nil(t, peekValue)
			} else {
				assert.Equal(t, tc.expectedPeekKey, peekKey)
				assert.Equal(t, tc.expectedPeekValue, peekValue)
			}

			for _, key := range tc.expectedContainsKey {
				assert.True(t, heap.ContainsKey(key))
			}

			for _, value := range tc.expectedContainsValue {
				assert.True(t, heap.ContainsValue(value))
			}

			for i := 0; i < len(tc.expectedDeleteKeys); i++ {
				deleteKey, deleteValue = heap.Delete()
				assert.Equal(t, tc.expectedDeleteKeys[i], deleteKey)
				assert.Equal(t, tc.expectedDeleteValues[i], deleteValue)
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
		})
	}
}

func BenchmarkMaxHeap(b *testing.B) {
	heapSize := 1024
	minInt := 0
	maxInt := 1000000

	rand.Seed(time.Now().UTC().UnixNano())

	b.Run("Insert", func(b *testing.B) {
		heap := NewMaxHeap(heapSize, compareInt, compareString)
		items := genIntSlice(b.N, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			heap.Insert(items[n], "")
		}
	})

	b.Run("Delete", func(b *testing.B) {
		heap := NewMaxHeap(heapSize, compareInt, compareString)
		items := genIntSlice(b.N, minInt, maxInt)
		for n := 0; n < b.N; n++ {
			heap.Insert(items[n], "")
		}
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			heap.Delete()
		}
	})
}
