package config

import (
	"flag"
	"path"

	"fmt"
	"io/ioutil"
	"os"

	"common/vilmysql"
	"common/vilogs"

	"github.com/BurntSushi/toml"
)

var (
	ProPath   = flag.String("home", "E:/GoWork/src/vil_tools/", "the project home path")
	AllConfig *Config
)

//type MysqlConfig struct {
//	HostName  string
//	UserName  string
//	Password  string
//	DataBases map[string]string
//}

type RedisConfig struct {
	HostName []string
	Password string
	TimeOut  int
}
type ProConfig struct {
	ProName string
}

//配置文件参数保存结构体
type Config struct {
	Project *ProConfig
	Mysql   *vilMysql.VilMysql
	Redis   *RedisConfig
	ViLog   *viLogs.ViLogsConfig
}

//读取配置文件
func NewConfig(fileName string) (conf *Config, err error) {
	var file *os.File
	var buf []byte

	conf = new(Config)
	if file, err = os.Open(fileName); err != nil {
		return
	}

	if buf, err = ioutil.ReadAll(file); err != nil {
		return
	}
	if err = toml.Unmarshal(buf, conf); err != nil {
		return
	}
	return
}

func ShowConfigInfo(conf *Config) {
	fmt.Println("Project config", *conf.Project)
	fmt.Println("MySql config", *conf.Mysql)
	fmt.Println("Redis config", *conf.Redis)
	fmt.Println("ViLog config", *conf.ViLog)
}

func init() {
	var err error
	AllConfig, err = NewConfig(path.Join(*ProPath, "./conf/config.toml"))
	if err != nil {
		fmt.Println("config.toml read failed.")
		panic(-1)
	}
	ShowConfigInfo(AllConfig)
	vilMysql.RegisterMysql(AllConfig.Mysql)
	viLogs.InitLogs(AllConfig.ViLog, AllConfig.Project.ProName)
}
