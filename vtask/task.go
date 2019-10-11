// @File     : task
// @Author   : Ville
// @Time     : 19-10-9 上午9:08
// vtask
package vtask

import (
	"context"
	"sync/atomic"
	"time"
	"vilgo/vqueue"
)

type TaskFunc func(params interface{}) error

type element struct {
	params interface{}
	times  int
	tf     TaskFunc
}

func (e *element) IncTimes() {

}

type Task struct {
	ps     bool         // 任务是否启动了
	scheme string       // 方案名称
	runCnt int32        // 当前任务数量
	ts     chan element // 任务执行的单元
	inl    time.Duration
	cache  vqueue.Queue
}

func NewTask(num ...int) *Task {
	n := 10
	if len(num) > 0 {
		n = num[0]
	}
	return &Task{
		ts:    make(chan element, n),
		cache: vqueue.New(),
		inl:   time.Second * 5,
	}
}

// Process is task process
func (t *Task) Process(ctx context.Context) {
	if t.ps {
		return
	}

	t.loopCache(ctx)

	go func() {
		t.ps = true
		for {
			select {
			case elm, ok := <-t.ts:
				if !ok {
					t.ps = false
					return
				}
				if err := elm.tf(elm.params); err != nil {
					//任务执行失败
					_ = t.cache.Push(elm)
					break
				}
				atomic.AddInt32(&t.runCnt, -1)

			case <-ctx.Done():
				t.ps = false
				return
			}
		}
	}()
}

func (t *Task) loopCache(ctx context.Context) {
	go func() {
		tkr := time.NewTicker(t.inl)

		for {
			for t.cache.Length() > 0 {
				nd := t.cache.Shift()
				if nd != nil {
					switch elem := nd.(type) {
					case element:
						_ = elem.tf(elem.params)
					}
					atomic.AddInt32(&t.runCnt, -1)
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-tkr.C:
			}
		}
	}()
	return
}

// Add a task , the param will pass on to tf
func (t *Task) Add(param interface{}, tf TaskFunc) {
	atomic.AddInt32(&t.runCnt, 1)
	t.ts <- element{params: param, tf: tf}
}

// CanStop return true is can stop task
func (t *Task) CanStop() bool {
	return atomic.LoadInt32(&t.runCnt) <= 0
}
