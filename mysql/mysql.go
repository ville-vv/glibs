package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	MaxIdleConnsDefautl = 100
	MaxOpenConnsDefault = 1000
)

type MySqlCnf struct {
	UserName     string   `json:"user_name"`
	Address      string   `json:"host"`
	Password     string   `json:"password"`
	Default      string   `json:"default"`
	MaxIdleConns int      `json:"max_idle_conns"`
	MaxOpenConns int      `json:"max_open_conns"`
	Databases    []string `json:"databases"`
}

func SqlConnStr(version, user, passwd, addr, dbname string) (cs string) {
	switch version {
	case "5":
		cs = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, addr, dbname)
	case "8":
		cs = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true", user, passwd, addr, dbname)
	default:
		cs = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, passwd, addr, dbname)
	}
	return
}
