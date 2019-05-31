package vredis

import (
	"testing"
)

func TestRedisDrive_Conn(t *testing.T) {
	rds := NewRedisGoDrive(&RedisCnf{
		Address:  "127.0.0.1:6379",
		MaxIdles: 100,
		MaxOpens: 1000,
	})
	defer rds.Close()
	if err := rds.Conn(); err != nil {
		t.Errorf("Redis 链接错误:%v", err.Error())
		return
	}
	if err := rds.Set("SUNNAME", "redis_ville", 5000); err != nil {
		t.Errorf("Redis 设置值错误:%v", err.Error())
		return
	}
	t.Log("获取到值：", rds.Get("SUNNAME"))
}
