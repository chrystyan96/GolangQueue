package main

import (
	"errors"
	"fmt"
)

//FIFO -> First In, First Out

type Item = int

type Queue struct {
	items []Item
	send  chan Item
	done  chan struct{}
}

func NewQueue() *Queue {
	q := new(Queue)
	q.items = []Item{}
	q.send = make(chan Item, 100)
	return q
}

func (q *Queue) add(item Item) {
	q.items = append(q.items, item)
}

func (q *Queue) Run() {
	for {
		select {
		case item := <-q.send:
			q.add(item)
			fmt.Printf("Got item from channel: %d/n", item)
		case <-q.done:
			return
		}
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Enqueue
// |			   F ADD
// |			   v
// [ A, B, C, D, E ]
func (q *Queue) Enqueue(item Item) {
	q.send <- item
}

// Dequeue
// |A POP
// |^
// [ B, C, D, E, F ]
func (q *Queue) Dequeue() (Item, error) {
	if len(q.items) == 0 {
		return 0, errors.New("queue is empty")
	}

	it := q.items[0]
	q.items = q.items[1:]
	return it, nil
}

func main() {
	q := NewQueue()
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	for {
		item, err := q.Dequeue()
		if err != nil {
			break
		}

		fmt.Println(item)
	}
}
