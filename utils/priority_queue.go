package utils

import (
	"cmp"
	"container/heap"
)

// An Item is something we manage in a priority queue.
type Item[T any, U cmp.Ordered] struct {
	Value    T // The value of the item; arbitrary.
	Priority U // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

type PriorityQueue[T any, U cmp.Ordered] struct {
	items priorityQueue[T, U]
}

func NewPriorityQueue[T any, U cmp.Ordered]() *PriorityQueue[T, U] {
	return &PriorityQueue[T, U]{
		items: make(priorityQueue[T, U], 0),
	}
}

func (pq *PriorityQueue[T, U]) Push(value T, prio U) {
	heap.Push(&pq.items, &Item[T, U]{Value: value, Priority: prio})
}

func (pq *PriorityQueue[T, U]) Pop() *Item[T, U] {
	return heap.Pop(&pq.items).(*Item[T, U])
}

func (pq PriorityQueue[T, U]) IsEmpty() bool {
	return pq.items.Len() == 0
}

type priorityQueue[T any, U cmp.Ordered] []*Item[T, U]

func (pq priorityQueue[T, U]) Len() int { return len(pq) }

func (pq priorityQueue[T, U]) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq priorityQueue[T, U]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue[T, U]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[T, U])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue[T, U]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *priorityQueue[T, U]) update(item *Item[T, U], value T, priority U) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.index)
}
