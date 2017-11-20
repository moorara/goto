package st

import (
	. "github.com/moorara/go-box/dt"
)

type patriciaNode struct {
	key   Generic
	value Generic
	left  *rbNode
	right *rbNode
	b     int
}

type patricia struct {
	root         *patricia
	size         int
	bitStringKey BitString
}

func newPatriciaNode(key, value Generic, b int) *patriciaNode {
	return &patriciaNode{
		key:   key,
		value: value,
		b:     b,
	}
}

// NewPatricia creates a new Patricia (Radix) Tree
func NewPatricia(bitStringKey BitString) OrderedSymbolTable {
	return &patricia{
		root:         nil,
		size:         0,
		bitStringKey: bitStringKey,
	}
}

func (t *patricia) verify() bool {
	return false
}

func (t *patricia) Size() int {
	return 0
}

func (t *patricia) Height() int {
	return 0
}

func (t *patricia) IsEmpty() bool {
	return false
}

func (t *patricia) Put(key, value Generic) {

}

func (t *patricia) Get(key Generic) (Generic, bool) {
	return nil, false
}

func (t *patricia) Delete(key Generic) (Generic, bool) {
	return nil, false
}

func (t *patricia) KeyValues() []KeyValue {
	return nil
}

func (t *patricia) Min() (Generic, Generic) {
	return nil, nil
}

func (t *patricia) Max() (Generic, Generic) {
	return nil, nil
}

func (t *patricia) Floor(key Generic) (Generic, Generic) {
	return nil, nil
}

func (t *patricia) Ceiling(key Generic) (Generic, Generic) {
	return nil, nil
}

func (t *patricia) Rank(key Generic) int {
	return 0
}

func (t *patricia) Select(rank int) (Generic, Generic) {
	return nil, nil
}

func (t *patricia) DeleteMin() (Generic, Generic) {
	return nil, nil
}

func (t *patricia) DeleteMax() (Generic, Generic) {
	return nil, nil
}

func (t *patricia) RangeSize(lo, hi Generic) int {
	return 0
}

func (t *patricia) Range(lo, hi Generic) []KeyValue {
	return nil
}

func (t *patricia) Traverse(order int, visit VisitFunc) {

}

func (t *patricia) Graphviz() string {
	return ""
}
