package vredis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

type RedisGoDrive struct {
	rcnf  *RedisCnf
	rPool *redis.Pool
}

func NewRedisGoDrive(cnf *RedisCnf) *RedisGoDrive {
	r := &RedisGoDrive{
		rcnf: cnf,
	}
	return r
}

// 链接 Redis
func (sel *RedisGoDrive) Conn() error {
	sel.rPool = &redis.Pool{
		MaxActive:   sel.rcnf.MaxOpens,
		MaxIdle:     sel.rcnf.MaxIdles,
		Wait:        sel.rcnf.IsMaxConnWait,
		IdleTimeout: time.Duration(sel.rcnf.IdleTimeout),
		Dial: func() (conn redis.Conn, e error) {
			conn, e = redis.Dial("tcp", sel.rcnf.Address)
			if sel.rcnf.Password != "" {
				if _, err := conn.Do("AUTH", sel.rcnf.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			if sel.rcnf.UserName != "" {
				if _, err := conn.Do("CLIENT", "SETNAME", sel.rcnf.UserName); err != nil {
					conn.Close()
					return nil, err
				}
			}
			if sel.rcnf.Db != 0 {
				if _, err := conn.Do("SELECT", sel.rcnf.Db); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return
		},
	}
	return nil
}
func (sel *RedisGoDrive) Close() error {
	return sel.rPool.Close()
}
func (sel *RedisGoDrive) Cli() redis.Conn {
	return sel.rPool.Get()
}
func (sel *RedisGoDrive) Get(k string) (str string) {
	conn := sel.rPool.Get()
	res, err := conn.Do("GET", k)
	if err != nil {
		return
	}
	switch res.(type) {
	case []byte:
		str = string(res.([]byte))
	}
	return
}
func (sel *RedisGoDrive) Set(k string, v string, expiration int64) (err error) {
	conn := sel.rPool.Get()
	_, err = conn.Do("SET", k, v)
	if err != nil {
		return
	}
	_, err = conn.Do("PEXPIRE", k, expiration)
	if err != nil {
		return
	}
	return
}
