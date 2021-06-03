// @File     : task
// @Author   : Ville
// @Time     : 19-10-9 上午9:08
// vtask
package vtask

import (
	"context"
	"errors"
	"fmt"
	"github.com/vilsongwei/vilgo/vqueue"
	"sync/atomic"
	"time"
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
	cache   Queue         //
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
		nd := t.cache.Pop()
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

// 重试机制
// 执行循序
// 分布式全局锁
// 自动扩容工作池
// 可持久存储

// 持久化存储
type Persistent interface {
	Load(interface{}) (int, error)
	Store(interface{}) error
}

type TaskOption struct {
	MaxQueueNum    int   // 最列表务数
	RetryNum       int   // 重试次数 0 不重试
	RetryInterval  int64 // 重试间隔时间，最小单位秒
	PersistentFlag bool  // 是否持久化
	ErrDealFunc    func(ctx interface{}, err error)
}

type retryTask struct {
	LastErr  error
	retryNum int
	Data     interface{}
}

func (sel *retryTask) Dec() {
	sel.retryNum--
	if sel.retryNum < 0 {
		sel.retryNum = 0
	}
}

func (sel *retryTask) Warn() bool {
	return sel.retryNum <= 0
}

type MiniTask struct {
	option    *TaskOption
	dataList  Queue
	retryList DelayQueue // 延迟重试
	pushCh    chan interface{}
	retryCh   chan interface{}
	New       func() interface{}
	pst       Persistent
}

func (t *MiniTask) Start() {
	t.retryList.Run()
}

func (t *MiniTask) Stop() {
}

func (t *MiniTask) loopPush() {
	for val := range t.pushCh {
		err := t.dataList.Push(val)
		if err != nil {
			// 出现错误先 放入重试
			Err := t.Retry(val)
			if Err != nil {
				// 放入重试也错误了，就调用错误处理函数
				t.option.ErrDealFunc(val, fmt.Errorf("%s,%s", err.Error(), Err.Error()))
			}
		}
	}
}

func (t *MiniTask) loopRetry() {
	var err error
	for val := range t.retryCh {
		nextTime := time.Now().Add(time.Second * time.Duration(t.option.RetryInterval))
		switch retryDt := val.(type) {
		case *retryTask:
			if retryDt.Warn() {
				t.option.ErrDealFunc(retryDt.Data, retryDt.LastErr)
				break
			}
			retryDt.Dec()
			err = t.retryList.Push(retryDt, nextTime)
			if err != nil {
				t.option.ErrDealFunc(retryDt.Data, err)
			}
		default:
			err = t.retryList.Push(&retryTask{
				retryNum: t.option.RetryNum,
				Data:     val,
			}, nextTime)
			if err != nil {
				t.option.ErrDealFunc(val, err)
			}
		}
	}
}

func (t *MiniTask) Push(ctx interface{}) error {
	select {
	case t.pushCh <- ctx:
	default:

	}
	return nil
}

func (t *MiniTask) Retry(ctx interface{}) error {
	select {
	case t.retryCh <- ctx:
	default:
		if t.option.PersistentFlag {
			// 如果满了就持久化存储
			return t.pst.Store(ctx)
		}
		// 如果么有开启持久处理，直接返回错误
		return errors.New("retry list is full")
	}
	return nil
}
