package st

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	symbolTableTest struct {
		name            string
		symbolTable     string
		compareKey      func(a, b interface{}) int
		keyValues       []KeyValue
		expectedSize    int
		expectedIsEmpty bool
	}

	orderedSymbolTableTest struct {
		name                      string
		symbolTable               string
		compareKey                func(a, b interface{}) int
		keyValues                 []KeyValue
		expectedSize              int
		expectedHeight            int
		expectedIsEmpty           bool
		expectedMinKey            interface{}
		expectedMinValue          interface{}
		expectedMaxKey            interface{}
		expectedMaxValue          interface{}
		floorKey                  string
		expectedFloorKey          interface{}
		expectedFloorValue        interface{}
		ceilingKey                string
		expectedCeilingKey        interface{}
		expectedCeilingValue      interface{}
		rankKey                   string
		expectedRank              int
		selectRank                int
		expectedSelectKey         interface{}
		expectedSelectValue       interface{}
		rangeKeyLo                string
		rangeKeyHi                string
		expectedRangeSize         int
		expectedRange             []KeyValue
		expectedPreOrderTraverse  []KeyValue
		expectedInOrderTraverse   []KeyValue
		expectedPostOrderTraverse []KeyValue
		expectedDotCode           string
	}
)

func getSymbolTableTests() []symbolTableTest {
	return []symbolTableTest{
		{
			name:            "",
			symbolTable:     "",
			compareKey:      compareString,
			keyValues:       []KeyValue{},
			expectedSize:    0,
			expectedIsEmpty: true,
		},
	}
}

func getOrderedSymbolTableTests() []orderedSymbolTableTest {
	return []orderedSymbolTableTest{
		{
			name:                 "Empty",
			compareKey:           compareString,
			keyValues:            []KeyValue{},
			expectedSize:         0,
			expectedIsEmpty:      true,
			expectedMinKey:       nil,
			expectedMinValue:     nil,
			expectedMaxKey:       nil,
			expectedMaxValue:     nil,
			floorKey:             "",
			expectedFloorKey:     nil,
			expectedFloorValue:   nil,
			ceilingKey:           "",
			expectedCeilingKey:   nil,
			expectedCeilingValue: nil,
			rankKey:              "",
			expectedRank:         0,
			selectRank:           0,
			expectedSelectKey:    nil,
			expectedSelectValue:  nil,
			rangeKeyLo:           "",
			rangeKeyHi:           "",
			expectedRangeSize:    0,
			expectedRange:        nil,
		},
		{
			name:       "ABC",
			compareKey: compareString,
			keyValues: []KeyValue{
				{"B", 2},
				{"A", 1},
				{"C", 3},
			},
			expectedSize:         3,
			expectedIsEmpty:      false,
			expectedMinKey:       "A",
			expectedMinValue:     1,
			expectedMaxKey:       "C",
			expectedMaxValue:     3,
			floorKey:             "A",
			expectedFloorKey:     "A",
			expectedFloorValue:   1,
			ceilingKey:           "C",
			expectedCeilingKey:   "C",
			expectedCeilingValue: 3,
			rankKey:              "C",
			expectedRank:         2,
			selectRank:           1,
			expectedSelectKey:    "B",
			expectedSelectValue:  2,
			rangeKeyLo:           "A",
			rangeKeyHi:           "C",
			expectedRangeSize:    3,
			expectedRange: []KeyValue{
				{"A", 1},
				{"B", 2},
				{"C", 3},
			},
		},
		{
			name:       "ABCDE",
			compareKey: compareString,
			keyValues: []KeyValue{
				{"B", 2},
				{"A", 1},
				{"C", 3},
				{"E", 5},
				{"D", 4},
			},
			expectedSize:         5,
			expectedIsEmpty:      false,
			expectedMinKey:       "A",
			expectedMinValue:     1,
			expectedMaxKey:       "E",
			expectedMaxValue:     5,
			floorKey:             "B",
			expectedFloorKey:     "B",
			expectedFloorValue:   2,
			ceilingKey:           "D",
			expectedCeilingKey:   "D",
			expectedCeilingValue: 4,
			rankKey:              "E",
			expectedRank:         4,
			selectRank:           2,
			expectedSelectKey:    "C",
			expectedSelectValue:  3,
			rangeKeyLo:           "B",
			rangeKeyHi:           "D",
			expectedRangeSize:    3,
			expectedRange: []KeyValue{
				{"B", 2},
				{"C", 3},
				{"D", 4},
			},
		},
		{
			name:       "ADGJMPS",
			compareKey: compareString,
			keyValues: []KeyValue{
				{"J", 10},
				{"A", 1},
				{"D", 4},
				{"S", 19},
				{"P", 16},
				{"M", 13},
				{"G", 7},
			},
			expectedSize:         7,
			expectedIsEmpty:      false,
			expectedMinKey:       "A",
			expectedMinValue:     1,
			expectedMaxKey:       "S",
			expectedMaxValue:     19,
			floorKey:             "C",
			expectedFloorKey:     "A",
			expectedFloorValue:   1,
			ceilingKey:           "R",
			expectedCeilingKey:   "S",
			expectedCeilingValue: 19,
			rankKey:              "S",
			expectedRank:         6,
			selectRank:           3,
			expectedSelectKey:    "J",
			expectedSelectValue:  10,
			rangeKeyLo:           "B",
			rangeKeyHi:           "R",
			expectedRangeSize:    5,
			expectedRange: []KeyValue{
				{"D", 4},
				{"G", 7},
				{"J", 10},
				{"M", 13},
				{"P", 16},
			},
		},
	}
}

