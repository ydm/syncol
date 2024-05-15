package syncol

import "sync"

type SynchronizedCollection[T any] struct {
	c               Collection[T]
	mu              sync.Mutex
	signal          *sync.Cond
	unfinishedTasks int
}

func NewSynchronizedCollection[T any](col Collection[T]) *SynchronizedCollection[T] {
	col.Init()
	ans := SynchronizedCollection[T]{
		c:               col,
		mu:              sync.Mutex{},
		signal:          nil,
		unfinishedTasks: 0,
	}
	ans.signal = sync.NewCond(&ans.mu)

	return &ans
}

func (q *SynchronizedCollection[T]) TaskDone() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.unfinishedTasks--

	if q.unfinishedTasks <= 0 {
		if q.unfinishedTasks < 0 {
			panic("TaskDone() called too many times")
		}

		q.signal.Broadcast()
	}
}

func (q *SynchronizedCollection[T]) Join() {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.unfinishedTasks > 0 {
		q.signal.Wait()
	}
}

func (q *SynchronizedCollection[T]) Put(item T) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.c.Put(item)
	q.unfinishedTasks++ //nolint:wsl
	q.signal.Broadcast()
}

// Get returns when there is a new item Put() in the queue (with ok
// set to true) or when there is nothing more to get (ok = false).
func (q *SynchronizedCollection[T]) Get() (item T, ok bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.unfinishedTasks > 0 {
		item, ok = q.c.Get()
		if ok {
			return item, ok
		}
		// Wait either for an element to be present or for a
		// task to be done.
		q.signal.Wait()
	}

	return item, false
}
