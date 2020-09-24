package main

import (
	"fmt"
	"github.com/ydm/syncol"
)

type Item struct {
	name string
	priority int
}

func (i Item) Less(other Item) bool {
	return i.priority < other.priority
}

func compareItems(lhs, rhs interface{}) bool {
	left := lhs.(*Item)
	right := rhs.(*Item)
	return right.Less(*left)
}

func main() {
	q := syncol.NewPriorityQueue(compareItems)
	q.Put(&Item{"цици", 2})
	q.Put(&Item{"пънк", 3})
	q.Put(&Item{"наркотици", 1})
	for {
		p, ok := q.Get()
		if !ok {
			break
		}
		item := p.(*Item)
		fmt.Printf("item=%v\n", item)
		q.TaskDone()
	}
	q.Join()
	fmt.Println("END")
}
