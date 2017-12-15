package toolModel

import (
	//	"common/config"
	"common/errcodes"
	"encoding/json"
	"time"
)

type CourrntTime struct {
}

func (ul *CourrntTime) DoHandle(data []byte) (interface{}, errcodes.ErrCodes) {
	map_data := make(map[string]interface{}, 0)
	json.Unmarshal(([]byte)(data), &map_data)
	return time.Now(), errcodes.Code_OK
}
