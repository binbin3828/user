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
	MaxIdle        int
	MaxOpen        int
	MaxLifetime    int
}

// NewMysql 创建 MySQL 连接，返回 *gorm.DB 实例
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

	db, err := gorm.Open(mysqlConf.DriveName, mysqlConf.DataSourceName)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(mysqlConf.MaxIdle)
	db.DB().SetMaxOpenConns(mysqlConf.MaxOpen)
	db.DB().SetConnMaxLifetime(time.Duration(mysqlConf.MaxLifetime) * time.Hour)

	gormLogger := &logger.GormLogger{Logger: log}
	db.LogMode(true)
	db.SetLogger(gormLogger)

	log.Debug("mysql init succ...")
	return db, nil
}
