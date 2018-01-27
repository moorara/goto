package graph

// DFS implements Depth-First Search for Graph
type DFS struct {
	g      Graph
	s      int
	marked []bool
	edgeTo []int
}

// NewDFS creates a new DFS
func NewDFS(g Graph, s int) *DFS {
	dfs := &DFS{g: g, s: s}
	dfs.runDFS(s)
	return dfs
}

func (dfs *DFS) runDFS(v int) {
	dfs.marked[v] = true
	for _, w := range dfs.g.Adjacency(v) {
		if !dfs.marked[w] {
			dfs.runDFS(w)
			dfs.edgeTo[w] = v
		}
	}
}
