package db

import (
	"cs/config"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var dbPool = DbPool{
	DbMap: make(map[string]*gorm.DB, 1),
}

// DbPool DB连接池
type DbPool struct {
	sync.RWMutex
	DbMap map[string]*gorm.DB
}

type dbConfig struct {
	dbConfigName    string
	userName        string
	password        string
	host            string
	port            int
	database        string
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifeTime time.Duration
	charset         string
}

// Default 获取默认gorm实例
func Default() *gorm.DB {
	return Get("mysql", "default")
}

// DefaultE 获取默认gorm实例，未获取到返回错误
func DefaultE() (*gorm.DB, error) {
	return GetE("mysql", "default")
}

// Get 获取gorm实例
func Get(driver, dbConfigName string) *gorm.DB {
	// 读锁
	dbPool.RLock()
	db, ok := dbPool.DbMap[dbConfigName]
	dbPool.RUnlock()
	if ok && db != nil {
		return db
	}

	// 写锁
	dbPool.Lock()
	defer dbPool.Unlock()
	if db, ok := dbPool.DbMap[dbConfigName]; ok && db != nil {
		return db
	}
	db, err := loadDbByConfig(driver, dbConfigName)
	if err != nil {
		// todo log
	}
	return db
}

// GetE 获取gorm实例，未获取到返回错误
func GetE(driver, dbConfigName string) (*gorm.DB, error) {
	// 读锁
	dbPool.RLock()
	db, ok := dbPool.DbMap[dbConfigName]
	dbPool.RUnlock()
	if ok && db != nil {
		return db, nil
	}

	// 写锁
	dbPool.Lock()
	defer dbPool.Unlock()
	if db, ok := dbPool.DbMap[dbConfigName]; ok && db != nil {
		return db, nil
	}
	db, err := loadDbByConfig(driver, dbConfigName)
	if err != nil {
		// todo log
		return nil, err
	}
	return db, err
}

// loadDbByConfig 根据配置文件价值数据库
func loadDbByConfig(driver, dbConfigName string) (*gorm.DB, error) {
	dbConfigs := config.Use("db")
	if dbConfigs == nil {
		return nil, errors.New("lost key 'db' in config file")
	}

	dbConf := dbConfigs.Use(dbConfigName)
	if dbConf == nil {
		return nil, errors.New("lost key " + dbConfigName + " in db config")
	}

	return load(dbConfig{
		dbConfigName:    dbConfigName,
		userName:        dbConf.GetString("username"),
		password:        dbConf.GetString("password"),
		host:            dbConf.GetString("host"),
		port:            dbConf.GetInt("port"),
		database:        dbConf.GetString("database"),
		maxIdleConns:    dbConf.GetInt("maxIdleConns"),
		maxOpenConns:    dbConf.GetInt("maxOpenConns"),
		connMaxLifeTime: dbConf.MustDuration("connMaxLifeTime", 1*time.Hour),
		charset:         dbConf.MustString("charset", "utf8m64,utf8"),
	})
}

// load 加载数据库
func load(conf dbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
		conf.userName, conf.password, conf.host, conf.port, conf.database, conf.charset)
	c := &gorm.Config{
		SkipDefaultTransaction: true, // 单次写入不需要使用事务
		PrepareStmt:            true, // 缓存预编译语句
		//Logger: todo
	}
	db, err := gorm.Open(mysql.Open(dsn), c)
	if err != nil {
		// todo log
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		// todo log
		return nil, err
	}

	sqlDb.SetMaxIdleConns(conf.maxIdleConns)
	sqlDb.SetMaxOpenConns(conf.maxOpenConns)
	sqlDb.SetConnMaxLifetime(conf.connMaxLifeTime)
	// todo log

	dbPool.DbMap[conf.dbConfigName] = db
	return db, nil
}
