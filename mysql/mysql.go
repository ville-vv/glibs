package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
)

const (
	MaxIdleConnsDefautl = 100
	MaxOpenConnsDefault = 1000
)

var (
	mysqlConns    = make(map[string]*gorm.DB)
	lock          sync.RWMutex
	defaultDbName = ""
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

func open(user, passwd, addr, dbname string, maxIdle, maxOpen int) error {
	cnStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true", user, passwd, addr, dbname)
	tempDb, err := gorm.Open("mysql", cnStr)
	if err != nil {
		return err
	}

	if maxIdle <= 0 {
		tempDb.DB().SetMaxIdleConns(MaxIdleConnsDefautl)
	}
	if maxOpen <= 0 {
		tempDb.DB().SetMaxOpenConns(MaxOpenConnsDefault)
	}
	lock.Lock()
	mysqlConns[dbname] = tempDb
	defer lock.Unlock()
	return nil
}

func Open(conf *MySqlCnf) error {
	defaultDbName = conf.Default
	// 连接默认数据库
	if err := open(conf.UserName, conf.Password, conf.Address, conf.Default, conf.MaxIdleConns, conf.MaxOpenConns); err != nil {
		return err
	}
	// 连接其他指定的数据库
	for _, dbv := range conf.Databases {
		if err := open(conf.UserName, conf.Password, conf.Address, dbv, conf.MaxIdleConns, conf.MaxOpenConns); err != nil {
			return err
		}
	}
	return nil
}

// 获取默认连接的数据库对象
func GetDefDb() *gorm.DB {
	return getDb(defaultDbName)
}

// 根据数据库名称获取数据库对象
func GetDb(dbName string) *gorm.DB {
	return getDb(dbName)
}

func getDb(dbName string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return mysqlConns[dbName]
}
