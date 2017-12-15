package models

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
)

func ToMd5(arg_str string) string {
	hd := md5.New()
	hd.Write([]byte(arg_str))
	return hex.EncodeToString(hd.Sum(nil))
}

//map[string]string 类型排序
func MapSorted(map_array map[string]string) map[string]string {
	var keys []string
	for k, _ := range map_array {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var resMap map[string]string
	resMap = make(map[string]string)
	for _, k := range keys {
		resMap[k] = map_array[k]
	}
	return resMap
}

//进行签名
func ToSign(agr_str map[string]string, user_key string) string {

	var tosignStr string = ""
	var num int = 0

	//先key排序
	agr_str = MapSorted(agr_str)
	//组合为 key=&value&key2=vale2&user_key
	for k, v := range agr_str {
		if num == 0 {
			tosignStr = fmt.Sprintf("%v=%v", k, v)
		} else {
			tosignStr = fmt.Sprintf("%s&%v=%v", tosignStr, k, v)
		}
		num++
	}
	tosignStr = fmt.Sprintf("%s&%v", tosignStr, user_key)
	fmt.Println("tosignStr=", tosignStr)
	return ToMd5(tosignStr)
}
