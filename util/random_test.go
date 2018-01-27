package util

import (
	"testing"

	. "github.com/moorara/goto/dt"
	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	tests := []struct {
		items []Generic
	}{
		{[]Generic{10, 20, 30, 40, 50, 60, 70, 80, 90}},
		{[]Generic{"Alice", "Bob", "Dan", "Edgar", "Helen", "Karen", "Milad", "Peter", "Sam", "Wesley"}},
	}

	SeedWithNow()

	for _, tc := range tests {
		orig := make([]Generic, len(tc.items))
		copy(orig, tc.items)
		Shuffle(tc.items)

		assert.NotEqual(t, orig, tc.items)
	}
}

func TestGenerateInt(t *testing.T) {
	tests := []struct {
		min int
		max int
	}{
		{0, 0},
		{1, 1},
		{0, 1000},
		{100, 100000},
	}

	SeedWithNow()

	for _, tc := range tests {
		n := GenerateInt(tc.min, tc.max)

		assert.True(t, tc.min <= n && n <= tc.max)
	}
}

func TestGenerateString(t *testing.T) {
	tests := []struct {
		minLen int
		maxLen int
	}{
		{0, 0},
		{1, 1},
		{10, 100},
		{100, 1000},
	}

	SeedWithNow()

	for _, tc := range tests {
		str := GenerateString(tc.minLen, tc.maxLen)

		assert.True(t, tc.minLen <= len(str) && len(str) <= tc.maxLen)
	}
}

func TestGenerateIntSlice(t *testing.T) {
	tests := []struct {
		size int
		min  int
		max  int
	}{
		{0, 0, 0},
		{1, 1, 1},
		{10, 0, 100},
		{100, 100, 1000},
	}

	SeedWithNow()

	for _, tc := range tests {
		items := GenerateIntSlice(tc.size, tc.min, tc.max)
		for _, item := range items {
			if CompareInt(item, tc.min) < 0 || CompareInt(item, tc.max) > 0 {
				t.Errorf("%d is not between %d and %d.", item, tc.min, tc.max)
			}
		}
	}
}

func TestGenerateStringSlice(t *testing.T) {
	tests := []struct {
		size   int
		minLen int
		maxLen int
	}{
		{0, 0, 0},
		{1, 1, 1},
		{10, 1, 10},
		{100, 10, 100},
	}

	SeedWithNow()

	for _, tc := range tests {
		items := GenerateStringSlice(tc.size, tc.minLen, tc.maxLen)
		for _, item := range items {
			if len(item.(string)) < tc.minLen || len(item.(string)) > tc.maxLen {
				t.Errorf("%s length is not between %d and %d.", item, tc.minLen, tc.maxLen)
			}
		}
	}
}
