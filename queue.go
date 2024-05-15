package syncol

import "container/list"

type queue[T any] struct {
	q *list.List
}

func newQueue[T any]() *queue[T] {
	return &queue[T]{list.New()}
}

func (q *queue[T]) Init() {
	q.q = list.New()
}

func (q *queue[T]) Put(item T) {
	q.q.PushBack(item)
}

func (q *queue[T]) Get() (item T, ok bool) {
	e := q.q.Front()
	if e != nil {
		item = q.q.Remove(e).(T)
		ok = true
	}
	return
}

func NewQueue[T any]() *SynchronizedCollection[T] {
	var col Collection[T] = newQueue[T]()
	return NewSynchronizedCollection(col)
}
