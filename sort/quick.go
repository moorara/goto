package sort

import (
	. "github.com/moorara/go-box/dt"
)

func partition(a []Generic, lo, hi int, compare Compare) int {
	v := a[lo]
	var i, j int = lo, hi + 1

	for true {
		for i++; i < hi && compare(a[i], v) < 0; i++ {
		}
		for j--; j > lo && compare(a[j], v) > 0; j-- {
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
func quickSort(a []Generic, lo, hi int, compare Compare) {
	if lo >= hi {
		return
	}

	j := partition(a, lo, hi, compare)
	quickSort(a, lo, j-1, compare)
	quickSort(a, j+1, hi, compare)
}

// QuickSort implements quick sort algorithm
func QuickSort(a []Generic, compare Compare) {
	Shuffle(a)
	quickSort(a, 0, len(a)-1, compare)
}

func quickSort3Way(a []Generic, lo, hi int, compare Compare) {
	if lo >= hi {
		return
	}

	v := a[lo]
	var lt, i, gt int = lo, lo + 1, hi

	for i <= gt {
		c := compare(a[i], v)
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

	quickSort3Way(a, lo, lt-1, compare)
	quickSort3Way(a, gt+1, hi, compare)
}

// QuickSort3Way implements 3-way quick sort algorithm
func QuickSort3Way(a []Generic, compare Compare) {
	quickSort3Way(a, 0, len(a)-1, compare)
}
