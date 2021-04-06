package vlock

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisLock struct {
	rds     *redis.Client
	timeout time.Duration
}

func NewRedisLock(rds *redis.Client) *RedisLock {
	return &RedisLock{
		rds: rds,
	}
}

func (r *RedisLock) lockKey(k string) string {
	return fmt.Sprintf("vlock:%s", k)
}

func (r *RedisLock) interceptKey(k string) string {
	return fmt.Sprintf("intercept::%s", k)
}

func (r *RedisLock) Lock(key string, timeout time.Duration) (string, error) {
	val := fmt.Sprintf("%d", time.Now().UnixNano())
	for {
		res, err := r.rds.SetNX(r.lockKey(key), val, timeout).Result()
		if err != nil {
			return val, err
		}
		if res {
			return val, nil
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (r *RedisLock) UnLock(key string, val string) error {
	key = r.lockKey(key)
	script := `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[2])
else
	return 0
end
`
	_, err := r.rds.Eval(script, []string{key, key}, val).Result()
	if err != nil {
		return err
	}
	return nil

}

//
func (r *RedisLock) Intercept(key string, timeout time.Duration) error {
	val := fmt.Sprintf("%d", time.Now().UnixNano())
	res, err := r.rds.SetNX(r.interceptKey(key), val, timeout).Result()
	if err != nil {
		return err
	}
	if !res {
		return errors.New("too many times, please try again later")
	}
	return nil
}
