package ds

// Stack represents a stack data structure
type Stack interface {
	Size() int
	IsEmpty() bool
	Push(Value)
	Pop() Value
	Peek() Value
	Contains(Value) bool
}

type arrayStack struct {
	listSize  int
	nodeSize  int
	nodeIndex int
	topNode   *arrayNode
}

// NewStack creates a new array-list stack
func NewStack(nodeSize int) Stack {
	return &arrayStack{
		listSize:  0,
		nodeSize:  nodeSize,
		nodeIndex: -1,
		topNode:   nil,
	}
}

func (s *arrayStack) Size() int {
	return s.listSize
}

func (s *arrayStack) IsEmpty() bool {
	return s.listSize == 0
}

func (s *arrayStack) Push(item Value) {
	s.listSize++
	s.nodeIndex++

	if s.topNode == nil {
		s.topNode = newArrayNode(s.nodeSize, nil)
	} else {
		if s.nodeIndex == s.nodeSize {
			s.nodeIndex = 0
			s.topNode = newArrayNode(s.nodeSize, s.topNode)
		}
	}

	s.topNode.block[s.nodeIndex] = item
}

func (s *arrayStack) Pop() Value {
	if s.listSize == 0 {
		return nil
	}

	item := s.topNode.block[s.nodeIndex]
	s.nodeIndex--
	s.listSize--

	if s.nodeIndex == -1 {
		s.topNode = s.topNode.next
		if s.topNode != nil {
			s.nodeIndex = s.nodeSize - 1
		}
	}

	return item
}

func (s *arrayStack) Peek() Value {
	if s.listSize == 0 {
		return nil
	}

	return s.topNode.block[s.nodeIndex]
}

func (s *arrayStack) Contains(item Value) bool {
	n := s.topNode
	i := s.nodeIndex

	for n != nil {
		if n.block[i].Compare(item) == 0 {
			return true
		}

		i--
		if i < 0 {
			n = n.next
			i = s.nodeSize - 1
		}
	}

	return false
}
