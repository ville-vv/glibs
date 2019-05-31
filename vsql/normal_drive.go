package vsql

import (
	"database/sql"
)

type NormalSqlDrive struct {
	SqlDrive
	dbConns map[string]*sql.DB
}

func NewNormalSqlDrive(cnf *MySqlCnf) (dri *NormalSqlDrive) {
	mcnf := &MySqlCnf{
		Version:  cnf.Version,
		UserName: cnf.UserName,
		Address:  cnf.Address,
		Password: cnf.Password,
		Default:  cnf.Default,
		MaxIdles: cnf.MaxIdles,
		MaxOpens: cnf.MaxOpens,
	}
	mcnf.Databases = append(mcnf.Databases, cnf.Databases[:]...)

	dri = &NormalSqlDrive{
		SqlDrive: SqlDrive{
			SqlCnf: mcnf,
		},
		dbConns: make(map[string]*sql.DB),
	}
	dri.Driver = dri
	return
}

func (sel *NormalSqlDrive) Conn(user, passwd, addr, dbname, version string, maxIdle, maxOpen int) error {
	cnStr := sel.SqlConnStr(version, user, passwd, addr, dbname)
	tempDb, err := sql.Open("mysql", cnStr)
	if err != nil {
		return err
	}

	if err := tempDb.Ping(); err != nil {
		tempDb.Close()
		return err
	}

	// 开启链接池，SetMaxOpenConns 设置最大链接数，
	tempDb.SetMaxOpenConns(maxOpen)
	// SetMaxIdleConns 用于设置闲置的连接数。
	tempDb.SetMaxIdleConns(maxIdle)
	sel.Lock.Lock()
	defer sel.Lock.Unlock()
	sel.dbConns[dbname] = tempDb
	return err
}
func (sel *NormalSqlDrive) BreakOff() {
	for _, v := range sel.dbConns {
		v.Close()
	}
}

func (sel *NormalSqlDrive) getDb(name string) *sql.DB {
	sel.Lock.RLock()
	defer sel.Lock.RUnlock()
	return sel.dbConns[name]
}

func (sel *NormalSqlDrive) GetDefDb() *sql.DB {
	return sel.getDb(sel.SqlCnf.Default)
}

func (sel *NormalSqlDrive) GetDb(name string) *sql.DB {
	return sel.getDb(name)
}
