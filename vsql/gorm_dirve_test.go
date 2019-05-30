package vsql

import "testing"

func TestGormDb_Connector(t *testing.T) {
	sqlCnf := &MySqlCnf{
		Version:      "8",
		UserName:     "root",
		Address:      "localhost:3306",
		Password:     "Root123",
		Default:      "",
		MaxIdleConns: 100,
		MaxOpenConns: 1000,
		Databases:    []string{},
	}
	norSql := NewGormDb(sqlCnf)
	err := norSql.Open()
	if err != nil {
		t.Errorf("链接数据库错误：%v", err)
		return
	}
}
