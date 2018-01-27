package graph

// BFS implements Breadth-First Search for Graph
type BFS struct {
	g Graph
	s int
}

// NewBFS creates a new BFS
func NewBFS(g Graph, s int) *BFS {
	bfs := &BFS{g: g, s: s}
	bfs.runBFS(s)
	return bfs
}

func (bfs *BFS) runBFS(v int) {
	// TODO
}
