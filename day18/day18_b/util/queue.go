// VV: Adapted from https://pkg.go.dev/container/heap
// Heap which pops element with minimum priority
package util

import (
	"container/heap"
)

// An HeapItem is something we manage in a priority queue.
type HeapItem struct {
	Value    any // The value of the item; arbitrary.
	Priority int // The priority of the item in the queue.
	index    int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*HeapItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*HeapItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *HeapItem, value any, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.index)
}

func (pq *PriorityQueue) Upsert(value any, priority int) {
	updated := false

	for idx, q := range *pq {
		if q.Value == value {
			updated = true

			q.Priority = priority
			heap.Fix(pq, idx)
			break
		}
	}

	if !updated {
		t := &HeapItem{
			Value:    value,
			Priority: priority,
		}

		heap.Push(pq, t)
	}
}
