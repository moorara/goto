package sort

import (
	"math/rand"
	"time"

	. "github.com/moorara/go-box/ds"
)

func toGenericArray(items ...interface{}) []Generic {
	genericItems := make([]Generic, len(items))
	for i := 0; i < len(genericItems); i++ {
		genericItems[i] = items[i]
	}

	return genericItems
}

func genGenericIntArray(size int) []Generic {
	rand.Seed(time.Now().UTC().UnixNano())
	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		items[i] = rand.Int()
	}

	return items
}

func genGenericStringArray(size, strMinLen, strMaxLen int) []Generic {
	rand.Seed(time.Now().UTC().UnixNano())
	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		strLen := strMinLen + rand.Intn(strMaxLen-strMinLen+1)
		bytes := make([]byte, strLen)
		for j := 0; j < strLen; j++ {
			bytes[j] = byte(65 + rand.Intn(90-65+1))
		}
		items[i] = string(bytes)
	}

	return items
}

func isSorted(items []Generic, cmp Comparator) bool {
	for i := 0; i < len(items)-1; i++ {
		if cmp.Compare(items[i], items[i+1]) > 0 {
			return false
		}
	}

	return true
}