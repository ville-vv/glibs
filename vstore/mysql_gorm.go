package vstore

import "github.com/jinzhu/gorm"

type IMySqlGorm interface {
	DB() *gorm.DB
}

type MySqlStore struct {
}
