package heap

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMinHeap(t *testing.T) {
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
			[]int{30, 10, 20}, []string{"thirty", "ten", "twenty"},
			3, false,
			10, "ten",
			[]int{10, 20, 30}, []string{"ten", "twenty", "thirty"},
			[]int{10, 20, 30}, []string{"ten", "twenty", "thirty"},
		},
		{
			"SomePairs",
			4,
			compareInt, compareString,
			[]int{50, 30, 40, 10, 20}, []string{"fifty", "thirty", "forty", "ten", "twenty"},
			5, false,
			10, "ten",
			[]int{10, 20, 30, 40, 50}, []string{"ten", "twenty", "thirty", "forty", "fifty"},
			[]int{10, 20, 30, 40, 50}, []string{"ten", "twenty", "thirty", "forty", "fifty"},
		},
		{
			"MorePairs",
			4,
			compareInt, compareString,
			[]int{90, 80, 70, 40, 50, 60, 30, 10, 20}, []string{"ninety", "eighty", "seventy", "forty", "fifty", "sixty", "thirty", "ten", "twenty"},
			9, false,
			10, "ten",
			[]int{10, 20, 30, 40, 50, 60, 70, 80, 90}, []string{"ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"},
			[]int{10, 20, 30, 40, 50, 60, 70, 80, 90}, []string{"ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			heap := NewMinHeap(tc.initialSize, tc.compareKey, tc.compareValue)

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

func BenchmarkMinHeap(b *testing.B) {
	heapSize := 1024
	minInt := 0
	maxInt := 1000000

	rand.Seed(time.Now().UTC().UnixNano())

	b.Run("Insert", func(b *testing.B) {
		heap := NewMinHeap(heapSize, compareInt, compareString)
		items := genIntSlice(b.N, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			heap.Insert(items[n], "")
		}
	})

	b.Run("Delete", func(b *testing.B) {
		heap := NewMinHeap(heapSize, compareInt, compareString)
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
