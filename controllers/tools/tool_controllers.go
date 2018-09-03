/*******************************************************************************
用户模块接口控制器
*******************************************************************************/
package toolControllers

import (
	"common/errcodes"
	"encoding/json"
	//"fmt"
	"vil_tools/controllers"
	"vil_tools/models/tools"
)

type ToolControllers struct {
	controllers.MainController
}

// @router /getCourrntTime [get]
func (u *ToolControllers) GetCourrntTime() {
	reqData := u.Ctx.Input.RequestBody
	map_data := make(map[string]interface{}, 0)
	json.Unmarshal(([]byte)(reqData), &map_data)
	var resData interface{}
	var code errcodes.ErrCodes
	defer func() {
		if code == errcodes.Code_OK {
			u.SendOk(resData)
		} else {
			u.SendError((int)(code), code.ErrMsg())
		}
	}()

	handle := toolModel.NewModel(toolModel.Cmd_CourrntTime)
	resData, code = handle.DoHandle(reqData)
	return
}

// @router /getWeChatOpenID [get,post]
func (u *ToolControllers) GetWeChatOpenID() {
	reqData := u.Ctx.Input.RequestBody
	map_data := make(map[string]interface{}, 0)
	json.Unmarshal(([]byte)(reqData), &map_data)
	var resData interface{}
	var code errcodes.ErrCodes
	defer func() {
		if code == errcodes.Code_OK {
			u.SendOk(resData)
		} else {
			u.SendError((int)(code), code.ErrMsg())
		}
	}()

	//根据Cmd_GetWeChatOpenID 获取不同的业务处理实例
	handle := toolModel.NewModel(toolModel.Cmd_GetWeChatOpenID)
	resData, code = handle.DoHandle(reqData)
	return
}
