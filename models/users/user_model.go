package userModel

import (
	"common/errcodes"
	"vil_tools/models/users/userLogin"
)

const (
	LoginModelCode int = iota
)

type IUserHandle interface {
	DoHandle(data []byte) (resDate interface{}, code errcodes.ErrCodes)
}

var (
	mp_Model map[int]IUserHandle
)

func init() {
	mp_Model = make(map[int]IUserHandle)
	mp_Model[LoginModelCode] = new(userLogin.UserLogin)
}

func NewUserModel(model int) IUserHandle {
	return mp_Model[model]
}
