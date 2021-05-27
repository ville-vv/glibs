package vtask

import (
	"errors"
	"sync"
)

func newNode(val interface{}) *node {
	return &node{value: val}
}

func New() Queue {
	return newChainQueue()
}

type Queue interface {
	Pop() interface{}
	Push(v interface{}) error
	Length() int64
}

type node struct {
	value    interface{}
	previous *node
	Next     *node
}

type ChainQueue struct {
	lock     sync.Mutex
	head     *node
	rear     *node
	length   int64
	capacity int64
}

func newChainQueue(args ...interface{}) *ChainQueue {
	max := int64(100000000)
	if len(args) > 0 {
		switch args[0].(type) {
		case int64:
			max = args[0].(int64)
		case int:
			max = int64(args[0].(int))
		}
	}
	return &ChainQueue{
		capacity: max,
	}
}
func (sel *ChainQueue) Pop() interface{} {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if sel.length <= 0 {
		sel.length = 0
		return nil
	}
	val := sel.head
	sel.head = sel.head.Next
	sel.length--
	val.previous = nil
	val.Next = nil
	return val.value
}
func (sel *ChainQueue) Shift() interface{} {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if sel.length <= 0 {
		sel.length = 0
		return nil
	}
	val := sel.rear
	if sel.rear.previous == nil {
		sel.rear = sel.head
	} else {
		sel.rear = sel.rear.previous
		sel.rear.Next = nil
	}
	val.previous = nil
	val.Next = nil
	sel.length--
	return val.value
}
func (sel *ChainQueue) Push(n interface{}) error {
	if sel.length >= sel.capacity {
		return errors.New("over max num for stack")
	}
	sel.push(&node{value: n})
	return nil
}
func (sel *ChainQueue) push(top *node) {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if 0 == sel.length {
		sel.head = top
		sel.rear = sel.head
	}
	top.Next = sel.head
	sel.head.previous = top
	sel.head = top
	sel.length++
	return
}
func (sel *ChainQueue) Length() int64 {
	return sel.length
}

var (
	defaultQueue = newChainQueue()
)

func Pop() interface{} {
	return defaultQueue.Pop()
}
func Push(v interface{}) error {
	return defaultQueue.Push(v)
}
func Shift() interface{} {
	return defaultQueue.Shift()
}
func Length() int64 {
	return defaultQueue.Length()
}
