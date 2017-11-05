package sort

import (
	. "github.com/moorara/go-box/ds"
)

func partition(a []Generic, lo, hi int, cmp Comparator) int {
	v := a[lo]
	var i, j int = lo, hi + 1

	for true {
		for i++; i < hi && cmp.Compare(a[i], v) < 0; i++ {
		}
		for j--; j > lo && cmp.Compare(a[j], v) > 0; j-- {
		}
		if i >= j {
			break
		}
		a[i], a[j] = a[j], a[i]
	}
	a[lo], a[j] = a[j], a[lo]

	return j
}

// QuickSort implements quick sort algorithm
func quickSort(a []Generic, lo, hi int, cmp Comparator) {
	if lo >= hi {
		return
	}

	j := partition(a, lo, hi, cmp)
	quickSort(a, lo, j-1, cmp)
	quickSort(a, j+1, hi, cmp)
}

// QuickSort implements quick sort algorithm
func QuickSort(a []Generic, cmp Comparator) {
	Shuffle(a)
	quickSort(a, 0, len(a)-1, cmp)
}

func quickSort3Way(a []Generic, lo, hi int, cmp Comparator) {
	if lo >= hi {
		return
	}

	v := a[lo]
	var lt, i, gt int = lo, lo + 1, hi

	for i <= gt {
		c := cmp.Compare(a[i], v)
		switch {
		case c < 0:
			a[lt], a[i] = a[i], a[lt]
			lt++
			i++
		case c > 0:
			a[i], a[gt] = a[gt], a[i]
			gt--
		default:
			i++
		}
	}

	quickSort3Way(a, lo, lt-1, cmp)
	quickSort3Way(a, gt+1, hi, cmp)
}

// QuickSort3Way implements 3-way quick sort algorithm
func QuickSort3Way(a []Generic, cmp Comparator) {
	quickSort3Way(a, 0, len(a)-1, cmp)
}
