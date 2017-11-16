package sort

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/stretchr/testify/assert"
)

func TestMergeSortInt(t *testing.T) {
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
		MergeSort(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
	}
}

func TestMergeSortString(t *testing.T) {
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
		MergeSort(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
	}
}

func TestMergeSortRecInt(t *testing.T) {
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
		MergeSortRec(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
	}
}

func TestMergeSortRecString(t *testing.T) {
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
		MergeSortRec(test.items, test.compare)

		assert.True(t, isSorted(test.items, test.compare))
	}
}

func BenchmarkMergeSortInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		MergeSort(items, CompareInt)
	}
}

func BenchmarkMergeSortString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		MergeSort(items, CompareString)
	}
}

func BenchmarkMergeSortRecInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericIntArray(1000)
		MergeSortRec(items, CompareInt)
	}
}

func BenchmarkMergeSortRecString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		items := genGenericStringArray(1000, 10, 20)
		MergeSortRec(items, CompareString)
	}
}
