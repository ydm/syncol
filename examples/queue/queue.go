package main

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ydm/syncol"
)

// +------+
// | node |
// +------+

type node struct {
	Value int
	left  *node
	right *node
}

func (n *node) Add(x int) {
	sub := &n.left
	if n.Value < x {
		sub = &n.right
	}
	if *sub != nil {
		(*sub).Add(x)
	} else {
		*sub = &node{x, nil, nil}
	}
}

func (n *node) String() string {
	left := "_"
	if n.left != nil {
		left = n.left.String()
	}
	right := "_"
	if n.right != nil {
		right = n.right.String()
	}
	return fmt.Sprintf("{%d %v %v}", n.Value, left, right)
}

// +-----------+
// | nodeStack |
// +-----------+

type nodeStack struct {
	stack *list.List
	mu    sync.Mutex
}

func newNodeStack(initial *node) *nodeStack {
	stack := list.New()
	stack.PushBack(initial)
	return &nodeStack{stack: stack}
}

func (s *nodeStack) Push(n *node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stack.PushBack(n)
}

func (s *nodeStack) Pop() (n *node, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	back := s.stack.Back()
	if back != nil {
		n = s.stack.Remove(back).(*node)
		ok = true
	}
	return
}

// +---------+
// | bounded |
// +---------+

var invocations int64

func bounded(root *node) {
	queue := syncol.NewQueue()
	queue.Put(root)

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go boundedWorker(queue)
	}
	queue.Join()
}

func boundedWorker(queue *syncol.SynchronizedCollection) {
	id := atomic.AddInt64(&invocations, 1)
	for {
		time.Sleep(100 * time.Millisecond)
		item, ok := queue.Get()
		if !ok {
			break
		}
		// Do something with item
		n := item.(*node)
		fmt.Printf("[%d] %d\n", id, n.Value)
		// Enqueue more items.
		if n.left != nil {
			queue.Put(n.left)
		}
		if n.right != nil {
			queue.Put(n.right)
		}
		queue.TaskDone()
	}
}

// +-----------+
// | unbounded |
// +-----------+

func unbounded(root *node) {
	var wg sync.WaitGroup
	wg.Add(1)
	go unboundedWorker(root, &wg)
	wg.Wait()
}

func unboundedWorker(n *node, wg *sync.WaitGroup) {
	defer wg.Done()
	id := atomic.AddInt64(&invocations, 1)
	for {
		fmt.Printf("[%d] %d\n", id, n.Value)
		if n.right != nil {
			if n.left != nil {
				wg.Add(1)
				go unboundedWorker(n.left, wg)
			}
			n = n.right
		} else if n.left != nil {
			n = n.left
		} else {
			// This is a leaf.
			break
		}
	}
}

func main() {
	// Make and populate a binary tree.
	root := &node{11, nil, nil}
	values := []int{10, 15, 7, 2, 16, 4, 1, 6, 3, 8, 12, 14, 13, 5, 9}
	for _, x := range values {
		root.Add(x)
	}

	fmt.Println("UNBOUNDED")
	invocations = 0
	unbounded(root)
	fmt.Printf("invocations=%d\n\n", invocations)

	fmt.Println("BOUNDED")
	invocations = 0
	bounded(root)
	fmt.Printf("invocations=%d\n\n", invocations)

	time.Sleep(1 * time.Second)
	fmt.Println("PUT A BREAKPOINT HERE")
}
