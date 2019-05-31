package vsql

import "testing"

func TestSql_Connector(t *testing.T) {
	sqlCnf := &MySqlCnf{
		Version:   "5",
		UserName:  "root",
		Address:   "localhost:3306",
		Password:  "Root123",
		Default:   "",
		MaxIdles:  100,
		MaxOpens:  1000,
		Databases: []string{},
	}
	norSql := NewNormalSqlDrive(sqlCnf)
	err := norSql.Open()
	if err != nil {
		t.Errorf("链接数据库错误：%v", err)
		return
	}
	res, err := norSql.GetDefDb().Prepare("show databases;")
	if err != nil {
		t.Errorf("执行数据库命令失败：%v", err)
		return
	}
	rows, err := res.Query()
	var dbname string
	for rows.Next() {
		if err = rows.Scan(&dbname); err != nil {
			t.Errorf("读取数据错误：%v", err)
			return
		}
		t.Log("读取到的数据库名称为：", dbname)
	}
}
