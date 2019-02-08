package st

const (
	// TraversePreOrder represents pre-order traversal order
	TraversePreOrder = 0
	// TraverseInOrder represents in-order traversal order
	TraverseInOrder = 1
	// TraversePostOrder represents post-order traversal order
	TraversePostOrder = 2
)

type (
	// VisitFunc represents the function for visting a key-value
	VisitFunc func(interface{}, interface{}) bool

	// KeyValue represents a key-value pair
	KeyValue struct {
		key   interface{}
		value interface{}
	}

	// SymbolTable represents an unordered symbol table (key-value collection)
	SymbolTable interface {
		verify() bool
		Size() int
		Height() int
		IsEmpty() bool
		Put(interface{}, interface{})
		Get(interface{}) (interface{}, bool)
		Delete(interface{}) (interface{}, bool)
		KeyValues() []KeyValue
	}

	// OrderedSymbolTable represents an ordered symbol table (key-value collection)
	OrderedSymbolTable interface {
		SymbolTable
		Min() (interface{}, interface{})
		Max() (interface{}, interface{})
		Floor(interface{}) (interface{}, interface{})
		Ceiling(interface{}) (interface{}, interface{})
		Rank(interface{}) int
		Select(int) (interface{}, interface{})
		DeleteMin() (interface{}, interface{})
		DeleteMax() (interface{}, interface{})
		RangeSize(interface{}, interface{}) int
		Range(interface{}, interface{}) []KeyValue
		Traverse(int, VisitFunc)
		Graphviz() string
	}
)
