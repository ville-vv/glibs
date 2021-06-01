package vtask

import (
	"errors"
	"sync"
)

type RingQueue struct {
	lock   sync.Mutex
	list   []interface{}
	qLen   AtomicInt64
	qCap   int64
	ptrIdx AtomicInt64 // 指定当前元素存放位置
	popIdx AtomicInt64
}

func NewRingQueue(qCap int64) *RingQueue {
	return &RingQueue{
		list:   make([]interface{}, qCap+1),
		qCap:   qCap,
		ptrIdx: AtomicInt64{},
		popIdx: AtomicInt64{},
	}
}

func (r *RingQueue) Push(val interface{}) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.qLen.Load() > r.qCap {
		return errors.New("")
	}
	r.list[r.ptrIdx.Load()] = val
	r.ptrIdx.Inc()
	r.qLen.Inc()
	// 如果指针索引超过 cap 就指向 0
	if r.ptrIdx.Load() >= r.qCap {
		r.ptrIdx.Store(0)
	}
	return nil
}

func (r *RingQueue) Pop() interface{} {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.qLen.Load() == 0 {
		return nil
	}
	val := r.list[r.popIdx.Load()]
	r.list[r.popIdx.Load()] = nil
	r.qLen.Dec()
	r.popIdx.Inc()
	if r.popIdx.Load() >= r.qCap {
		r.popIdx.Store(0)
	}
	return val
}

func (r *RingQueue) Length() int64 {
	return r.qLen.Load()
}
