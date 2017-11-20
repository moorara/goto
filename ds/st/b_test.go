package st

import (
	"testing"
)

const (
	bTreeOrder = 4
)

func getBTreeTests() []orderedSymbolTableTest {
	tests := getOrderedSymbolTableTests()

	tests[0].SymbolTable = "B-Tree"
	tests[0].expectedHeight = 0
	tests[0].expectedPreOrderTraverse = nil
	tests[0].expectedInOrderTraverse = nil
	tests[0].expectedPostOrderTraverse = nil
	tests[0].expectedDotCode = ``

	tests[1].SymbolTable = "B-Tree"
	tests[1].expectedHeight = 0
	tests[1].expectedPreOrderTraverse = []KeyValue{}
	tests[1].expectedInOrderTraverse = []KeyValue{}
	tests[1].expectedPostOrderTraverse = []KeyValue{}
	tests[1].expectedDotCode = ``

	tests[2].SymbolTable = "B-Tree"
	tests[2].expectedHeight = 0
	tests[2].expectedPreOrderTraverse = []KeyValue{}
	tests[2].expectedInOrderTraverse = []KeyValue{}
	tests[2].expectedPostOrderTraverse = []KeyValue{}
	tests[2].expectedDotCode = ``

	tests[3].SymbolTable = "B-Tree"
	tests[3].expectedHeight = 0
	tests[3].expectedPreOrderTraverse = []KeyValue{}
	tests[3].expectedInOrderTraverse = []KeyValue{}
	tests[3].expectedPostOrderTraverse = []KeyValue{}
	tests[3].expectedDotCode = ``

	return tests
}

func TestBTree(t *testing.T) {
	tests := getBTreeTests()

	for _, test := range tests {
		bt := NewBTree(bTreeOrder, test.compareKey)
		runOrderedSymbolTableTest(t, bt, test)
	}
}
