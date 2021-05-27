package vtask

import (
	"fmt"
	"testing"
)

type asb struct {
	ad int
}

func TestRingQueue_Pop(t *testing.T) {
	rQ := NewRingQueue(1000)
	for i := 0; i < 100; i++ {
		_ = rQ.Push(&asb{ad: i})
	}

	for {
		if rQ.Length() == 0 {
			break
		}
		fmt.Println(rQ.Pop())
	}
}
