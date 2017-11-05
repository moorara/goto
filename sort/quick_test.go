package sort

import (
	"testing"

	. "github.com/moorara/go-box/ds"
	"github.com/stretchr/testify/assert"
)

func TestQuickSortInt(t *testing.T) {
	tests := []struct {
		cmp   Comparator
		items []Generic
	}{
		{&IntComparator{}, toGenericArray()},
		{&IntComparator{}, toGenericArray(20, 10, 30)},
		{&IntComparator{}, toGenericArray(30, 20, 10, 40, 50)},
		{&IntComparator{}, toGenericArray(90, 80, 70, 60, 50, 40, 30, 20, 10)},
	}

	for _, test := range tests {
		QuickSort(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func TestQuickSort3WayInt(t *testing.T) {
	tests := []struct {
		cmp   Comparator
		items []Generic
	}{
		{&IntComparator{}, toGenericArray()},
		{&IntComparator{}, toGenericArray(20, 10, 10, 20, 30, 30, 30)},
		{&IntComparator{}, toGenericArray(30, 20, 30, 20, 10, 40, 40, 40, 50, 50)},
		{&IntComparator{}, toGenericArray(90, 10, 80, 20, 70, 30, 60, 40, 50, 50, 40, 60, 30, 70, 20, 80, 10, 90)},
	}

	for _, test := range tests {
		QuickSort3Way(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func TestQuickSortString(t *testing.T) {
	tests := []struct {
		cmp   Comparator
		items []Generic
	}{
		{&StringComparator{}, toGenericArray()},
		{&StringComparator{}, toGenericArray("Milad", "Mona")},
		{&StringComparator{}, toGenericArray("Alice", "Bob", "Alex", "Jackie")},
		{&StringComparator{}, toGenericArray("Docker", "Kubernetes", "Go", "JavaScript", "Elixir", "React", "Redux", "Vue")},
	}

	for _, test := range tests {
		QuickSort(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func TestQuickSort3WayString(t *testing.T) {
	tests := []struct {
		cmp   Comparator
		items []Generic
	}{
		{&StringComparator{}, toGenericArray()},
		{&StringComparator{}, toGenericArray("Milad", "Mona", "Milad", "Mona")},
		{&StringComparator{}, toGenericArray("Alice", "Bob", "Alex", "Jackie", "Jackie", "Alex", "Bob", "Alice")},
		{&StringComparator{}, toGenericArray("Docker", "Kubernetes", "Docker", "Go", "JavaScript", "Go", "React", "Redux", "Vue", "Redux", "React")},
	}

	for _, test := range tests {
		QuickSort3Way(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func BenchmarkQuickSortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		QuickSort(items, &IntComparator{})
	}
}

func BenchmarkQuick3WaySortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		QuickSort3Way(items, &IntComparator{})
	}
}

func BenchmarkQuickSortString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		QuickSort(items, &StringComparator{})
	}
}

func BenchmarkQuickSort3WayString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		QuickSort3Way(items, &StringComparator{})
	}
}
