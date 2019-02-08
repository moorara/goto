package heap

import (
	"math/rand"
	"strings"
)

func compareInt(a, b interface{}) int {
	intA, _ := a.(int)
	intB, _ := b.(int)
	diff := intA - intB
	switch {
	case diff < 0:
		return -1
	case diff > 0:
		return 1
	default:
		return 0
	}
}

func compareString(a, b interface{}) int {
	strA, _ := a.(string)
	strB, _ := b.(string)
	return strings.Compare(strA, strB)
}

func genIntSlice(size, min, max int) []interface{} {
	items := make([]interface{}, size)
	for i := 0; i < len(items); i++ {
		items[i] = min + rand.Intn(max-min+1)
	}

	return items
}
