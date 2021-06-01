package vtask

import (
	"fmt"
	"testing"
	"time"
)

type asb struct {
	ad int
}

func TestRingQueue_Pop(t *testing.T) {
	rQ := NewRingQueue(10000)

	l := AtomicInt64{}

	add := AtomicInt64{}

	for i := 0; i < 20; i++ {
		go func(n int) {
			for {
				if rQ.Length() == 0 {
					continue
				}
				val := rQ.Pop()
				fmt.Println(val)
				if val != nil {

					l.Inc()
				}
				time.Sleep(time.Millisecond * 10)
			}
		}(i)
	}

	for i := 1; i <= 20; i++ {
		go func(n int) {
			for i := 0; i < 10; i++ {
				err := rQ.Push(&asb{ad: i * n})
				if err == nil {
					add.Inc()
				}
				time.Sleep(time.Millisecond * 10)

			}

		}(i)
	}

	time.Sleep(time.Second * 5)
	fmt.Println(l.Load(), add.Load())

}
