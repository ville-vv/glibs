package vstore

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"time"
)

var logConfig = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
	SlowThreshold: 1 * time.Second,
	LogLevel:      logger.Warn,
	Colorful:      true,
})

func MakeDb(dbConfig DbConfig, replicas ...DbConfig) DB {
	return NewGormDb(dbConfig, false, replicas...)
}

func MakeDBUtil(dbConfig DbConfig) DBUtil {
	return NewGormDb(dbConfig, true)
}

type dbConf struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DbName       string
	MaxIdleConns int
	MaxOpenConns int
	DbCharset    string
}

func (sel *dbConf) GetUserName() string {
	return sel.Username
}

func (sel *dbConf) GetPassword() string {
	return sel.Password
}

func (sel *dbConf) GetHost() string {
	return sel.Host
}

func (sel *dbConf) GetPort() string {
	return sel.Port
}

func (sel *dbConf) GetDbName() string {
	return sel.Port
}

func (sel *dbConf) GetMaxIdleConn() int {
	return sel.MaxIdleConns
}

func (sel *dbConf) GetMaxOpenConn() int {
	return sel.MaxOpenConns
}

func (sel *dbConf) GetCharset() string {
	return sel.DbCharset
}

func (sel *dbConf) ToSet(db DbConfig) {
	sel.Username = db.GetUserName()
	sel.Password = db.GetPassword()
	sel.Port = db.GetPort()
	sel.Host = db.GetHost()
	sel.DbName = db.GetDbName()
	sel.MaxIdleConns = db.GetMaxIdleConn()
	sel.MaxOpenConns = db.GetMaxOpenConn()
	sel.DbCharset = db.GetCharset()
	if sel.DbCharset == "" {
		sel.DbCharset = "utf8"
	}
}

type DbConfig interface {
	GetUserName() string
	GetPassword() string
	GetHost() string
	GetPort() string
	GetDbName() string
	GetMaxIdleConn() int
	GetMaxOpenConn() int
	GetCharset() string
}

type GormDb struct {
	dbConfig *dbConf
	replicas []*dbConf
	db       *gorm.DB
	utilDB   *gorm.DB
}

func NewGormDb(dbC DbConfig, forUtil bool, repcs ...DbConfig) *GormDb {
	mainDb := &dbConf{}
	mainDb.ToSet(dbC)
	replicas := make([]*dbConf, len(repcs))
	for i, rep := range replicas {
		rep.ToSet(repcs[i])
	}
	return newGormDb(mainDb, forUtil, replicas...)
}

func newGormDb(dbConfig *dbConf, forUtil bool, replicas ...*dbConf) *GormDb {
	gm := &GormDb{dbConfig: dbConfig, replicas: replicas}
	if forUtil {
		gm.initCdDb()
		return gm
	}
	// init db
	gm.initGormDB()
	return gm
}

// information_schema 数据库表
func (gm *GormDb) initCdDb() {
	if gm.db != nil {
		panic("gorm db should nil")
	}
	// 链接mysql 固定的数据库
	cStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, "information_schema", gm.dbConfig.DbCharset)
	openedDb, err := gorm.Open(mysql.Open(cStr), &gorm.Config{Logger: logConfig})
	if err != nil {
		panic("连接数据库出错:" + err.Error())
	}
	gm.utilDB = openedDb
}

func (gm *GormDb) initGormDB() {
	if gm.db != nil {
		panic("gorm db should nil")
	}

	mysqlDialer := mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName, gm.dbConfig.DbCharset))
	gormLogLevel := logger.Silent
	// 如果不是生产数据库则打开详细日志
	//if substr(gm.dbConfig.DbName, len(gm.dbConfig.DbName)-4, 4) != "prod" {
	//	gormLogLevel = logger.Info
	//}
	openedDb, err := gorm.Open(mysqlDialer, &gorm.Config{
		Logger: logConfig.LogMode(gormLogLevel),
	})
	if err != nil {
		panic("数据库连接出错：" + err.Error())
	}
	dbPool, err := openedDb.DB()
	if err != nil {
		panic("获取数据库连接池出错：" + err.Error())
	}
	dbPool.SetMaxIdleConns(gm.dbConfig.MaxIdleConns)
	dbPool.SetMaxOpenConns(gm.dbConfig.MaxOpenConns)
	// 避免久了不使用，导致连接被mysql断掉的问题
	dbPool.SetConnMaxLifetime(time.Hour * 1)

	gm.db = openedDb

	// 读分离 - 丛库
	gm.addReplicas()
}

func (gm *GormDb) addReplicas() {
	if len(gm.replicas) <= 0 {
		return
	}
	var replicaDialers = make([]gorm.Dialector, 0, len(gm.replicas))
	for _, config := range gm.replicas {
		if config == nil {
			return
		}
		cStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", gm.dbConfig.Username, gm.dbConfig.Password, gm.dbConfig.Host, gm.dbConfig.Port, gm.dbConfig.DbName, gm.dbConfig.DbCharset)
		replicaDialers = append(replicaDialers, mysql.Open(cStr))
	}
	err := gm.db.Use(
		dbresolver.Register(dbresolver.Config{
			Replicas: replicaDialers,
			Policy:   dbresolver.RandomPolicy{}}). // 随机负载均衡
			SetMaxIdleConns(gm.dbConfig.MaxIdleConns).
			SetMaxOpenConns(gm.dbConfig.MaxOpenConns).
			SetConnMaxLifetime(time.Hour * 1))
	if err != nil {
		panic("初始化丛库出错: " + err.Error())
	}
}

func (gm *GormDb) ClearAllData() {

	dbEnv := os.Getenv("DB_ENV")
	if dbEnv != "test" {
		panic("not allow to clear all data in un-test env")
	}

	tmpDb := gm.db
	if tmpDb == nil {
		panic("please init databases")
	}
	rs, err := tmpDb.Raw("show tables;").Rows()
	if err != nil {
		panic("get table list failure：" + err.Error())
	}
	var tName string
	for rs.Next() {
		if err := rs.Scan(&tName); err != nil || tName == "" {
			panic("get table name failure")
		}
		if err := tmpDb.Exec(fmt.Sprintf("delete from %s", tName)).Error; err != nil {
			panic("clear table data failure:" + err.Error())
		}
	}
}

func (gm *GormDb) GetDB() *gorm.DB {
	return gm.db
}

func (gm *GormDb) CreateDB() {
	createDbSQL := "CREATE DATABASE IF NOT EXISTS " + gm.dbConfig.DbName + " DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;"

	err := gm.utilDB.Exec(createDbSQL).Error
	if err != nil {
		fmt.Println("创建失败：" + err.Error() + " sql:" + createDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库创建成功")
}

func (gm *GormDb) DropDB() {
	dropDbSQL := "DROP DATABASE IF EXISTS " + gm.dbConfig.DbName + ";"
	err := gm.utilDB.Exec(dropDbSQL).Error
	if err != nil {
		fmt.Println("删除失败：" + err.Error() + " sql:" + dropDbSQL)
		return
	}
	fmt.Println(gm.dbConfig.DbName + "数据库删除成功")
}

func (gm *GormDb) GetUtilDB() *gorm.DB {
	return gm.utilDB
}
