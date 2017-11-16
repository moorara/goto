package sort

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/stretchr/testify/assert"
)

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

	for _, test := range tests {
		QuickSort(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
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

	for _, test := range tests {
		QuickSort3Way(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
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

	for _, test := range tests {
		QuickSort(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
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

	for _, test := range tests {
		QuickSort3Way(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
	}
}

func BenchmarkQuickSortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		QuickSort(items, CompareInt)
	}
}

func BenchmarkQuick3WaySortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		QuickSort3Way(items, CompareInt)
	}
}

func BenchmarkQuickSortString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		QuickSort(items, CompareString)
	}
}

func BenchmarkQuickSort3WayString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		QuickSort3Way(items, CompareString)
	}
}
