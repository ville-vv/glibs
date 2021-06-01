package vtask

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type delayNode struct {
	Value interface{}
	Next  *delayNode
	TmStp int64 // 时间戳
}

func (sel *delayNode) SetTmStp(duration time.Duration) {
	sel.TmStp = time.Now().UnixNano() + int64(duration)
}

func newDelayNode(value interface{}, duration time.Duration) *delayNode {
	nd := &delayNode{Value: value}
	nd.SetTmStp(duration)
	return nd
}

// DelayQueue 延迟队列，使用环型的方式一直在循环检测队列中的数据是否超时
// 如果超时了就取出来
type DelayQueue struct {
	lock    sync.Mutex
	prtNode *delayNode
	head    *delayNode
	rear    *delayNode
	qLen    AtomicInt64
}

func NewDelayQueue() *DelayQueue {
	return &DelayQueue{}
}

func (sel *DelayQueue) Push(v interface{}, delay time.Duration) error {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	fmt.Println("Push:", v)
	node := newDelayNode(v, delay)
	if sel.head == nil || sel.qLen.Load() == 0 {
		sel.head = newDelayNode(0, 0)
		sel.head.Next = node
		sel.rear = node
		sel.prtNode = sel.head
	}
	sel.rear.Next = node
	sel.rear = node
	node.Next = sel.head.Next
	sel.qLen.Inc()
	return nil
}

func (sel *DelayQueue) Loop(ctx context.Context, f func(interface{})) {
	tmr := time.NewTicker(time.Millisecond * 10)
	var dataNode *delayNode
	for {
		select {
		case <-tmr.C:

			sel.lock.Lock()
			if sel.qLen.Load() <= 0 || sel.prtNode == nil {
				// 空的和当前长度为 0 就 直接跳过
				break
			}
			dataNode = sel.prtNode.Next
			if dataNode == nil {
				break
			}
			//
			if dataNode.TmStp > time.Now().UnixNano() {
				// 未超时就获取下一个
				sel.prtNode = sel.prtNode.Next
				break
			}
			f(dataNode.Value)
			sel.prtNode.Next = sel.prtNode.Next.Next
			dataNode.Next = nil
			sel.qLen.Dec()
			fmt.Println("当前队列大小：", sel.qLen.Load())
		case <-ctx.Done():
			tmr.Stop()
			sel.lock.Unlock()
			return
		}
		sel.lock.Unlock()
	}
}

func (sel *DelayQueue) del(prtNode *delayNode) {
	sel.lock.Lock()
	defer sel.lock.Unlock()
	sel.prtNode.Next = prtNode.Next
	prtNode = nil
	sel.qLen.Dec()
	return
}

func (sel *DelayQueue) Length() int64 {
	return sel.qLen.Load()
}
