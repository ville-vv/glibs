package vuid

import (
	"math/rand"
	"time"
)

var (
	uuid UUid
)

type UUid interface {
	Generate() (id int64)
}

func newGenerator(workId int64) (UUid, error) {
	return NewSnowFlake(workId)
}

// 获取 Uuid
func GenUUid() (id int64) {
	rand.Seed(time.Now().UnixNano())
	id, _ = genWithId(rand.Int63n(1023))
	return
}

// 指定一个自己的 work id 生成uuid, work 最大值为 1023
func genWithId(workId int64) (id int64, err error) {
	if uuid == nil {
		uuid, err = newGenerator(workId)
	}
	id = uuid.Generate()
	return
}
