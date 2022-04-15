package db

import (
	"cs/config"
	"cs/constants"
	"cs/log"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	driver 			string
}

// Default 获取默认gorm实例
func Default() *gorm.DB {
	return Get(constants.MySQLDriver, "default")
}

// DefaultE 获取默认gorm实例，未获取到返回错误
func DefaultE() (*gorm.DB, error) {
	return GetE(constants.MySQLDriver, "default")
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
		log.Errorf("load datebase error: %v", err)
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
		log.Errorf("load database error: %v", err)
		return nil, err
	}
	return db, err
}

// loadDbByConfig 根据配置文件加载数据库
func loadDbByConfig(driver, dbConfigName string) (*gorm.DB, error) {
	dbConfigs := config.Use("db")
	if dbConfigs == nil {
		return nil, errors.New("lost key 'db' in config file")
	}

	dbConf := dbConfigs.Use(dbConfigName)
	if dbConf == nil {
		return nil, errors.New("lost key " + dbConfigName + " in db config")
	}
	disableLog := dbConf.MustBool("disableLog", false)
	log.Infof("=%v", disableLog)

	return load(dbConfig{
		dbConfigName:    dbConfigName,
		userName:        dbConf.GetString("username"),
		password:        dbConf.GetString("password"),
		host:            dbConf.GetString("host"),
		port:            dbConf.GetInt("port"),
		database:        dbConf.GetString("database"),
		maxIdleConns:    dbConf.MustInt("maxIdleConns", 4),
		maxOpenConns:    dbConf.MustInt("maxOpenConns", 32),
		connMaxLifeTime: dbConf.MustDuration("connMaxLifeTime", 1*time.Hour),
		charset:         dbConf.MustString("charset", "utf8m64,utf8"),
		driver: 		 driver,
	}, disableLog)
}

// load 加载数据库
func load(conf dbConfig, disableLog bool) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=3s",
		conf.userName, conf.password, conf.host, conf.port, conf.database, conf.charset)
	c := &gorm.Config{
		SkipDefaultTransaction: true, // 单次写入不需要使用事务
		PrepareStmt:            true, // 缓存预编译语句
		Logger: 				log.NewGormLogger(3 * time.Second),
	}
	if disableLog {
		c.Logger = logger.Default.LogMode(logger.Silent)
	}
	var db *gorm.DB
	var err error

	switch conf.driver {
	case constants.MySQLDriver:
		db, err = gorm.Open(mysql.Open(dsn), c)
	case constants.PostgreSQLDriver:
		db, err = gorm.Open(postgres.Open(dsn), c)
	case constants.SQLiteDriver:
		db, err = gorm.Open(sqlite.Open(dsn), c)
	case constants.SQLServerDriver:
		db, err = gorm.Open(sqlserver.Open(dsn), c)
	default:
		db, err = gorm.Open(mysql.Open(dsn), c)
	}
	if err != nil {
		log.Errorf("failed to connect database, config name: %s, err: %v", conf.dbConfigName, err)
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		log.Errorf("failed to get *sql.DB: %v", err)
		return nil, err
	}

	sqlDb.SetMaxIdleConns(conf.maxIdleConns)
	sqlDb.SetMaxOpenConns(conf.maxOpenConns)
	sqlDb.SetConnMaxLifetime(conf.connMaxLifeTime)
	log.Infof("Load database configName=%s, host=%s, port=%d, database=%s, maxIdleConns=%d, maxOpenConns=%d",
		conf.dbConfigName, conf.host, conf.port, conf.database, conf.maxIdleConns, conf.maxOpenConns)

	dbPool.DbMap[conf.dbConfigName] = db
	return db, nil
}
