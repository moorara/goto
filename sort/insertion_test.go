package sort

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/moorara/go-box/util"
	"github.com/stretchr/testify/assert"
)

func TestInsertionSortInt(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareInt, []Generic{}},
		{CompareInt, []Generic{20, 10, 30}},
		{CompareInt, []Generic{30, 20, 10, 40, 50}},
		{CompareInt, []Generic{90, 80, 70, 60, 50, 40, 30, 20, 10}},
	}

	for _, test := range tests {
		InsertionSort(test.items, test.compare)

		assert.True(t, util.IsSorted(test.items, test.compare))
	}
}

func TestInsertionSortString(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareString, []Generic{}},
		{CompareString, []Generic{"Milad", "Mona"}},
		{CompareString, []Generic{"Alice", "Bob", "Alex", "Jackie"}},
		{CompareString, []Generic{"Docker", "Kubernetes", "Go", "JavaScript", "Elixir", "React", "Redux", "Vue"}},
	}

	for _, test := range tests {
		InsertionSort(test.items, test.compare)

		assert.True(t, util.IsSorted(test.items, test.compare))
	}
}

func BenchmarkInsertionSortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := util.GenerateIntArray(1000, -1000, 1000)
		InsertionSort(items, CompareInt)
	}
}

func BenchmarkInsertionSortString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := util.GenerateStringArray(1000, 10, 50)
		InsertionSort(items, CompareString)
	}
}
