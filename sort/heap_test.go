package sort

import (
	"testing"

	. "github.com/moorara/goto/dt"
	"github.com/moorara/goto/util"
	"github.com/stretchr/testify/assert"
)

func TestHeapSortInt(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareInt, []Generic{}},
		{CompareInt, []Generic{20, 10, 30}},
		{CompareInt, []Generic{30, 20, 10, 40, 50}},
		{CompareInt, []Generic{90, 80, 70, 60, 50, 40, 30, 20, 10}},
	}

	for _, tc := range tests {
		HeapSort(tc.items, tc.compare)
		assert.True(t, util.IsSorted(tc.items, tc.compare))
	}
}

func TestHeapSortString(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareString, []Generic{}},
		{CompareString, []Generic{"Milad", "Mona"}},
		{CompareString, []Generic{"Alice", "Bob", "Alex", "Jackie"}},
		{CompareString, []Generic{"Docker", "Kubernetes", "Go", "JavaScript", "Elixir", "React", "Redux", "Vue"}},
	}

	for _, tc := range tests {
		HeapSort(tc.items, tc.compare)
		assert.True(t, util.IsSorted(tc.items, tc.compare))
	}
}
