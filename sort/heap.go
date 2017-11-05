package sort

import (
	. "github.com/moorara/go-box/ds"
)

func sink(a []Generic, k, n int, cmp Comparator) {
	for 2*k <= n {
		j := 2 * k
		if j < n && cmp.Compare(a[j], a[j+1]) < 0 {
			j++
		}
		if cmp.Compare(a[k], a[j]) >= 0 {
			break
		}
		a[k], a[j] = a[j], a[k]
		k = j
	}
}

func heapSort(a []Generic, cmp Comparator) {
	n := len(a) - 1

	for k := n / 2; k >= 1; k-- { // build max-heap bottom-up
		sink(a, k, n, cmp)
	}
	for n > 1 { // remove the maximum, one at a time
		a[1], a[n] = a[n], a[1]
		n--
		sink(a, 1, n, cmp)
	}
}

// HeapSort implements heap sort algorithm
func HeapSort(a []Generic, cmp Comparator) {
	// Heap elements need to start from position 1
	aux := append([]Generic{nil}, a...)
	heapSort(aux, cmp)
	copy(a, aux[1:])
}
