package errcodes

import (
	"fmt"
)

type ErrCodes int

// ErrCodes 错误字符串
func (e ErrCodes) ErrMsg() string {
	return errorMsg[int(e)]
}

func ErrInternal(msg interface{}) ErrCodes {
	//var errcode ErrCodes
	//errcode = (ErrCodes)(Code_InternalError)
	errorMsg[Code_InternalError] = fmt.Sprintf("%v", msg)
	return Code_InternalError
}

/*******************************************************************************
// @describe : 新建一个错误码对应的错误信息
// @param	: code 错误代码
// @param	: msg 错误消息
// @return	: ErrCodes 错误类
*******************************************************************************/
func NewErrCodes(code int, msg interface{}) ErrCodes {
	var errcode ErrCodes
	errcode = (ErrCodes)(code)
	errorMsg[code] = fmt.Sprintf("%v", msg)
	return errcode
}

const (
	Code_OK            = 200  //成功
	Code_PwdErr        = 1001 //密码错误
	Code_UserNotExist  = 1002 //用户不存在
	Code_InternalError = 5000 //内部错误(数据库，redis 等)
)

var (
	errorMsg = map[int]string{
		Code_OK:            "success",
		Code_PwdErr:        "password is error",
		Code_UserNotExist:  "user not exist",
		Code_InternalError: "internal error",
	}
)
