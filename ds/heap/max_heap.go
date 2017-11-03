package heap

import (
	. "github.com/moorara/go-box/ds"
)

// MaxHeap represents a max-heap (priority queue) data structure
type MaxHeap struct {
	last            int
	keys            []Generic
	values          []Generic
	keyComparator   Comparator
	valueComparator Comparator
}

// NewMaxHeap creates a new max-heap (priority queue)
func NewMaxHeap(initialSize int, keyComparator Comparator, valueComparator Comparator) Heap {
	return &MaxHeap{
		last:            0,
		keys:            make([]Generic, initialSize),
		values:          make([]Generic, initialSize),
		keyComparator:   keyComparator,
		valueComparator: valueComparator,
	}
}

func (h *MaxHeap) resize(newSize int) {
	newKeys := make([]Generic, newSize)
	newValues := make([]Generic, newSize)

	copy(newKeys, h.keys)
	copy(newValues, h.values)

	h.keys = newKeys
	h.values = newValues
}

// Size returns size of heap
func (h *MaxHeap) Size() int {
	return h.last
}

// IsEmpty determines if heap is empty or not
func (h *MaxHeap) IsEmpty() bool {
	return h.last == 0
}

// Insert inserts a new key into heap
func (h *MaxHeap) Insert(key Generic, value Generic) {
	if h.last == len(h.keys)-1 {
		h.resize(len(h.keys) * 2)
	}

	h.last++
	var i int

	for i = h.last; true; i /= 2 {
		if i == 1 || h.keyComparator.Compare(key, h.keys[i/2]) <= 0 {
			break
		}
		h.keys[i] = h.keys[i/2]
		h.values[i] = h.values[i/2]
	}

	h.keys[i] = key
	h.values[i] = value
}

// Delete deletes the maximum key from heap
func (h *MaxHeap) Delete() (Generic, Generic) {
	if h.last == 0 {
		return nil, nil
	}

	maxKey := h.keys[1]
	maxValue := h.values[1]
	lastKey := h.keys[h.last]
	lastValue := h.values[h.last]

	h.last--
	var i, j int

	for i, j = 1, 2; j <= h.last; i, j = j, j*2 {
		if j < h.last && h.keyComparator.Compare(h.keys[j], h.keys[j+1]) < 0 {
			j++
		}
		if h.keyComparator.Compare(lastKey, h.keys[j]) >= 0 {
			break
		}
		h.keys[i] = h.keys[j]
		h.values[i] = h.values[j]
	}

	h.keys[i] = lastKey
	h.values[i] = lastValue

	if h.last < len(h.keys)/4 {
		h.resize(len(h.keys) / 2)
	}

	return maxKey, maxValue
}

// Peek returns the maximum key from heap without deleting it
func (h *MaxHeap) Peek() (Generic, Generic) {
	if h.last == 0 {
		return nil, nil
	}

	return h.keys[1], h.values[1]
}

// Contains checks if a key exists in heap
func (h *MaxHeap) ContainsKey(key Generic) bool {
	for i := 1; i <= h.last; i++ {
		if h.keyComparator.Compare(h.keys[i], key) == 0 {
			return true
		}
	}

	return false
}

// Contains checks if a value exists in heap
func (h *MaxHeap) ContainsValue(value Generic) bool {
	for i := 1; i <= h.last; i++ {
		if h.valueComparator.Compare(h.values[i], value) == 0 {
			return true
		}
	}

	return false
}
