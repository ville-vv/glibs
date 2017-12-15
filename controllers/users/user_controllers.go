/*******************************************************************************
用户模块接口控制器
*******************************************************************************/
package userControllers

import (
	"common/errcodes"
	"encoding/json"
	//"fmt"
	log "common/vilogs"
	"vil_tools/controllers"
	"vil_tools/models/users"
)

type UserControllers struct {
	controllers.MainController
}

// @router /login [post]
func (u *UserControllers) Login() {
	reqData := u.Ctx.Input.RequestBody
	map_data := make(map[string]interface{}, 0)
	json.Unmarshal(([]byte)(reqData), &map_data)
	var resData interface{}
	var code errcodes.ErrCodes
	defer func() {
		if code == errcodes.Code_OK {
			log.LOGI("response success %v", resData)
			u.SendOk(resData)
		} else {
			log.LOGE("response faild code=%v errmsg=%v", code, code.ErrMsg())
			u.SendError((int)(code), code.ErrMsg())
		}
	}()

	handle := userModel.NewUserModel(userModel.LoginModelCode)
	resData, code = handle.DoHandle(reqData)
	return
}
