package utils

import (
	"cmp"
	"slices"
)

type WithCost[T any, U cmp.Ordered] struct {
	Value T
	Cost  U
}

type Visit[T any] struct {
	Value  T
	Parent T
}

func Dijkstra[T comparable, U cmp.Ordered](start WithCost[T, U], isDone func(T) bool, neighbors func(WithCost[T, U]) []WithCost[T, U]) ([]T, U, bool) {

	pq := NewPriorityQueue[T, U]()
	pq.Push(start.Value, start.Cost)
	visited := NewSet[T]()
	parent := make(map[T]T, 0)
	cost := make(map[T]U, 0)

	for !pq.IsEmpty() {
		current := pq.Pop()
		if visited.Exists(current.Value) {
			continue
		}

		visited.Add(current.Value)

		for _, n := range neighbors(WithCost[T, U]{current.Value, current.Priority}) {
			if visited.Exists(n.Value) {
				continue
			}

			if isDone(n.Value) {
				path := make([]T, 0)
				path = append(path, n.Value)
				node := current.Value
				ok := true
				for ok {
					path = append(path, node)
					node, ok = parent[node]
				}

				slices.Reverse(path)
				return path, n.Cost, true
			}

			if c, ok := cost[n.Value]; !ok || n.Cost < c {
				parent[n.Value] = current.Value
				cost[n.Value] = n.Cost
				pq.Push(n.Value, n.Cost)
			}
		}
	}
	return nil, start.Cost, false
}
