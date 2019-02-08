package st

import (
	"math/rand"
	"testing"

	"github.com/moorara/goto/sort"
)

const (
	seed = 27
)

func getIntSlice(size int) []interface{} {
	items := make([]interface{}, size)
	for i := 0; i < len(items); i++ {
		items[i] = i
	}
	sort.Shuffle(items)

	return items
}

func runPutBenchmark(b *testing.B, ost OrderedSymbolTable) {
	items := getIntSlice(b.N)
	rand.Seed(seed)
	sort.Shuffle(items)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ost.Put(items[n], "")
	}
}

func runGetBenchmark(b *testing.B, ost OrderedSymbolTable) {
	items := getIntSlice(b.N)
	rand.Seed(seed)
	sort.Shuffle(items)
	for n := 0; n < b.N; n++ {
		ost.Put(items[n], "")
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		ost.Get(items[n])
	}
}

func BenchmarkOrderedSymbolTable(b *testing.B) {
	b.Run("BST.Put", func(b *testing.B) {
		ost := NewBST(compareInt)
		runPutBenchmark(b, ost)
	})

	b.Run("BST.Get", func(b *testing.B) {
		ost := NewBST(compareInt)
		runGetBenchmark(b, ost)
	})

	b.Run("AVL.Put", func(b *testing.B) {
		ost := NewAVL(compareInt)
		runPutBenchmark(b, ost)
	})

	b.Run("AVL.Get", func(b *testing.B) {
		ost := NewAVL(compareInt)
		runGetBenchmark(b, ost)
	})

	b.Run("RedBlack.Put", func(b *testing.B) {
		ost := NewRedBlack(compareInt)
		runPutBenchmark(b, ost)
	})

	b.Run("RedBlack.Get", func(b *testing.B) {
		ost := NewRedBlack(compareInt)
		runGetBenchmark(b, ost)
	})
}
