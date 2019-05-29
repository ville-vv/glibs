package mysql

import "database/sql"

type DbSqlDb struct {
}

func (sel *DbSqlDb) open(user, passwd, addr, dbname string, maxIdle, maxOpen int) error {
	cnStr := SqlConnStr("8", user, passwd, addr, dbname)
	tempDb, err := sql.Open("mysql", cnStr)
	if err != nil {
		panic(err.Error())
	}
	// 开启链接池，SetMaxOpenConns 设置最大链接数， SetMaxIdleConns 用于设置闲置的连接数。
	tempDb.SetMaxOpenConns(1000)
	tempDb.SetMaxIdleConns(100)
	// TODO
	return err
}

func (sel *DbSqlDb) Open() {

}
