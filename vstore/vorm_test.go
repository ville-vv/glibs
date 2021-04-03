package vstore

import (
	"testing"
)

func TestGormDb_CreateDB(t *testing.T) {
	db := MakeDBUtil(&dbConf{
		Username: "root",
		Password: "Root123.",
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "test_for_auto_create",
	})
	db.CreateDB()
	db.DropDB()
}
