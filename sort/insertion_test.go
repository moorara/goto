package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertionSortInt(t *testing.T) {
	tests := []struct {
		compare func(a, b interface{}) int
		items   []interface{}
	}{
		{compareInt, []interface{}{}},
		{compareInt, []interface{}{20, 10, 30}},
		{compareInt, []interface{}{30, 20, 10, 40, 50}},
		{compareInt, []interface{}{90, 80, 70, 60, 50, 40, 30, 20, 10}},
	}

	for _, tc := range tests {
		InsertionSort(tc.items, tc.compare)
		assert.True(t, sorted(tc.items, tc.compare))
	}
}

func TestInsertionSortString(t *testing.T) {
	tests := []struct {
		compare func(a, b interface{}) int
		items   []interface{}
	}{
		{compareString, []interface{}{}},
		{compareString, []interface{}{"Milad", "Mona"}},
		{compareString, []interface{}{"Alice", "Bob", "Alex", "Jackie"}},
		{compareString, []interface{}{"Docker", "Kubernetes", "Go", "JavaScript", "Elixir", "React", "Redux", "Vue"}},
	}

	for _, tc := range tests {
		InsertionSort(tc.items, tc.compare)
		assert.True(t, sorted(tc.items, tc.compare))
	}
}
