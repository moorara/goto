package sort

import (
	"testing"

	. "github.com/moorara/goto/dt"
	"github.com/moorara/goto/util"
	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	tests := []struct {
		compare       Compare
		items         []Generic
		expectedItems []Generic
	}{
		{CompareInt, []Generic{}, nil},
		{CompareInt, []Generic{20, 10, 30}, []Generic{10, 20, 30}},
		{CompareInt, []Generic{20, 10, 30, 40, 50}, []Generic{10, 20, 30, 40, 50}},
		{CompareInt, []Generic{20, 10, 30, 40, 50, 80, 60, 70, 90}, []Generic{10, 20, 30, 40, 50, 60, 70, 80, 90}},
	}

	for _, tc := range tests {
		for k := 0; k < len(tc.items); k++ {
			item := Select(tc.items, k, tc.compare)
			assert.Equal(t, tc.expectedItems[k], item)
		}
	}
}

func TestQuickSortInt(t *testing.T) {
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
		QuickSort(tc.items, tc.compare)
		assert.True(t, util.IsSorted(tc.items, tc.compare))
	}
}

func TestQuickSort3WayInt(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareInt, []Generic{}},
		{CompareInt, []Generic{20, 10, 10, 20, 30, 30, 30}},
		{CompareInt, []Generic{30, 20, 30, 20, 10, 40, 40, 40, 50, 50}},
		{CompareInt, []Generic{90, 10, 80, 20, 70, 30, 60, 40, 50, 50, 40, 60, 30, 70, 20, 80, 10, 90}},
	}

	for _, tc := range tests {
		QuickSort3Way(tc.items, tc.compare)
		assert.True(t, util.IsSorted(tc.items, tc.compare))
	}
}

func TestQuickSortString(t *testing.T) {
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
		QuickSort(tc.items, tc.compare)
		assert.True(t, util.IsSorted(tc.items, tc.compare))
	}
}

func TestQuickSort3WayString(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareString, []Generic{}},
		{CompareString, []Generic{"Milad", "Mona", "Milad", "Mona"}},
		{CompareString, []Generic{"Alice", "Bob", "Alex", "Jackie", "Jackie", "Alex", "Bob", "Alice"}},
		{CompareString, []Generic{"Docker", "Kubernetes", "Docker", "Go", "JavaScript", "Go", "React", "Redux", "Vue", "Redux", "React"}},
	}

	for _, tc := range tests {
		QuickSort3Way(tc.items, tc.compare)
		assert.True(t, util.IsSorted(tc.items, tc.compare))
	}
}
