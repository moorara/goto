package util

import (
	"math/rand"
	"time"

	. "github.com/moorara/go-box/dt"
)

// Shuffle shuffles an array in O(n) time
func Shuffle(a []Generic) {
	rand.Seed(time.Now().UTC().UnixNano())

	n := len(a)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[i], a[r] = a[r], a[i]
	}
}

// GenerateIntArray generates an array with random integers
func GenerateIntArray(size, min, max int) []Generic {
	rand.Seed(time.Now().UTC().UnixNano())

	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		items[i] = min + rand.Intn(max-min+1)
	}

	return items
}

// GenerateStringArray generates an array with random strings
func GenerateStringArray(size, minLen, maxLen int) []Generic {
	rand.Seed(time.Now().UTC().UnixNano())

	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		strLen := minLen + rand.Intn(maxLen-minLen+1)
		bytes := make([]byte, strLen)
		for j := 0; j < strLen; j++ {
			bytes[j] = byte(65 + rand.Intn(90-65+1))
		}
		items[i] = string(bytes)
	}

	return items
}
