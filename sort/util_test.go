package sort

import (
	"testing"

	. "github.com/moorara/go-box/dt"
	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	tests := []struct {
		compare Compare
		items   []Generic
	}{
		{CompareInt, []Generic{10, 15, 20, 25, 30, 35, 40, 45, 50}},
		{CompareInt, []Generic{10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95}},
	}

	for _, test := range tests {
		Shuffle(test.items)

		assert.False(t, isSorted(test.items, test.compare))
	}
}

func TestSelect(t *testing.T) {
	tests := []struct {
		compare       Compare
		items         []Generic
		expectedItems []Generic
	}{
		{CompareInt, []Generic{}, nil},
		{CompareInt, []Generic{20, 10, 30}, []Generic{10, 20, 30}},
		{CompareInt, []Generic{20, 10, 30, 40, 50}, []Generic{10, 20, 30, 40, 50}},
		{CompareInt, []Generic{20, 10, 30, 40, 50, 80, 60, 70, 90}, []Generic{10, 20, 30, 40, 50, 60, 70, 80, 90}},
	}

	for _, test := range tests {
		for k := 0; k < len(test.items); k++ {
			item := Select(test.items, k, test.compare)

			assert.Equal(t, test.expectedItems[k], item)
		}
	}
}
