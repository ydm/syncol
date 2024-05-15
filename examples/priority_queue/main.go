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

func compareItems(lhs, rhs *Item) bool {
	return rhs.Less(*lhs)
}

func main() {
	q := syncol.NewPriorityQueue(compareItems)
	q.Put(&Item{"нихилизъм", 2})
	q.Put(&Item{"пънк", 3})
	q.Put(&Item{"анархия", 1})
	for {
		item, ok := q.Get()
		if !ok {
			break
		}
		fmt.Printf("item=%v\n", item)
		q.TaskDone()
	}
	q.Join()
	fmt.Println("END")
}
