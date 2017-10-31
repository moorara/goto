package ds

type arrayNode struct {
	block []Generic
	next  *arrayNode
}

func newArrayNode(size int, next *arrayNode) *arrayNode {
	return &arrayNode{
		block: make([]Generic, size),
		next:  next,
	}
}
