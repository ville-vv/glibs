package mysql

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type GormDb struct {
	mcnf       *MySqlCnf
	mysqlConns map[string]*gorm.DB
	lock       sync.RWMutex
}

func NewGormDb(conf *MySqlCnf) *GormDb {
	if conf == nil {
		panic("mysql config is nil")
	}
	mcnf := &MySqlCnf{
		UserName:     "",
		Address:      "",
		Password:     "",
		Default:      "",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
	}
	mcnf.Databases = append(mcnf.Databases, conf.Databases[:]...)
	return &GormDb{
		mcnf:       mcnf,
		mysqlConns: make(map[string]*gorm.DB),
	}
}

func (sel *GormDb) open(user, passwd, addr, dbname string, maxIdle, maxOpen int) error {
	cnStr := SqlConnStr("8", user, passwd, addr, dbname)
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
	sel.lock.Lock()
	defer sel.lock.Unlock()
	sel.mysqlConns[dbname] = tempDb
	return nil
}

func (sel *GormDb) Open() error {
	// 连接默认数据库
	if err := sel.open(sel.mcnf.UserName, sel.mcnf.Password, sel.mcnf.Address, sel.mcnf.Default, sel.mcnf.MaxIdleConns, sel.mcnf.MaxOpenConns); err != nil {
		return err
	}
	// 连接其他指定的数据库
	for _, dbv := range sel.mcnf.Databases {
		if err := sel.open(sel.mcnf.UserName, sel.mcnf.Password, sel.mcnf.Address, dbv, sel.mcnf.MaxIdleConns, sel.mcnf.MaxOpenConns); err != nil {
			return err
		}
	}
	return nil
}

// 获取默认连接的数据库对象
func (sel *GormDb) GetDefDb() *gorm.DB {
	return sel.getDb(sel.mcnf.Default)
}

// 根据数据库名称获取数据库对象
func (sel *GormDb) GetDb(dbName string) *gorm.DB {
	return sel.getDb(dbName)
}

func (sel *GormDb) getDb(dbName string) *gorm.DB {
	sel.lock.RLock()
	defer sel.lock.RUnlock()
	return sel.mysqlConns[dbName]
}
