package syncol

import "container/heap"

type items struct {
	ary []interface{}
	cmp Comparator
}

func (xs items) Len() int { return len(xs.ary) }

func (xs items) Less(i, j int) bool {
	return xs.cmp(xs.ary[i], xs.ary[j])
}

func (xs items) Swap(i, j int) {
	xs.ary[i], xs.ary[j] = xs.ary[j], xs.ary[i]
}

func (xs *items) Push(x interface{}) {
	xs.ary = append(xs.ary, x)
}

func (xs *items) Pop() interface{} {
	n := len(xs.ary)
	if n <= 0 {
		return nil
	}
	m := n-1
	ans := xs.ary[m]
	xs.ary[m] = nil // Avoid memory leak.
	xs.ary = xs.ary[:m]
	return ans
}

type priorityQueue struct {
	xs items
}

func (q *priorityQueue) Init() {
}

func (q *priorityQueue) Put(item interface{}) {
	heap.Push(&q.xs, item)
}

func (q *priorityQueue) Get() (item interface{}, ok bool) {
	n := len(q.xs.ary)
	if n <= 0 {
		return nil, false
	}
	return heap.Pop(&q.xs), true
}

func NewPriorityQueue(cmp Comparator) *SynchronizedCollection {
	xs := items{make([]interface{}, 0, 16), cmp}
	pq := priorityQueue{xs}
	return NewSynchronizedCollection(&pq)
}
