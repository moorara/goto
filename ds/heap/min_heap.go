package heap

import (
	. "github.com/moorara/go-box/ds"
)

// MinHeap represents a min-heap (priority queue) data structure
type MinHeap struct {
	last            int
	keys            []Generic
	values          []Generic
	keyComparator   Comparator
	valueComparator Comparator
}

// NewMinHeap creates a new min-heap (priority queue)
func NewMinHeap(initialSize int, keyComparator Comparator, valueComparator Comparator) Heap {
	return &MinHeap{
		last:            0,
		keys:            make([]Generic, initialSize),
		values:          make([]Generic, initialSize),
		keyComparator:   keyComparator,
		valueComparator: valueComparator,
	}
}

func (h *MinHeap) resize(newSize int) {
	newKeys := make([]Generic, newSize)
	newValues := make([]Generic, newSize)

	copy(newKeys, h.keys)
	copy(newValues, h.values)

	h.keys = newKeys
	h.values = newValues
}

// Size returns size of heap
func (h *MinHeap) Size() int {
	return h.last
}

// IsEmpty determines if heap is empty or not
func (h *MinHeap) IsEmpty() bool {
	return h.last == 0
}

// Insert inserts a new key into heap
func (h *MinHeap) Insert(key Generic, value Generic) {
	if h.last == len(h.keys)-1 {
		h.resize(len(h.keys) * 2)
	}

	h.last++
	var i int

	for i = h.last; true; i /= 2 {
		if i == 1 || h.keyComparator.Compare(key, h.keys[i/2]) >= 0 {
			break
		}
		h.keys[i] = h.keys[i/2]
		h.values[i] = h.values[i/2]
	}

	h.keys[i] = key
	h.values[i] = value
}

// Delete deletes the minimum key from heap
func (h *MinHeap) Delete() (Generic, Generic) {
	if h.last == 0 {
		return nil, nil
	}

	minKey := h.keys[1]
	minValue := h.values[1]
	lastKey := h.keys[h.last]
	lastValue := h.values[h.last]

	h.last--
	var i, j int

	for i, j = 1, 2; j <= h.last; i, j = j, j*2 {
		if j < h.last && h.keyComparator.Compare(h.keys[j], h.keys[j+1]) > 0 {
			j++
		}
		if h.keyComparator.Compare(lastKey, h.keys[j]) <= 0 {
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

	return minKey, minValue
}

// Peek returns the minimum key from heap without deleting it
func (h *MinHeap) Peek() (Generic, Generic) {
	if h.last == 0 {
		return nil, nil
	}

	return h.keys[1], h.values[1]
}

// Contains checks if a key exists in heap
func (h *MinHeap) ContainsKey(key Generic) bool {
	for i := 1; i <= h.last; i++ {
		if h.keyComparator.Compare(h.keys[i], key) == 0 {
			return true
		}
	}

	return false
}

// Contains checks if a value exists in heap
func (h *MinHeap) ContainsValue(value Generic) bool {
	for i := 1; i <= h.last; i++ {
		if h.valueComparator.Compare(h.values[i], value) == 0 {
			return true
		}
	}

	return false
}
