package toolModel

import (
	"common/errcodes"
	"encoding/json"
	"time"
)

type GetWeChatOpenID struct{

}

/**
 * data 请求的数据
 */
func(self * GetWeChatOpenID)DoHandle(data []byte) (interface{}, errcodes.ErrCodes){
	//把请求的参数转换为map数据格式
	map_data := make(map[string]interface{}, 0)
	json.Unmarshal(([]byte)(data), &map_data)

	return time.Now(), errcodes.Code_OK
}
