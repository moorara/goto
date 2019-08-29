// Package heap provides min heap and max heap data structures.
//
// Deprecated: this package has been frozen and deprecated in favor of github.com/moorara/algo/ds/heap
package heap

// Heap represents a heap (priority queue) data structure
type Heap interface {
	Size() int
	IsEmpty() bool
	Insert(interface{}, interface{})
	Delete() (interface{}, interface{})
	Peek() (interface{}, interface{})
	ContainsKey(interface{}) bool
	ContainsValue(interface{}) bool
}
