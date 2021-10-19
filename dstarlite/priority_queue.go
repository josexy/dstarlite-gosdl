package dstarlite

import "container/heap"

type NodeHeap []*Node

func (list *NodeHeap) Clear() {
	*list = (*list)[:0]
}

func (list NodeHeap) Len() int {
	return len(list)
}

func (list NodeHeap) Less(i, j int) bool {
	n1, n2 := list[i], list[j]
	if n1.k.k1 == n2.k.k1 {
		return n1.k.k2 < n2.k.k2
	}
	return n1.k.k1 < n2.k.k1
}

func (list NodeHeap) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list *NodeHeap) Push(x interface{}) {
	*list = append(*list, x.(*Node))
}

func (list *NodeHeap) Pop() interface{} {
	n := list.Len()
	value := (*list)[n-1]
	*list = (*list)[:n-1]
	return value
}

func (list *NodeHeap) Top() interface{} {
	if list.Len() == 0 {
		return nil
	}
	return (*list)[0]
}

type PriorityQueue struct {
	list *NodeHeap
}

func NewPriorityQueue() *PriorityQueue {
	return &PriorityQueue{list: &NodeHeap{}}
}

func (p *PriorityQueue) Clear() {
	p.list.Clear()
}

func (p *PriorityQueue) Size() int {
	return p.list.Len()
}

func (p *PriorityQueue) IsEmpty() bool {
	return p.list.Len() == 0
}

func (p *PriorityQueue) Push(node *Node) {
	heap.Push(p.list, node)
}

func (p *PriorityQueue) Pop() *Node {
	return heap.Pop(p.list).(*Node)
}

func (p *PriorityQueue) Top() *Node {
	if p.IsEmpty() {
		return nil
	}
	return p.list.Top().(*Node)
}

func (p *PriorityQueue) Find(node *Node) int {
	for i, s := range *p.list {
		if s.Equal(node) {
			return i
		}
	}
	return -1
}

func (p *PriorityQueue) Remove(i int) {
	heap.Remove(p.list, i)
}
