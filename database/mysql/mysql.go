package mysql

import (
	"database/sql"
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
	db      *sql.DB
	dbmap   = make(map[string]*gorm.DB)
)

// Init initializes the default MySQL connection.
// It checks if MySQL is enabled and selected as the default database.
// If the connection is successful, it initializes the database and returns nil.
func Init() (err error) {
	cfg := config.App.MySQLConfig
	if !cfg.Enable || config.App.ServerConfig.DB != config.DBMySQL {
		return
	}

	if Default, err = New(cfg); err != nil {
		zap.S().Error(err)
		return err
	}
	if db, err = Default.DB(); err != nil {
		zap.S().Error(err)
		return err
	}
	db.SetMaxIdleConns(config.App.DatabaseConfig.MaxIdleConns)
	db.SetMaxOpenConns(config.App.DatabaseConfig.MaxOpenConns)
	db.SetConnMaxLifetime(config.App.DatabaseConfig.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.App.DatabaseConfig.ConnMaxIdleTime)

	zap.S().Infow("successfully connect to mysql", "host", cfg.Host, "port", cfg.Port, "database", cfg.Database)
	return helper.InitDatabase(Default, dbmap)
}

// New creates and returns a new MySQL database connection with the given configuration.
// Returns (*gorm.DB, error) where error is non-nil if the connection fails.
func New(cfg config.MySQLConfig) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(buildDSN(cfg)), &gorm.Config{Logger: logger.Gorm})
}

func buildDSN(cfg config.MySQLConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset,
	)
}

func Transaction(fn func(tx *gorm.DB) error) error { return helper.Transaction(Default, fn) }
func Exec(sql string, values any) error            { return helper.Exec(Default, sql, values) }
