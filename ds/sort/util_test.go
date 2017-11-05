package sort

import (
	"testing"

	. "github.com/moorara/go-box/ds"
	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	tests := []struct {
		cmp   Comparator
		items []Generic
	}{
		{&IntComparator{}, toGenericArray(10, 15, 20, 25, 30, 35, 40, 45, 50)},
		{&IntComparator{}, toGenericArray(10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95)},
	}

	for _, test := range tests {
		Shuffle(test.items)

		assert.False(t, isSorted(test.items, test.cmp))
	}
}

func TestSelect(t *testing.T) {
	tests := []struct {
		cmp           Comparator
		items         []Generic
		expectedItems []Generic
	}{
		{&IntComparator{}, toGenericArray(), nil},
		{&IntComparator{}, toGenericArray(20, 10, 30), toGenericArray(10, 20, 30)},
		{&IntComparator{}, toGenericArray(20, 10, 30, 40, 50), toGenericArray(10, 20, 30, 40, 50)},
		{&IntComparator{}, toGenericArray(20, 10, 30, 40, 50, 80, 60, 70, 90), toGenericArray(10, 20, 30, 40, 50, 60, 70, 80, 90)},
	}

	for _, test := range tests {
		for k := 0; k < len(test.items); k++ {
			item := Select(test.items, k, test.cmp)

			assert.Equal(t, test.expectedItems[k], item)
		}
	}
}
