package vtask

import (
	"errors"
	"sync"
)

func newNode(val interface{}) *Node {
	return &Node{Value: val}
}

func New() Queue {
	return NewList()
}

type Queue interface {
	Pop() interface{}
	Push(v interface{}) error
	Length() int64
}

type Node struct {
	Value interface{}
	Front *Node
	Next  *Node
}

func NewNode(value interface{}) *Node {
	return &Node{Value: value}
}

type List struct {
	lock     sync.Mutex
	head     *Node
	rear     *Node
	length   int64
	capacity int64
}

func NewList(args ...interface{}) *List {
	max := int64(100000000)
	if len(args) > 0 {
		switch args[0].(type) {
		case int64:
			max = args[0].(int64)
		case int:
			max = int64(args[0].(int))
		}
	}
	return &List{
		capacity: max,
	}
}
func (sel *List) Pop() interface{} {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if sel.length <= 0 {
		sel.length = 0
		return nil
	}
	val := sel.head
	sel.head = sel.head.Next
	sel.length--
	val.Front = nil
	val.Next = nil
	return val.Value
}
func (sel *List) Shift() interface{} {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if sel.length <= 0 {
		sel.length = 0
		return nil
	}
	val := sel.rear
	if sel.rear.Front == nil {
		sel.rear = sel.head
	} else {
		sel.rear = sel.rear.Front
		sel.rear.Next = nil
	}
	val.Front = nil
	val.Next = nil
	sel.length--
	return val.Value
}
func (sel *List) Push(n interface{}) error {
	if sel.length >= sel.capacity {
		return errors.New("over max num for stack")
	}
	sel.push(&Node{Value: n})
	return nil
}
func (sel *List) push(top *Node) {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	if 0 == sel.length {
		sel.head = top
		sel.rear = sel.head
	}
	top.Next = sel.head
	sel.head.Front = top
	sel.head = top
	sel.length++
	return
}
func (sel *List) Length() int64 {
	return sel.length
}

var (
	defaultQueue = NewList()
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
