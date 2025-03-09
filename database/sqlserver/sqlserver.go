package sqlserver

import (
	"database/sql"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/helper"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	Default *gorm.DB
	db      *sql.DB
	dbmap   = make(map[string]*gorm.DB)
)

// Init initializes the default SQLServer connection.
// It checks if SQLServer is enabled and selected as the default database.
// If the connection is successful, it initializes the database and returns nil.
func Init() (err error) {
	cfg := config.App.SQLServer
	if !cfg.Enable || config.App.Server.DB != config.DBSQLServer {
		return
	}

	if Default, err = New(cfg); err != nil {
		return errors.Wrap(err, "failed to connect to sqlserver")
	}
	if db, err = Default.DB(); err != nil {
		return errors.Wrap(err, "failed to get sqlserver db")
	}
	db.SetMaxIdleConns(config.App.Database.MaxIdleConns)
	db.SetMaxOpenConns(config.App.Database.MaxOpenConns)
	db.SetConnMaxLifetime(config.App.Database.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.App.Database.ConnMaxIdleTime)

	zap.S().Infow("successfully connect to sqlserver", "host", cfg.Host, "port", cfg.Port, "database", cfg.Database)
	return helper.InitDatabase(Default, dbmap)
}

// New creates and returns a new SQLServer database connection with the given configuration.
// Returns (*gorm.DB, error) where error is non-nil if the connection fails.
func New(cfg config.SQLServer) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(buildDSN(cfg)), &gorm.Config{Logger: logger.Gorm})
}

func buildDSN(cfg config.SQLServer) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=%v&trustServerCertificate=%v",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
		cfg.Encrypt, cfg.TrustServer,
	)
}

func Transaction(fn func(tx *gorm.DB) error) error { return helper.Transaction(Default, fn) }
func Exec(sql string, values any) error            { return helper.Exec(Default, sql, values) }
