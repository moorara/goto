package graph

/* import (
	. "github.com/moorara/goto/dt"
) */

type (
	// Graph represents an undirected graph data type
	Graph interface {
		V() int
		E() int
		// AddVertex(int)
		// RemoveVertex(int)
		AddEdge(int, int)
		// RemoveEdge(int, int)
		Adjacency(int) []int
		Adjacent(int, int) bool
		Degree(int) int
		MaxDegree() int
		AvgDegree() float64
		// SetVertexValue(int, Generic)
		// GetVertexValue(int, Generic)
		// SetEdgeValue(int, int, Generic)
		// GetEdgeValue(int, int, Generic)
		// String() string
		// Graphviz() string
	}

	graph struct {
		e   int
		adj [][]int
	}
)

// NewGraph creates a new undirected graph
func NewGraph(v int) Graph {
	adj := make([][]int, v)
	for v := range adj {
		adj[v] = make([]int, 0)
	}

	return &graph{
		adj: adj,
	}
}

func (g *graph) V() int {
	return len(g.adj)
}

func (g *graph) E() int {
	return g.e
}

func (g *graph) AddEdge(v, w int) {
	g.e++
	g.adj[v] = append(g.adj[v], w)
	g.adj[w] = append(g.adj[w], v)
}

func (g *graph) Adjacency(v int) []int {
	return g.adj[v]
}

func (g *graph) Adjacent(v, w int) bool {
	adj := g.Adjacency(v)
	for _, x := range adj {
		if x == w {
			return true
		}
	}
	return false
}

func (g *graph) Degree(v int) int {
	return len(g.adj[v])
}

func (g *graph) MaxDegree() int {
	max := 0
	for v := 0; v < g.V(); v++ {
		d := g.Degree(v)
		if d > max {
			max = d
		}
	}
	return max
}

func (g *graph) AvgDegree() float64 {
	if g.V() == 0 {
		return 0
	}
	return 2.0 * float64(g.E()) / float64(g.V())
}