func runSymbolTableTest(t *testing.T, st SymbolTable, test symbolTableTest) {
	t.Run(test.name, func(t *testing.T) {
		// Tree initially should be empty
		assert.True(t, st.IsEmpty())
		assert.Zero(t, st.Size())

		// TODO: verify
		assert.NotEmpty(t, test.symbolTable)
		assert.NotEmpty(t, test.compareKey)
		assert.NotEmpty(t, test.keyValues)
		assert.NotEmpty(t, test.expectedSize)
		assert.NotEmpty(t, test.expectedIsEmpty)

		// Tree should be empty at the end
		assert.Zero(t, st.Size())
		assert.True(t, st.IsEmpty())
	})
}

func runOrderedSymbolTableTest(t *testing.T, ost OrderedSymbolTable, test orderedSymbolTableTest) {
	t.Run(test.name, func(t *testing.T) {
		var i int
		var kvs []KeyValue
		var minKey, minValue interface{}
		var maxKey, maxValue interface{}
		var floorKey, floorValue interface{}
		var ceilingKey, ceilingValue interface{}
		var selectKey, selectValue interface{}

		// Tree initially should be empty
		assert.True(t, ost.verify())
		assert.Zero(t, ost.Size())
		assert.Zero(t, ost.Height())
		assert.True(t, ost.IsEmpty())
		minKey, minValue = ost.Min()
		assert.Nil(t, minKey)
		assert.Nil(t, minValue)
		maxKey, maxValue = ost.Max()
		assert.Nil(t, maxKey)
		assert.Nil(t, maxValue)
		floorKey, floorValue = ost.Floor(nil)
		assert.Nil(t, floorKey)
		assert.Nil(t, floorValue)
		ceilingKey, ceilingValue = ost.Ceiling(nil)
		assert.Nil(t, ceilingKey)
		assert.Nil(t, ceilingValue)
		assert.Equal(t, -1, ost.Rank(nil))
		selectKey, selectValue = ost.Select(0)
		assert.Nil(t, selectKey)
		assert.Nil(t, selectValue)
		assert.Equal(t, -1, ost.RangeSize(nil, nil))
		assert.Nil(t, ost.Range(nil, nil))

		// Put
		for _, kv := range test.keyValues {
			ost.Put(kv.key, kv.value)
			ost.Put(kv.key, kv.value) // Update existing key-value
			assert.True(t, ost.verify())
		}

		// Get
		for _, expected := range test.keyValues {
			value, ok := ost.Get(expected.key)
			assert.True(t, ok)
			assert.Equal(t, expected.value, value)
		}

		assert.Equal(t, test.expectedSize, ost.Size())
		assert.Equal(t, test.expectedHeight, ost.Height())
		assert.Equal(t, test.expectedIsEmpty, ost.IsEmpty())
		minKey, minValue = ost.Min()
		assert.Equal(t, test.expectedMinKey, minKey)
		assert.Equal(t, test.expectedMinValue, minValue)
		maxKey, maxValue = ost.Max()
		assert.Equal(t, test.expectedMaxKey, maxKey)
		assert.Equal(t, test.expectedMaxValue, maxValue)
		floorKey, floorValue = ost.Floor(test.floorKey)
		assert.Equal(t, test.expectedFloorKey, floorKey)
		assert.Equal(t, test.expectedFloorValue, floorValue)
		ceilingKey, ceilingValue = ost.Ceiling(test.ceilingKey)
		assert.Equal(t, test.expectedCeilingKey, ceilingKey)
		assert.Equal(t, test.expectedCeilingValue, ceilingValue)
		assert.Equal(t, test.expectedRank, ost.Rank(test.rankKey))
		selectKey, selectValue = ost.Select(test.selectRank)
		assert.Equal(t, test.expectedSelectKey, selectKey)
		assert.Equal(t, test.expectedSelectValue, selectValue)

		minKey, minValue = ost.DeleteMin()
		assert.Equal(t, test.expectedMinKey, minKey)
		assert.Equal(t, test.expectedMinValue, minValue)
		assert.True(t, ost.verify())
		ost.Put(minKey, minValue)
		maxKey, maxValue = ost.DeleteMax()
		assert.Equal(t, test.expectedMaxKey, maxKey)
		assert.Equal(t, test.expectedMaxValue, maxValue)
		assert.True(t, ost.verify())
		ost.Put(maxKey, maxValue)

		kvs = ost.KeyValues()
		for _, kv := range kvs { // Soundness
			assert.Contains(t, test.keyValues, kv)
		}
		for _, kv := range test.keyValues { // Completeness
			assert.Contains(t, kvs, kv)
		}
		for i = 0; i < len(kvs)-1; i++ { // Sorted Ascending
			assert.Equal(t, -1, test.compareKey(kvs[i].key, kvs[i+1].key))
		}

		assert.Equal(t, test.expectedRangeSize, ost.RangeSize(test.rangeKeyLo, test.rangeKeyHi))
		kvs = ost.Range(test.rangeKeyLo, test.rangeKeyHi)
		for _, kv := range kvs { // Soundness
			assert.Contains(t, test.expectedRange, kv)
		}
		for _, kv := range test.expectedRange { // Completeness
			assert.Contains(t, kvs, kv)
		}
		for i = 0; i < len(kvs)-1; i++ { // Sorted Ascending
			assert.Equal(t, -1, test.compareKey(kvs[i].key, kvs[i+1].key))
		}

		// Invalid Traversal
		i = 0
		ost.Traverse(-1, func(key, value interface{}) bool {
			i++
			return true
		})
		assert.Zero(t, i)

		// Pre-Order Traversal
		i = 0
		ost.Traverse(TraversePreOrder, func(key, value interface{}) bool {
			assert.Equal(t, test.expectedPreOrderTraverse[i].key, key)
			assert.Equal(t, test.expectedPreOrderTraverse[i].value, value)
			i++
			return true
		})

		// In-Order Traversal
		i = 0
		ost.Traverse(TraverseInOrder, func(key, value interface{}) bool {
			assert.Equal(t, test.expectedInOrderTraverse[i].key, key)
			assert.Equal(t, test.expectedInOrderTraverse[i].value, value)
			i++
			return true
		})

		// Post-Order Traversal
		i = 0
		ost.Traverse(TraversePostOrder, func(key, value interface{}) bool {
			assert.Equal(t, test.expectedPostOrderTraverse[i].key, key)
			assert.Equal(t, test.expectedPostOrderTraverse[i].value, value)
			i++
			return true
		})

		// Graphviz dot language code
		assert.Equal(t, test.expectedDotCode, ost.Graphviz())

		// Delete
		value, ok := ost.Delete(nil)
		assert.False(t, ok)
		assert.Nil(t, value)
		for _, expected := range test.keyValues {
			value, ok := ost.Delete(expected.key)
			assert.True(t, ok)
			assert.Equal(t, expected.value, value)
			assert.True(t, ost.verify())
		}

		// Tree should be empty at the end
		assert.True(t, ost.verify())
		assert.Zero(t, ost.Size())
		assert.Zero(t, ost.Height())
		assert.True(t, ost.IsEmpty())
		minKey, minValue = ost.Min()
		assert.Nil(t, minKey)
		assert.Nil(t, minValue)
		maxKey, maxValue = ost.Max()
		assert.Nil(t, maxKey)
		assert.Nil(t, maxValue)
		floorKey, floorValue = ost.Floor(nil)
		assert.Nil(t, floorKey)
		assert.Nil(t, floorValue)
		ceilingKey, ceilingValue = ost.Ceiling(nil)
		assert.Nil(t, ceilingKey)
		assert.Nil(t, ceilingValue)
		assert.Equal(t, -1, ost.Rank(nil))
		selectKey, selectValue = ost.Select(0)
		assert.Nil(t, selectKey)
		assert.Nil(t, selectValue)
		assert.Equal(t, -1, ost.RangeSize(nil, nil))
		assert.Nil(t, ost.Range(nil, nil))
	})
}
