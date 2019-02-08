package sort

import (
	"math/rand"
	"sort"
	"testing"
)

const (
	seed   = 27
	size   = 1000
	minInt = 0
	maxInt = 1000000
)

type GenricSlice []interface{}

func (s GenricSlice) Len() int {
	return len(s)
}

func (s GenricSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GenricSlice) Less(i, j int) bool {
	return compareInt(s[i], s[j]) < 0
}

func BenchmarkSort(b *testing.B) {
	b.Run("sort.Sort", func(b *testing.B) {
		rand.Seed(seed)
		items := GenricSlice(genIntSlice(size, minInt, maxInt))
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			sort.Sort(items)
		}
	})

	b.Run("HeapSort", func(b *testing.B) {
		rand.Seed(seed)
		items := genIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			HeapSort(items, compareInt)
		}
	})

	b.Run("InsertionSort", func(b *testing.B) {
		rand.Seed(seed)
		items := genIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			InsertionSort(items, compareInt)
		}
	})

	b.Run("MergeSort", func(b *testing.B) {
		rand.Seed(seed)
		items := genIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			MergeSort(items, compareInt)
		}
	})

	b.Run("MergeSortRec", func(b *testing.B) {
		rand.Seed(seed)
		items := genIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			MergeSortRec(items, compareInt)
		}
	})

	b.Run("QuickSort", func(b *testing.B) {
		rand.Seed(seed)
		items := genIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			QuickSort(items, compareInt)
		}
	})

	b.Run("QuickSort3Way", func(b *testing.B) {
		rand.Seed(seed)
		items := genIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			QuickSort3Way(items, compareInt)
		}
	})

	b.Run("ShellSort", func(b *testing.B) {
		rand.Seed(seed)
		items := genIntSlice(size, minInt, maxInt)
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			Shuffle(items)
			ShellSort(items, compareInt)
		}
	})
}
