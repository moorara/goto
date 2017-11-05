package sort

import (
	"testing"

	. "github.com/moorara/go-box/ds"
	"github.com/stretchr/testify/assert"
)

func TestMergeSortInt(t *testing.T) {
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
		MergeSort(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func TestMergeSortString(t *testing.T) {
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
		MergeSort(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func TestMergeSortRecInt(t *testing.T) {
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
		MergeSortRec(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func TestMergeSortRecString(t *testing.T) {
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
		MergeSortRec(test.items, test.cmp)

		assert.True(t, isSorted(test.items, test.cmp))
	}
}

func BenchmarkMergeSortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		MergeSort(items, &IntComparator{})
	}
}

func BenchmarkMergeSortString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		MergeSort(items, &StringComparator{})
	}
}

func BenchmarkMergeSortRecInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		MergeSortRec(items, &IntComparator{})
	}
}

func BenchmarkMergeSortRecString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		MergeSortRec(items, &StringComparator{})
	}
}
