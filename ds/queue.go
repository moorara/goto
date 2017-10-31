package ds

// Queue represents a stack data structure
type Queue interface {
	Size() int
	IsEmpty() bool
	Enqueue(Value)
	Dequeue() Value
	Peek() Value
	Contains(Value) bool
}

type arrayQueue struct {
	listSize       int
	nodeSize       int
	frontNodeIndex int
	rearNodeIndex  int
	frontNode      *arrayNode
	rearNode       *arrayNode
}

// NewQueue creates a new array-list queue
func NewQueue(nodeSize int) Queue {
	return &arrayQueue{
		listSize:       0,
		nodeSize:       nodeSize,
		frontNodeIndex: -1,
		rearNodeIndex:  -1,
		frontNode:      nil,
		rearNode:       nil,
	}
}

func (q *arrayQueue) Size() int {
	return q.listSize
}

func (q *arrayQueue) IsEmpty() bool {
	return q.listSize == 0
}

func (q *arrayQueue) Enqueue(item Value) {
	if q.frontNode == nil && q.rearNode == nil {
		q.frontNodeIndex = 0
		q.rearNodeIndex = 0
		q.frontNode = newArrayNode(q.nodeSize, nil)
		q.rearNode = q.frontNode
	}

	q.listSize++
	q.rearNode.block[q.rearNodeIndex] = item
	q.rearNodeIndex++

	if q.rearNodeIndex == q.nodeSize {
		q.rearNodeIndex = 0
		q.rearNode.next = newArrayNode(q.nodeSize, nil)
		q.rearNode = q.rearNode.next
	}
}

func (q *arrayQueue) Dequeue() Value {
	if q.listSize == 0 {
		return nil
	}

	q.listSize--
	item := q.frontNode.block[q.frontNodeIndex]
	q.frontNodeIndex++

	if q.frontNodeIndex == q.nodeSize {
		q.frontNodeIndex = 0
		q.frontNode = q.frontNode.next
	}

	return item
}

func (q *arrayQueue) Peek() Value {
	if q.listSize == 0 {
		return nil
	}

	return q.frontNode.block[q.frontNodeIndex]
}

func (q *arrayQueue) Contains(item Value) bool {
	n := q.frontNode
	i := q.frontNodeIndex

	for n != nil && (n != q.rearNode || i <= q.rearNodeIndex) {
		if n.block[i].Compare(item) == 0 {
			return true
		}

		i++
		if i == q.nodeSize {
			n = n.next
			i = 0
		}
	}

	return false
}
