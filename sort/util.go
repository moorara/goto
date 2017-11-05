package sort

import (
	"math/rand"
	"time"

	. "github.com/moorara/go-box/ds"
)

// Shuffle shuffles an array in O(n) time
func Shuffle(a []Generic) {
	rand.Seed(time.Now().UTC().UnixNano())
	n := len(a)

	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[i], a[r] = a[r], a[i]
	}
}

// Select finds the kth smallest item of an array in O(n) time on average
func Select(a []Generic, k int, cmp Comparator) Generic {
	Shuffle(a)
	var lo, hi int = 0, len(a) - 1
	for lo < hi {
		j := partition(a, lo, hi, cmp)
		switch {
		case j < k:
			lo = j + 1
		case j > k:
			hi = j - 1
		default:
			return a[k]
		}
	}

	return a[k]
}
