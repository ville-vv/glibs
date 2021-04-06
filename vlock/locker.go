package vlock

import "time"

type Locker interface {
	Lock(key string, timeout time.Duration) (string, error)
	UnLock(key string, val string) error
}

//
type Interceptor interface {
	Intercept(key string, timeout time.Duration) error
}
