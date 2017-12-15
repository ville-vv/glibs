package vilMysql

import (
	"fmt"
	"log"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type VilMysql struct {
	HostName  string
	UserName  string
	Password  string
	DataBases map[string]string
}

func RegisterMysql(mysql *VilMysql) (err error) {
	log.Println("begin to connect mysql...")
	if err = orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		log.Println("RegisterDriver err:", err)
		return
	}
	dataBaseSource := fmt.Sprintf("%v:%v@(%v)/%v?charset=utf8",
		mysql.UserName,
		mysql.Password,
		mysql.HostName,
		mysql.DataBases["default"])
	if err = orm.RegisterDataBase("default", "mysql", dataBaseSource); err != nil {
		log.Println("RegisterDataBase err:", err)
		return
	}
	log.Println("mysql connect success!")
	return
}
