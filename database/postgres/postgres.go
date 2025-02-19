package postgres

import (
	"database/sql"
	"fmt"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/helper"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Default *gorm.DB
	db      *sql.DB
	dbmap   = make(map[string]*gorm.DB)
)

// Init initializes the default PostgreSQL connection.
// It checks if PostgreSQL is enabled and selected as the default database.
// If the connection is successful, it initializes the database and returns nil.
func Init() (err error) {
	cfg := config.App.PostgreConfig
	if !cfg.Enable || config.App.ServerConfig.DB != config.DBPostgre {
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

	zap.S().Infow("successfully connect to postgres", "host", cfg.Host, "port", cfg.Port, "database", cfg.Database, "sslmode", cfg.SSLMode, "timezone", cfg.TimeZone)
	return helper.InitDatabase(Default, dbmap)
}

// New creates and returns a new PostgreSQL database connection with the given configuration.
// Returns (*gorm.DB, error) where error is non-nil if the connection fails.
func New(cfg config.PostgreConfig) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(buildDSN(cfg)), &gorm.Config{Logger: logger.Gorm})
}

func buildDSN(cfg config.PostgreConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Username, cfg.Password, cfg.Database, cfg.Port, cfg.SSLMode, cfg.TimeZone,
	)
}

func Transaction(fn func(tx *gorm.DB) error) error { return helper.Transaction(Default, fn) }
func Exec(sql string, values any) error            { return helper.Exec(Default, sql, values) }
