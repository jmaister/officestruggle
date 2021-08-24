package astar

import (
	"container/heap"
	"fmt"
)

type NodeFScore struct {
	Node   *Node
	FScore int
	Index  int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*NodeFScore

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].FScore > pq[j].FScore
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(NodeFScore)
	item.Index = n
	*pq = append(*pq, &item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *NodeFScore, node *Node, priority int) {
	item.Node = node
	item.FScore = priority
	heap.Fix(pq, item.Index)
	fmt.Println("pq", pq)
}
