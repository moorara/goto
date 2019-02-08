package sort

import (
	"math/rand"
)

// Shuffle shuffles an array in O(n) time
func Shuffle(a []interface{}) {
	n := len(a)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[i], a[r] = a[r], a[i]
	}
}
