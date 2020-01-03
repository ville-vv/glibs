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
	e.times += 1
}

type Task struct {
	ps      bool          // 任务是否启动了
	scheme  string        // 方案名称
	runCnt  int32         // 当前任务数量
	ts      chan element  // 任务执行的单元
	inl     time.Duration // 对执行错误的任务循环处理一次的定时
	reLoops int           //
	cache   vqueue.Queue  //
}

func NewTask(num ...int) *Task {
	n := 10
	if len(num) > 0 {
		n = num[0]
	}
	return &Task{
		ts:      make(chan element, n),
		cache:   vqueue.New(),
		inl:     time.Second * 5,
		reLoops: 3,
	}
}

// Process is task process
func (t *Task) Process(ctx context.Context) {
	if t.ps {
		return
	}

	go t.loopCache(ctx)
	go t.process(ctx)
}

func (t *Task) process(ctx context.Context) {
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
				if err = t.cache.Push(elm); err == nil {
					break
				}
			}
			atomic.AddInt32(&t.runCnt, -1)

		case <-ctx.Done():
			t.ps = false
			return
		}
	}
}

func (t *Task) digestCache() {
	reloops := make([]element, 0, t.cache.Length())
	for t.cache.Length() > 0 {
		nd := t.cache.Shift()
		if nd != nil {
			switch elem := nd.(type) {
			case element:
				elem.IncTimes()
				if err := elem.tf(elem.params); err != nil {
					if elem.times < t.reLoops {
						// 执行错误再次重试
						reloops = append(reloops, elem)
						break
					}
				}
				atomic.AddInt32(&t.runCnt, -1)
			}
		}
	}

	// 重新放入缓存列表中
	for i := range reloops {
		_ = t.cache.Push(reloops[i])
	}
}

func (t *Task) loopCache(ctx context.Context) {
	tkr := time.NewTicker(t.inl)
	for {
		t.digestCache()
		select {
		case <-ctx.Done():
			return
		case <-tkr.C:
		}
	}
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
