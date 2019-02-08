package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelect(t *testing.T) {
	tests := []struct {
		compare       func(a, b interface{}) int
		items         []interface{}
		expectedItems []interface{}
	}{
		{compareInt, []interface{}{}, nil},
		{compareInt, []interface{}{20, 10, 30}, []interface{}{10, 20, 30}},
		{compareInt, []interface{}{20, 10, 30, 40, 50}, []interface{}{10, 20, 30, 40, 50}},
		{compareInt, []interface{}{20, 10, 30, 40, 50, 80, 60, 70, 90}, []interface{}{10, 20, 30, 40, 50, 60, 70, 80, 90}},
	}

	for _, tc := range tests {
		for k := 0; k < len(tc.items); k++ {
			item := Select(tc.items, k, tc.compare)
			assert.Equal(t, tc.expectedItems[k], item)
		}
	}
}

func TestQuickSortInt(t *testing.T) {
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
		QuickSort(tc.items, tc.compare)
		assert.True(t, sorted(tc.items, tc.compare))
	}
}

func TestQuickSort3WayInt(t *testing.T) {
	tests := []struct {
		compare func(a, b interface{}) int
		items   []interface{}
	}{
		{compareInt, []interface{}{}},
		{compareInt, []interface{}{20, 10, 10, 20, 30, 30, 30}},
		{compareInt, []interface{}{30, 20, 30, 20, 10, 40, 40, 40, 50, 50}},
		{compareInt, []interface{}{90, 10, 80, 20, 70, 30, 60, 40, 50, 50, 40, 60, 30, 70, 20, 80, 10, 90}},
	}

	for _, tc := range tests {
		QuickSort3Way(tc.items, tc.compare)
		assert.True(t, sorted(tc.items, tc.compare))
	}
}

func TestQuickSortString(t *testing.T) {
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
		QuickSort(tc.items, tc.compare)
		assert.True(t, sorted(tc.items, tc.compare))
	}
}

func TestQuickSort3WayString(t *testing.T) {
	tests := []struct {
		compare func(a, b interface{}) int
		items   []interface{}
	}{
		{compareString, []interface{}{}},
		{compareString, []interface{}{"Milad", "Mona", "Milad", "Mona"}},
		{compareString, []interface{}{"Alice", "Bob", "Alex", "Jackie", "Jackie", "Alex", "Bob", "Alice"}},
		{compareString, []interface{}{"Docker", "Kubernetes", "Docker", "Go", "JavaScript", "Go", "React", "Redux", "Vue", "Redux", "React"}},
	}

	for _, tc := range tests {
		QuickSort3Way(tc.items, tc.compare)
		assert.True(t, sorted(tc.items, tc.compare))
	}
}
