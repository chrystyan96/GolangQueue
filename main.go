package main

import (
	"fmt"
	"os"
)

//FIFO -> First In, First Out

type Item = int
type deqResult struct {
	item Item
	err  error
}

type Queue struct {
	items []Item
	enq   chan Item
	deq   chan struct{}
	ret   chan deqResult
	done  chan struct{}
}

func NewQueue() *Queue {
	q := new(Queue)
	q.items = []Item{}
	q.enq = make(chan Item)
	q.deq = make(chan struct{})
	q.ret = make(chan deqResult)
	q.done = make(chan struct{})
	return q
}

func (q *Queue) add(item Item) {
	q.items = append(q.items, item)
}

func (q *Queue) Run() {
	for {
		select {
		case item := <-q.enq:
			q.add(item)
			fmt.Printf("Enqueued: %d\n", item)
		case <-q.deq:
			res := deqResult{}
			if len(q.items) == 0 {
				res.item = 0
				res.err = fmt.Errorf("no items in the queue")
				q.ret <- res
				continue
			}

			res.item = q.items[0]
			res.err = nil
			q.items = q.items[1:]
			q.ret <- res
		case <-q.done:
			return
		}
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0 && len(q.deq) == 0
}

// Enqueue
// |			   F ADD
// |			   v
// [ A, B, C, D, E ]
func (q *Queue) Enqueue(item Item) {
	q.enq <- item
}

// Dequeue
// |A POP
// |^
// [ B, C, D, E, F ]
func (q *Queue) Dequeue() (Item, error) {
	q.deq <- struct{}{}
	res := <-q.ret
	if res.err != nil {
		return 0, fmt.Errorf("dequeue() error: %s", res.err)
	}

	return res.item, nil
}

func (q *Queue) Kill() {
	q.done <- struct{}{}
}

func main() {
	q := NewQueue()
	go q.Run()

	q.Enqueue(1) // -> q.enq
	q.Enqueue(2) // -> q.enq
	q.Enqueue(3) // -> q.enq

	for !q.IsEmpty() {
		item, err := q.Dequeue() // -> q.deq
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Dequeued:", item)
	}

	q.Kill()
}
