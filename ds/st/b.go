/*
 * B-Tree order M gengerilizes 2-3 Tree by allowing M-1 key-link pairs per node.
 *   - At least 2 key-link pairs at root
 *   - At least M/2 key-link pairs in other nodes
 *   - External nodes contain client keys
 *   - Internal nodes contain copy of client keys to guide search
 * Choose M as large as enough. so that M links fit in a page (e.g., M = 1024)
 */

package st

import (
	. "github.com/moorara/go-box/dt"
)

type btEntry struct {
	key   Generic
	value Generic
	next  *btNode
}

type btNode struct {
	m        int
	children []btEntry
}

type bTree struct {
	M          int
	root       *btNode
	size       int
	height     int
	compareKey Compare
}

func newBTEntry(key, value Generic, next *btNode) btEntry {
	return btEntry{
		key:   key,
		value: value,
		next:  next,
	}
}

func newBTNode(M, m int) *btNode {
	return &btNode{
		m:        m,
		children: make([]btEntry, M),
	}
}

// NewBTree creates a new B-Tree
func NewBTree(M int, compareKey Compare) OrderedSymbolTable {
	return &bTree{
		M:          M,
		root:       newBTNode(M, 0),
		size:       0,
		height:     0,
		compareKey: compareKey,
	}
}

func (t *bTree) verify() bool {
	return false
}

func (t *bTree) split(n *btNode) *btNode {
	n.m = t.M / 2
	x := newBTNode(t.M, t.M/2)
	for i := 0; i < x.m; i++ {
		x.children[i] = n.children[n.m+i]
	}

	return x
}

func (t *bTree) Size() int {
	return t.size
}

func (t *bTree) Height() int {
	return t.height
}

func (t *bTree) IsEmpty() bool {
	return t.Size() == 0
}

func (t *bTree) _put(n *btNode, key, value Generic, h int) *btNode {
	var i int
	e := newBTEntry(key, value, nil)

	if h == 0 { // External node
		for i = 0; i < n.m; i++ {
			if t.compareKey(key, n.children[i].key) < 0 {
				break
			}
		}
	} else { // Internal node
		for i = 0; i < n.m; i++ {
			if i+1 == n.m || t.compareKey(key, n.children[i+1].key) < 0 {
				x := t._put(n.children[i].next, key, value, h-1)
				if x == nil {
					return nil
				}
				i++
				e.key = x.children[0].key
				e.next = x
				break
			}
		}
	}

	for j := n.m; j > i; j-- {
		n.children[j] = n.children[j-1]
	}
	n.children[i] = e
	n.m++
	if n.m < t.M {
		return nil
	}
	return t.split(n)
}

func (t *bTree) Put(key, value Generic) {
	if key == nil {
		return
	}

	x := t._put(t.root, key, value, t.height)
	t.size++
	if x == nil {
		return
	}

	// Need to split root
	y := newBTNode(t.M, 2)
	y.children[0] = newBTEntry(t.root.children[0].key, nil, t.root)
	y.children[1] = newBTEntry(x.children[0].key, nil, x)
	t.root = y
	t.height++
}

func (t *bTree) Get(key Generic) (Generic, bool) {
	return nil, false
}

func (t *bTree) Delete(key Generic) (Generic, bool) {
	return nil, false
}

func (t *bTree) KeyValues() []KeyValue {
	return nil
}

func (t *bTree) Min() (Generic, Generic) {
	return nil, nil
}

func (t *bTree) Max() (Generic, Generic) {
	return nil, nil
}

func (t *bTree) Floor(key Generic) (Generic, Generic) {
	return nil, nil
}

func (t *bTree) Ceiling(key Generic) (Generic, Generic) {
	return nil, nil
}

func (t *bTree) Rank(key Generic) int {
	return 0
}

func (t *bTree) Select(rank int) (Generic, Generic) {
	return nil, nil
}

func (t *bTree) DeleteMin() (Generic, Generic) {
	return nil, nil
}

func (t *bTree) DeleteMax() (Generic, Generic) {
	return nil, nil
}

func (t *bTree) RangeSize(lo, hi Generic) int {
	return 0
}

func (t *bTree) Range(lo, hi Generic) []KeyValue {
	return nil
}

func (t *bTree) Traverse(order int, visit VisitFunc) {

}

func (t *bTree) Graphviz() string {
	return ""
}
