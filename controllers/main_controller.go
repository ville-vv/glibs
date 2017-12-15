package controllers

import (
	"common/response"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

/*******************************************************************************
// @param    : 接口需要返回的数据
// @describe : 接口执行成功返回的数据
// @return   : 无
*******************************************************************************/
func (mc *MainController) SendOk(dataok interface{}) {
	res := &response.Response{}
	res.SendOk(dataok)
	mc.Data["json"] = res.GetOK()
	mc.ServeJSON()
	return
}

/*******************************************************************************
// @param    : code 返回的错误码，msg返回错误消息
// @describe : 接口执行失败调用
// @return   : 无
*******************************************************************************/
func (mc *MainController) SendError(code int, msg string) {
	res := &response.Response{}
	res.SendError(code, msg)
	mc.Data["json"] = res.GetError()
	mc.ServeJSON()
	return
}
