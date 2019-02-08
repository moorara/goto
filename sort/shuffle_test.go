package sort

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	tests := []struct {
		items []interface{}
	}{
		{[]interface{}{10, 20, 30, 40, 50, 60, 70, 80, 90}},
		{[]interface{}{"Alice", "Bob", "Dan", "Edgar", "Helen", "Karen", "Milad", "Peter", "Sam", "Wesley"}},
	}

	rand.Seed(time.Now().UTC().UnixNano())

	for _, tc := range tests {
		orig := make([]interface{}, len(tc.items))
		copy(orig, tc.items)
		Shuffle(tc.items)

		assert.NotEqual(t, orig, tc.items)
	}
}
