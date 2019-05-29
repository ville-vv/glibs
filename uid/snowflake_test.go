package uid

import (
	"fmt"
	"testing"
	"time"
)

func TestEpoch(t *testing.T) {
	ep, _ := time.ParseDuration(time.Now().Format(TimeFormatLayout))

	fmt.Println(time.Now().Format(TimeFormatLayout), ep, time.Now().UnixNano(), "------", fmt.Sprintf("%d", time.Now().UnixNano()))
	fmt.Println(time.Now().UnixNano() / 1000000)
	fmt.Println(time.Now().Unix())
}

func TestSnowFlake_Generate(t *testing.T) {
	// 测试脚本

	// 生成节点实例
	node, err := NewSnowFlake(1)

	if err != nil {
		fmt.Println(err)
		return
	}

	ch := make(chan int64, 100000)
	count := 100000
	// 并发 count 个 goroutine 进行 snowflake ID 生成
	for i := 0; i < count; i++ {
		go func() {
			id := node.Generate()
			ch <- id
		}()
	}

	defer close(ch)
	startTime := time.Now().UnixNano()
	m := make(map[int64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		// 如果 map 中存在为 id 的 key, 说明生成的 snowflake ID 有重复
		//fmt.Println(id)
		_, ok := m[id]
		if ok {
			fmt.Printf("ID is not unique! count=%d, id=%d \n", i, id)
			return
		}
		// 将 id 作为 key 存入 map
		m[id] = i
	}
	entTime := time.Now().UnixNano()
	// 成功生成 snowflake ID
	fmt.Println("All ", count, " snowflake ID generate successed!", (entTime-startTime)/1000000)
}
