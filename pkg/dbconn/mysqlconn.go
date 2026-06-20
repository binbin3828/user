package dbconn

import (
	"fmt"
	"time"
	"user/pkg/config"
	"user/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type MysqlConf struct {
	DriveName      string
	DataSourceName string
	MaxIdle        int
	MaxOpen        int
	MaxLifetime    int
}

func NewMysql(log logger.Logger) (*gorm.DB, error) {
	var mysqlConf MysqlConf
	mysqlConf.DriveName = config.Get("config.mysql.driveName").(string)
	mysqlConf.DataSourceName = config.Get("config.mysql.dataSourceName").(string)
	mysqlConf.MaxIdle = config.Get("config.mysql.maxIdle").(int)
	mysqlConf.MaxOpen = config.Get("config.mysql.maxOpen").(int)
	mysqlConf.MaxLifetime = config.Get("config.mysql.maxLifetime").(int)

	log.Debug("read mysql config...")
	log.Debugf("mysqlConf.DriveName : %v", mysqlConf.DriveName)
	log.Debugf("mysqlConf.MaxIdle : %v", mysqlConf.MaxIdle)
	log.Debugf("mysqlConf.MaxOpen : %v", mysqlConf.MaxOpen)
	log.Debugf("mysqlConf.MaxLifetime : %v", mysqlConf.MaxLifetime)

	gormLogger := &logger.GormLogger{Logger: log}

	var db *gorm.DB
	var err error

	for i := 0; i < 3; i++ {
		db, err = gorm.Open(mysql.Open(mysqlConf.DataSourceName), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
			Logger:         gormLogger,
		})
		if err == nil {
			break
		}
		log.Warnf("mysql connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("mysql connection failed after 3 attempts: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdle)
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlConf.MaxLifetime) * time.Hour)

	log.Debug("mysql init succ...")
	return db, nil
}
