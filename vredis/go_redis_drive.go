package vredis

import (
	"github.com/go-redis/redis"
	"time"
)

type GoRedisDrive struct {
	rcnf   *RedisCnf
	client *redis.Client
}

func NewGoRedisDrive(cnf *RedisCnf) *GoRedisDrive {
	return &GoRedisDrive{
		rcnf: cnf,
	}
}

func (sel *GoRedisDrive) Conn() error {
	sel.client = redis.NewClient(&redis.Options{
		Addr:         sel.rcnf.Address,
		IdleTimeout:  time.Duration(sel.rcnf.IdleTimeout),
		Password:     sel.rcnf.Password,
		PoolSize:     sel.rcnf.MaxOpens,
		MinIdleConns: sel.rcnf.MaxIdles,
	})
	return nil
}
func (sel *GoRedisDrive) Close() error {
	return sel.client.Close()
}
func (sel *GoRedisDrive) Cli() *redis.Client {
	return sel.client
}
func (sel *GoRedisDrive) Get(k string) string {
	return sel.client.Get(k).Val()
}
func (sel *GoRedisDrive) Set(k, v string, expiration uint64) error {
	return sel.client.Set(k, v, time.Millisecond*time.Duration(expiration)).Err()
}
