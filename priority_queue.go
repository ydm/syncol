package syncol

import "container/heap"

type items[T any] struct {
	ary []T
	cmp Comparator[T]
}

func (xs items[T]) Len() int { return len(xs.ary) }

func (xs items[T]) Less(i, j int) bool {
	return xs.cmp(xs.ary[i], xs.ary[j])
}

func (xs items[T]) Swap(i, j int) {
	xs.ary[i], xs.ary[j] = xs.ary[j], xs.ary[i]
}

func (xs *items[T]) Push(x any) {
	xs.ary = append(xs.ary, x.(T))
}

func (xs *items[T]) Pop() any {
	n := len(xs.ary)
	if n <= 0 {
		return nil
	}
	m := n - 1
	ans := xs.ary[m]
	// Avoid memory leaks. --+
	var zero T       //      |
	xs.ary[m] = zero //      |
	// ----------------------+
	xs.ary = xs.ary[:m]
	return ans
}

type priorityQueue[T any] struct {
	xs items[T]
}

func (q *priorityQueue[T]) Init() {
}

func (q *priorityQueue[T]) Put(item T) {
	heap.Push(&q.xs, item)
}

func (q *priorityQueue[T]) Get() (item T, ok bool) {
	n := len(q.xs.ary)
	if n <= 0 {
		return item, false
	}
	return heap.Pop(&q.xs).(T), true
}

func NewPriorityQueue[T any](cmp Comparator[T]) *SynchronizedCollection[T] {
	xs := items[T]{make([]T, 0, 16), cmp}
	pq := priorityQueue[T]{xs}
	return NewSynchronizedCollection[T](&pq)
}
