package vuid

import (
	"testing"
	"time"
)

func TestGenUUid(t *testing.T) {
	t.Log("生成的Uuid:", GenUUid())
}

func Test_genWithId(t *testing.T) {
	count := 0
	idMap2 := make(map[int64]bool)
	start := time.Now().UnixNano() / 1e6
	for {
		end := time.Now().UnixNano() / 1e6
		if end-start >= 100 {
			t.Log(start, "---", end)
			break
		}
		count++
		id, _ := genWithId(1)
		if _, ok := idMap2[id]; ok {
			t.Errorf("生成出重复ID: %d", id)
		}
		idMap2[id] = true
	}
	t.Log("生成的count:", count)
}
