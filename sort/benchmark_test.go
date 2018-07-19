package sort

import (
	"math/rand"
	"sort"
	"testing"

	. "github.com/moorara/goto/dt"
	"github.com/moorara/goto/math"
)

const (
	seed   = 27
	size   = 1000
	minInt = 0
	maxInt = 1000000
)

type GenricSlice []Generic

func (s GenricSlice) Len() int {
	return len(s)
}

func (s GenricSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GenricSlice) Less(i, j int) bool {
	return CompareInt(s[i], s[j]) < 0
}

func BenchmarkSort(b *testing.B) {
	b.Run("sort.Sort", func(b *testing.B) {
		rand.Seed(seed)
		items := GenricSlice(math.GenerateIntSlice(size, minInt, maxInt))
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			sort.Sort(items)
		}
	})

	b.Run("HeapSort", func(b *testing.B) {
		rand.Seed(seed)
		items := math.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			HeapSort(items, CompareInt)
		}
	})

	b.Run("InsertionSort", func(b *testing.B) {
		rand.Seed(seed)
		items := math.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			InsertionSort(items, CompareInt)
		}
	})

	b.Run("MergeSort", func(b *testing.B) {
		rand.Seed(seed)
		items := math.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			MergeSort(items, CompareInt)
		}
	})

	b.Run("MergeSortRec", func(b *testing.B) {
		rand.Seed(seed)
		items := math.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			MergeSortRec(items, CompareInt)
		}
	})

	b.Run("QuickSort", func(b *testing.B) {
		rand.Seed(seed)
		items := math.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			QuickSort(items, CompareInt)
		}
	})

	b.Run("QuickSort3Way", func(b *testing.B) {
		rand.Seed(seed)
		items := math.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			QuickSort3Way(items, CompareInt)
		}
	})

	b.Run("ShellSort", func(b *testing.B) {
		rand.Seed(seed)
		items := math.GenerateIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			math.Shuffle(items)
			ShellSort(items, CompareInt)
		}
	})
}
