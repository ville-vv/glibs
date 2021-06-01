package vutil

import (
	"fmt"
	"testing"
)

func TestGenVCode(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(GenVCode(10))
	}

}
