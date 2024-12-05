package utils

type ID any

type Graph[ID comparable] struct {
	Nodes  []ID
	Lookup map[ID]int
	Edges  []map[int]struct{}
}

func NewGraph[ID comparable]() *Graph[ID] {
	return &Graph[ID]{Lookup: make(map[ID]int)}
}

func (g *Graph[ID]) AddNode(id ID) {
	Assert(!g.HasNode(id))

	g.addNode(id)
}

func (g *Graph[ID]) addNode(id ID) {
	g.Lookup[id] = len(g.Nodes)
	g.Nodes = append(g.Nodes, id)
	g.Edges = append(g.Edges, make(map[int]struct{}))
}

func (g *Graph[ID]) AddEdge(a, b ID) {
	if !g.HasNode(a) {
		g.addNode(a)
	}
	if !g.HasNode(b) {
		g.addNode(b)
	}
	Assert(!g.HasEdge(a, b))
	g.addEdge(a, b)
}

func (g *Graph[ID]) addEdge(a, b ID) {
	g.Edges[g.Lookup[a]][g.Lookup[b]] = struct{}{}
}

func (g *Graph[ID]) HasNode(id ID) bool {
	_, ok := g.Lookup[id]
	return ok
}

func (g *Graph[ID]) HasEdge(a, b ID) bool {
	if !g.HasNode(a) {
		return false
	}
	if !g.HasNode(b) {
		return false
	}
	_, ok := g.Edges[g.Lookup[a]][g.Lookup[b]]
	return ok
}
