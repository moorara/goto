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
