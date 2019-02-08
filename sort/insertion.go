package sort

// InsertionSort implements insertion sort algorithm
func InsertionSort(a []interface{}, compare func(a, b interface{}) int) {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := i; j > 0 && compare(a[j], a[j-1]) < 0; j-- {
			a[j], a[j-1] = a[j-1], a[j]
		}
	}
}
