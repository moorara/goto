package sort

import (
	. "github.com/moorara/go-box/ds"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func merge(a, aux []Generic, lo, mid, hi int, cmp Comparator) {
	var i, j int = lo, mid + 1
	copy(aux[lo:hi+1], a[lo:hi+1])
	for k := lo; k <= hi; k++ {
		switch {
		case i > mid:
			a[k] = aux[j]
			j++
		case j > hi:
			a[k] = aux[i]
			i++
		case cmp.Compare(aux[j], aux[i]) < 0:
			a[k] = aux[j]
			j++
		default:
			a[k] = aux[i]
			i++
		}
	}
}

// MergeSort implements merge sort algorithm in an iterative manner
func MergeSort(a []Generic, cmp Comparator) {
	n := len(a)
	aux := make([]Generic, n)
	for sz := 1; sz < n; sz += sz {
		for lo := 0; lo < n-sz; lo += sz + sz {
			merge(a, aux, lo, lo+sz-1, min(lo+sz+sz-1, n-1), cmp)
		}
	}
}

func mergeSortRec(a, aux []Generic, lo, hi int, cmp Comparator) {
	if hi <= lo {
		return
	}

	mid := lo + (hi-lo)/2
	mergeSortRec(a, aux, lo, mid, cmp)
	mergeSortRec(a, aux, mid+1, hi, cmp)
	if cmp.Compare(a[mid+1], a[mid]) >= 0 {
		return
	}
	merge(a, aux, lo, mid, hi, cmp)
}

// MergeSortRec implements merge sort algorithm in a recursive manner
func MergeSortRec(a []Generic, cmp Comparator) {
	n := len(a)
	aux := make([]Generic, n)

	mergeSortRec(a, aux, 0, n-1, cmp)
}
