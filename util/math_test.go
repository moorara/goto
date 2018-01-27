package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinInt(t *testing.T) {
	tests := []struct {
		nums        []int
		expectedMin int
	}{
		{[]int{}, minInt},
		{[]int{7}, 7},
		{[]int{10, 20}, 10},
		{[]int{40, 50, 20, 30, 10}, 10},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expectedMin, MinInt(tc.nums...))
	}
}

func TestMinFloat64(t *testing.T) {
	tests := []struct {
		nums        []float64
		expectedMin float64
	}{
		{[]float64{}, minFloat64},
		{[]float64{3.14}, 3.14},
		{[]float64{.10, .20}, .10},
		{[]float64{.40, .50, .20, .30, .10}, .10},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expectedMin, MinFloat64(tc.nums...))
	}
}

func TestMaxInt(t *testing.T) {
	tests := []struct {
		nums        []int
		expectedMax int
	}{
		{[]int{}, maxInt},
		{[]int{7}, 7},
		{[]int{10, 20}, 20},
		{[]int{40, 50, 20, 30, 10}, 50},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expectedMax, MaxInt(tc.nums...))
	}
}

func TestMaxFloat64(t *testing.T) {
	tests := []struct {
		nums        []float64
		expectedMax float64
	}{
		{[]float64{}, maxFloat64},
		{[]float64{3.14}, 3.14},
		{[]float64{.10, .20}, .20},
		{[]float64{.40, .50, .20, .30, .10}, .50},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expectedMax, MaxFloat64(tc.nums...))
	}
}

func TestIsIntIn(t *testing.T) {
	tests := []struct {
		num            int
		list           []int
		expectedResult bool
	}{
		{5, []int{}, false},
		{5, []int{5}, true},
		{5, []int{10, 20}, false},
		{10, []int{10, 20}, true},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expectedResult, IsIntIn(tc.num, tc.list...))
	}
}
