package vtask

import (
	"fmt"
	"testing"
	"time"
)

func TestChan(t *testing.T) {

	testCh := make(chan int, 100)
	for i := 0; i < 100; i++ {
		testCh <- i
	}
	close(testCh)
	go func() {
		time.Sleep(time.Millisecond * 2)
		for {
			select {
			case val, ok := <-testCh:
				if !ok {
					return
				}
				fmt.Println(val)
			}
		}
	}()
	time.Sleep(time.Millisecond)
}
