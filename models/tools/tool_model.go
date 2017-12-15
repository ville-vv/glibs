package toolModel

import (
	"common/errcodes"
)

type IToolHandle interface {
	DoHandle(data []byte) (resDate interface{}, code errcodes.ErrCodes)
}

var (
	mp_Model map[int]IToolHandle
)

const (
	Cmd_CourrntTime int = iota
)

func init() {
	mp_Model = make(map[int]IToolHandle)
	mp_Model[Cmd_CourrntTime] = new(CourrntTime)
}

func NewModel(model int) IToolHandle {
	return mp_Model[model]
}
