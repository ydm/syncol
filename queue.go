package syncol

import "container/list"

type queue struct {
	q *list.List
}

func newQueue() *queue {
	return &queue{list.New()}
}

func (q *queue) Init() {
	q.q = list.New()
}

func (q *queue) Put(item interface{}) {
	q.q.PushBack(item)
}

func (q *queue) Get() (item interface{}, ok bool) {
	e := q.q.Front()
	if e != nil {
		item = q.q.Remove(e)
		ok = true
	}
	return
}

func NewQueue() *SynchronizedCollection {
	return NewSynchronizedCollection(newQueue())
}
