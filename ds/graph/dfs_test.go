package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDFS(t *testing.T) {
	graphTests := getGraphTests()

	tests := []struct {
		graph graphTest
	}{
		{
			graph: graphTests[0],
		},
		{
			graph: graphTests[1],
		},
		{
			graph: graphTests[2],
		},
	}

	for _, tc := range tests {
		t.Run(tc.graph.name, func(t *testing.T) {
			g := NewGraph(tc.graph.vertices)
			for _, e := range tc.graph.edges {
				g.AddEdge(e[0], e[1])
			}

			assert.NotNil(t, g)
		})
	}
}
