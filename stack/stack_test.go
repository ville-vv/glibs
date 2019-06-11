package stack

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type Do func(a, b int) int

func sub(a, b int) int {
	return a - b
}

func sum(a, b int) int {
	return a + b
}

func TestPool_Push(t *testing.T) {
	po := NewPool()
	num := 100000000
	start := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		po.Push(i)
	}
	end := time.Now().UnixNano()
	fmt.Println("时间1：", (end-start)/1e6)

	fmt.Println("长度：", po.Length())

	start = time.Now().UnixNano()
	for i := 0; i < num; i++ {
		po.Pop()
	}
	end = time.Now().UnixNano()
	fmt.Println("时间2：", (end-start)/1e6)
}

func TestPool_Push2(t *testing.T) {
	po := NewPool()
	num := 20
	start := time.Now().UnixNano()
	for i := int64(1); i < int64(num); i++ {
		po.Push(strconv.Itoa(int(i)))
	}
	end := time.Now().UnixNano()

	fmt.Println("时间1：", (end-start)/1e6)
	start = time.Now().UnixNano()
	for i := 1; i < 10; i++ {
		fmt.Println(po.Pop().ToInt())
	}
	end = time.Now().UnixNano()
	fmt.Println("时间3：", (end-start)/1e6)

	start = time.Now().UnixNano()
	for i := 1; i < 5; i++ {
		fmt.Println(po.Shift().ToInt64())
	}
	end = time.Now().UnixNano()
	fmt.Println("时间4：", (end-start)/1e6)
}

func TestStack_PushChan(t *testing.T) {
	po := NewPoolChan()
	num := 100000000
	start := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		po.PushChan(i)
	}
	end := time.Now().UnixNano()
	fmt.Println("时间1：", (end-start)/1e6)

	time.Sleep(time.Second * 15)
	fmt.Println("长度：", po.Length())
	start = time.Now().UnixNano()
	for i := 0; i < num; i++ {
		po.Pop()
	}
	end = time.Now().UnixNano()
	fmt.Println("时间2：", (end-start)/1e6)
	time.Sleep(time.Second * 15)
}
