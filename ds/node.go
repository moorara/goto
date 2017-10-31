package ds

type arrayNode struct {
	block []Value
	next  *arrayNode
}

func newArrayNode(size int, next *arrayNode) *arrayNode {
	return &arrayNode{
		block: make([]Value, size),
		next:  next,
	}
}
