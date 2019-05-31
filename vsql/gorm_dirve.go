package vsql

import (
	"github.com/jinzhu/gorm"
)

type GormDb struct {
	SqlDrive
	dbConns map[string]*gorm.DB
}

func NewGormDb(conf *MySqlCnf) (dri *GormDb) {
	if conf == nil {
		panic("mysql config is nil")
	}
	mcnf := &MySqlCnf{
		Version:  conf.Version,
		UserName: conf.UserName,
		Address:  conf.Address,
		Password: conf.Password,
		Default:  conf.Default,
		MaxIdles: conf.MaxIdles,
		MaxOpens: conf.MaxOpens,
	}
	mcnf.Databases = append(mcnf.Databases, conf.Databases[:]...)

	dri = &GormDb{
		dbConns: make(map[string]*gorm.DB),
		SqlDrive: SqlDrive{
			SqlCnf: mcnf,
		},
	}
	dri.Driver = dri
	return
}

func (sel *GormDb) Conn(user, passwd, addr, dbname, version string, maxIdle, maxOpen int) error {
	cnStr := sel.SqlConnStr("8", user, passwd, addr, dbname)
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
	sel.Lock.Lock()
	defer sel.Lock.Unlock()
	sel.dbConns[dbname] = tempDb
	return nil
}
func (sel *GormDb) BreakOff() {
	for _, v := range sel.dbConns {
		v.Close()
	}
}

// 获取默认连接的数据库对象
func (sel *GormDb) GetDefDb() *gorm.DB {
	return sel.getDb(sel.SqlCnf.Default)
}

// 根据数据库名称获取数据库对象
func (sel *GormDb) GetDb(dbName string) *gorm.DB {
	return sel.getDb(dbName)
}

func (sel *GormDb) getDb(dbName string) *gorm.DB {
	sel.Lock.RLock()
	defer sel.Lock.RUnlock()
	return sel.dbConns[dbName]
}
