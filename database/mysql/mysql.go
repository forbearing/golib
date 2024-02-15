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
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Default *gorm.DB

func Init() (err error) {
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

	for _, m := range model.Tables {
		if len(m.GetTableName()) > 0 {
			if err = Default.Table(m.GetTableName()).AutoMigrate(m); err != nil {
				return
			}
		} else {
			if err = Default.AutoMigrate(m); err != nil {
				return
			}
		}
	}
	// create the table records that must be pre-exists before database curds.
	for _, r := range model.Records {
		// FIXME: 如何 preload, 来递归创建表数据
		// for i := range r.Expands {
		// 	DB = DB.Preload(r.Expands[i])
		// }
		if err = Default.Model(r.Table).Save(r.Rows).Error; err != nil {
			return err
		}
	}
	// set default database to 'Default'.
	database.DB = Default

	return
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
// Transaction executes an arbitrary number of commands in fc within a transaction.
// On success the changes are committed; if an error occurs they are rolled back.
func Transaction(fn func(tx *gorm.DB) error) error {
	return Default.Transaction(fn)
}

// Exec executes raw sql without return rows
func Exec(sql string, values any) error {
	return Default.Exec(sql, values).Error
}
