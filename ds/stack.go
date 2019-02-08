package ds

// Stack represents a stack data structure
type Stack interface {
	Size() int
	IsEmpty() bool
	Push(interface{})
	Pop() interface{}
	Peek() interface{}
	Contains(interface{}) bool
}

type arrayStack struct {
	listSize  int
	nodeSize  int
	nodeIndex int
	topNode   *arrayNode
	compare   func(a, b interface{}) int
}

// NewStack creates a new array-list stack
func NewStack(nodeSize int, compare func(a, b interface{}) int) Stack {
	return &arrayStack{
		listSize:  0,
		nodeSize:  nodeSize,
		nodeIndex: -1,
		topNode:   nil,
		compare:   compare,
	}
}

func (s *arrayStack) Size() int {
	return s.listSize
}

func (s *arrayStack) IsEmpty() bool {
	return s.listSize == 0
}

func (s *arrayStack) Push(item interface{}) {
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

func (s *arrayStack) Pop() interface{} {
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

func (s *arrayStack) Peek() interface{} {
	if s.listSize == 0 {
		return nil
	}

	return s.topNode.block[s.nodeIndex]
}

func (s *arrayStack) Contains(item interface{}) bool {
	n := s.topNode
	i := s.nodeIndex

	for n != nil {
		if s.compare(n.block[i], item) == 0 {
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
