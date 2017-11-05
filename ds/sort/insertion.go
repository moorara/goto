package sort

import (
	. "github.com/moorara/go-box/ds"
)

// InsertionSort implements insertion sort algorithm
func InsertionSort(a []Generic, cmp Comparator) {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := i; j > 0 && cmp.Compare(a[j], a[j-1]) < 0; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
