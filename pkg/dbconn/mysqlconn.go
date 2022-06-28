/*
 * @Autor: Bobby
 * @Description: function for create and connect mysql pool
 * @Date: 2022-06-06 17:00:19
 * @LastEditTime: 2022-06-09 21:24:11
 * @FilePath: \user\pkg\dbconn\mysqlconn.go
 */
package dbconn

import (
	"time"
	"user/pkg/config"
	"user/pkg/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MysqlConf struct {
	DriveName      string
	DataSourceName string
	MaxIdle        int // Set the maximum number of connections in the free connection pool
	MaxOpen        int // Set the maximum number of open database connections
	MaxLifetime    int // Set the maximum time that the connection can be reused
}

var db = new(gorm.DB)

func GetMysql() *gorm.DB {
	return db
}

func InitMysql() {
	var mysqlConf MysqlConf
	mysqlConf.DriveName = config.Get("config.mysql.driveName").(string)
	mysqlConf.DataSourceName = config.Get("config.mysql.dataSourceName").(string)
	mysqlConf.MaxIdle = config.Get("config.mysql.maxIdle").(int)
	mysqlConf.MaxOpen = config.Get("config.mysql.maxOpen").(int)
	mysqlConf.MaxLifetime = config.Get("config.mysql.maxLifetime").(int)

	logger.SugarLogger.Debug("read mysql config...")
	logger.SugarLogger.Debugf("mysqlConf.DriveName : %v", mysqlConf.DriveName)
	logger.SugarLogger.Debugf("mysqlConf.MaxIdle : %v", mysqlConf.MaxIdle)
	logger.SugarLogger.Debugf("mysqlConf.MaxOpen : %v", mysqlConf.MaxOpen)
	logger.SugarLogger.Debugf("mysqlConf.MaxLifetime : %v", mysqlConf.MaxLifetime)

	var err error
	db, err = gorm.Open(mysqlConf.DriveName, mysqlConf.DataSourceName)
	if err != nil {
		panic(err)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(mysqlConf.MaxIdle)
	db.DB().SetMaxOpenConns(mysqlConf.MaxOpen)
	db.DB().SetConnMaxLifetime(time.Duration(mysqlConf.MaxLifetime) * time.Hour)

	gormLogger := &logger.GormLogger{}
	db.LogMode(true)
	db.SetLogger(gormLogger)

	logger.SugarLogger.Debug("mysql init succ...")
}
