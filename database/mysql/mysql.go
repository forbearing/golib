package mysql

/*
https://www.cnblogs.com/bfmq/p/17494295.html
https://casbin.org/zh/docs/supported-models
https://blog.csdn.net/LeoForBest/article/details/133610889	Casbin权限管理实现
https://blog.csdn.net/LeoForBest/article/details/133607878	Casbin权限管理实现
*/

import (
	"fmt"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/helper"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Default *gorm.DB
	dbmap   = make(map[string]*gorm.DB)
)

func Init() (err error) {
	if !config.App.MySQLConfig.Enable {
		return
	}

	dsnDefault := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.App.MySQLConfig.Username,
		config.App.MySQLConfig.Password,
		config.App.MySQLConfig.Host,
		config.App.MySQLConfig.Port,
		config.App.MySQLConfig.Database,
		config.App.MySQLConfig.Charset,
	)
	zap.S().Infow("database info", "host", config.App.MySQLConfig.Host, "port", config.App.MySQLConfig.Port, "database", config.App.MySQLConfig.Database)
	if Default, err = gorm.Open(mysql.Open(dsnDefault), &gorm.Config{Logger: logger.Gorm}); err != nil {
		zap.S().Error(err)
		return err
	}

	return helper.InitDatabase(Default, dbmap)
}

func Transaction(fn func(tx *gorm.DB) error) error { return helper.Transaction(Default, fn) }
func Exec(sql string, values any) error            { return helper.Exec(Default, sql, values) }
