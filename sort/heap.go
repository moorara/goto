package sort

func sink(a []interface{}, k, n int, compare func(a, b interface{}) int) {
	for 2*k <= n {
		j := 2 * k
		if j < n && compare(a[j], a[j+1]) < 0 {
			j++
		}
		if compare(a[k], a[j]) >= 0 {
			break
		}
		a[k], a[j] = a[j], a[k]
		k = j
	}
}

func heapSort(a []interface{}, compare func(a, b interface{}) int) {
	n := len(a) - 1

	for k := n / 2; k >= 1; k-- { // build max-heap bottom-up
		sink(a, k, n, compare)
	}
	for n > 1 { // remove the maximum, one at a time
		a[1], a[n] = a[n], a[1]
		n--
		sink(a, 1, n, compare)
	}
}

// HeapSort implements heap sort algorithm
func HeapSort(a []interface{}, compare func(a, b interface{}) int) {
	// Heap elements need to start from position 1
	aux := append([]interface{}{nil}, a...)
	heapSort(aux, compare)
	copy(a, aux[1:])
}
