package ds

type arrayNode struct {
	block []interface{}
	next  *arrayNode
}

func newArrayNode(size int, next *arrayNode) *arrayNode {
	return &arrayNode{
		block: make([]interface{}, size),
		next:  next,
	}
}
