package userLogin

import (
	//	"common/config"
	"common/errcodes"
	"encoding/json"
	"fmt"
	//"time"

	log "common/vilogs"
	"vil_tools/models"

	"github.com/astaxie/beego/orm"
	//"errors"
	//"fmt"
)

type UserLogin struct {
}

type RequestParam struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	CurrTime string `json:"curr_time"`
	Sign     string `json:"sign"`
}

type VilUserPrivacy struct {
	UserID   string `orm:"column(user_id)" json:"user_id"`
	Password string `orm:"column(password)" json:"password"`
	UserKey  string `orm:"column(user_key)" json:"user_key"`
}

/*******************************************************************************
// @describe :
// @param	: data 请求入参
// @return	: resp 返回参数
// @return	: code 返回错误码
*******************************************************************************/
func (ul *UserLogin) DoHandle(data []byte) (interface{}, errcodes.ErrCodes) {
	var reqParam RequestParam
	var resp string = "success"
	json.Unmarshal(([]byte)(data), &reqParam)

	code := login(&reqParam)
	return resp, code
}

func login(param *RequestParam) (code errcodes.ErrCodes) {
	ormMysql := orm.NewOrm()
	ormMysql.Using("default")
	sqlcmd := fmt.Sprintf("SELECT user_id, password, user_key FROM vil_user.user_privacy WHERE user_id=%d", param.UserID)
	var userinfo []VilUserPrivacy
	num, err := ormMysql.Raw(sqlcmd).QueryRows(&userinfo)
	if err != nil {
		log.LOGI("sql exec error:%v", err)
		code = errcodes.ErrInternal(err)
		return
	}
	if num == 0 {
		log.LOGI("error: user not exist")
		code = errcodes.Code_UserNotExist
		return
	}
	map_res := make(map[string]string)
	map_res["user_id"] = userinfo[0].UserID
	map_res["password"] = userinfo[0].Password
	map_res["curr_time"] = param.CurrTime
	log.LOGI("sign=%v", models.ToSign(map_res, userinfo[0].UserKey))
	return errcodes.Code_OK
}
