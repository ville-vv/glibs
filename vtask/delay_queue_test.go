package vtask

import (
	"context"
	"fmt"
	"github.com/ville-vv/vilgo/vutil"
	"testing"
	"time"
)

func TestDelayQueue_Loop(t *testing.T) {
	dqu := DelayQueue{}
	inc := AtomicInt64{}

	//for i := 0; i < 20; i++ {
	//	go func(n int) {
	//		dqu.Loop(context.Background(), func(val interface{}) {
	//			if val == nil {
	//				t.Error()
	//				panic("")
	//			}
	//			fmt.Println("消费者:", n, val)
	//		})
	//	}(i)
	//}
	go func() {
		dqu.Loop(context.Background(), func(i interface{}) {
			if i == nil {
				t.Error()
				panic("")
			}
			fmt.Println("消费者1:", i)
		})
	}()

	go func() {
		for {
			time.Sleep(time.Millisecond)
			val := vutil.GenVCode(7)
			fmt.Println("生产者1:", val)
			dqu.Push(val, time.Millisecond)
			inc.Inc()
			if inc.Load() > 10888 {
				return
			}
		}

	}()

	go func() {
		for {
			time.Sleep(time.Millisecond)
			val := vutil.GenVCode(7)
			fmt.Println("生产者1:", val)
			dqu.Push(val, time.Millisecond)
			inc.Inc()
			if inc.Load() > 1088 {
				return
			}
		}

		//for i := 1; i < 8; i++ {
		//	time.Sleep(time.Millisecond)
		//	dqu.Push(i, time.Millisecond)
		//}

	}()
	time.Sleep(time.Second)

	for i := 1; i < 8; i++ {
		time.Sleep(time.Millisecond)
		dqu.Push(i, time.Millisecond)
	}

	select {}

}
