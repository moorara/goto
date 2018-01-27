package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type graphTest struct {
	name     string
	vertices int
	edges    [][2]int
}

func getGraphTests() []graphTest {
	return []graphTest{
		{
			name:     "Empty",
			vertices: 0,
			edges:    [][2]int{},
		},
		{
			name:     "SimpleGraph",
			vertices: 4,
			edges: [][2]int{
				{0, 1}, {0, 2},
				{1, 3},
				{2, 3},
			},
		},
		{
			name:     "ComplexGraph",
			vertices: 10,
			edges: [][2]int{
				{0, 2}, {0, 4}, {0, 8},
				{1, 2}, {1, 3}, {1, 4}, {1, 5},
				{2, 4}, {2, 8},
				{3, 6}, {3, 9},
				{4, 8},
				{5, 6}, {5, 7}, {5, 8}, {5, 9},
				{6, 7}, {6, 8}, {6, 9},
				{7, 8}, {7, 9},
				{8, 9},
			},
		},
	}
}

func TestGraph(t *testing.T) {
	graphTests := getGraphTests()

	tests := []struct {
		graph             graphTest
		expectedAdjacency [][]int
		expectedAdjacents [][]int
		expectedDegrees   []int
		expectedMaxDegree int
		expectedAvgDegree float64
	}{
		{
			graph:             graphTests[0],
			expectedAdjacency: [][]int{},
			expectedAdjacents: [][]int{},
			expectedDegrees:   []int{},
			expectedMaxDegree: 0,
			expectedAvgDegree: 0,
		},
		{

			graph: graphTests[1],
			expectedAdjacency: [][]int{
				{1, 2},
				{0, 3},
				{0, 3},
				{1, 2},
			},
			expectedAdjacents: [][]int{
				{0, 1}, {0, 2},
				{1, 0}, {1, 3},
				{2, 0}, {2, 3},
				{3, 1}, {3, 2},
			},
			expectedDegrees:   []int{2, 2, 2, 2},
			expectedMaxDegree: 2,
			expectedAvgDegree: 2,
		},
		{

			graph: graphTests[2],
			expectedAdjacency: [][]int{
				{2, 4, 8},
				{2, 3, 4, 5},
				{0, 1, 4, 8},
				{1, 6, 9},
				{0, 1, 2, 8},
				{1, 6, 7, 8, 9},
				{3, 5, 7, 8, 9},
				{5, 6, 8, 9},
				{0, 2, 4, 5, 6, 7, 9},
				{3, 5, 6, 7, 8},
			},
			expectedAdjacents: [][]int{
				{0, 2}, {0, 4}, {0, 8},
				{1, 2}, {1, 3}, {1, 4}, {1, 5},
				{2, 0}, {2, 1}, {2, 4}, {2, 8},
				{3, 1}, {3, 6}, {3, 9},
				{4, 0}, {4, 1}, {4, 2}, {4, 8},
				{5, 1}, {5, 6}, {5, 7}, {5, 8}, {5, 9},
				{6, 3}, {6, 5}, {6, 7}, {6, 8}, {6, 9},
				{7, 5}, {7, 6}, {7, 8}, {7, 9},
				{8, 0}, {8, 2}, {8, 4}, {8, 5}, {8, 6}, {8, 7}, {8, 9},
				{9, 3}, {9, 5}, {9, 6}, {9, 7}, {9, 8},
			},
			expectedDegrees:   []int{3, 4, 4, 3, 4, 5, 5, 4, 7, 5},
			expectedMaxDegree: 7,
			expectedAvgDegree: 4.4,
		},
	}

	for _, tc := range tests {
		t.Run(tc.graph.name, func(t *testing.T) {
			g := NewGraph(tc.graph.vertices)
			for _, e := range tc.graph.edges {
				g.AddEdge(e[0], e[1])
			}

			assert.NotNil(t, g)
			assert.Equal(t, tc.graph.vertices, g.V())
			assert.Equal(t, len(tc.graph.edges), g.E())

			for v, expectedAdj := range tc.expectedAdjacency {
				assert.Equal(t, expectedAdj, g.Adjacency(v))
			}

			for _, e := range tc.expectedAdjacents {
				assert.True(t, g.Adjacent(e[0], e[1]))
			}

			for v, expectedDegree := range tc.expectedDegrees {
				assert.Equal(t, expectedDegree, g.Degree(v))
			}

			assert.Equal(t, tc.expectedMaxDegree, g.MaxDegree())
			assert.Equal(t, tc.expectedAvgDegree, g.AvgDegree())
		})
	}
}
