package util

import (
	"math/rand"
	"time"

	. "github.com/moorara/go-box/dt"
)

// SeedWithNow seeds the random generator with time now
func SeedWithNow() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Shuffle shuffles an array in O(n) time
func Shuffle(a []Generic) {
	n := len(a)
	for i := 0; i < n; i++ {
		r := i + rand.Intn(n-i)
		a[i], a[r] = a[r], a[i]
	}
}

// GenerateInt generates a random integer
func GenerateInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

// GenerateString generates a random string
func GenerateString(minLen, maxLen int) string {
	strLen := minLen + rand.Intn(maxLen-minLen+1)
	bytes := make([]byte, strLen)
	for j := 0; j < strLen; j++ {
		bytes[j] = byte(65 + rand.Intn(90-65+1))
	}

	return string(bytes)
}

// GenerateIntArray generates an array with random integers
func GenerateIntArray(size, min, max int) []Generic {
	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		items[i] = min + rand.Intn(max-min+1)
	}

	return items
}

// GenerateStringArray generates an array with random strings
func GenerateStringArray(size, minLen, maxLen int) []Generic {
	items := make([]Generic, size)
	for i := 0; i < len(items); i++ {
		items[i] = GenerateString(minLen, maxLen)
	}

	return items
}
