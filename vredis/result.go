package vredis

import (
	"fmt"
	"strconv"
)

type Result struct {
	val interface{}
}

func (sel *Result) Int() (int, error) {
	switch sel.val.(type) {
	case []byte:
		return strconv.Atoi(string(sel.val.([]byte)))
	default:
		err := fmt.Errorf("redis: type error %v", sel.val)
		return 0, err
	}
}
