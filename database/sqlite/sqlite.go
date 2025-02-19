package sqlite

import (
	"database/sql"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/helper"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Default *gorm.DB
	db      *sql.DB
	dbmap   = make(map[string]*gorm.DB)
)

// Init initializes the default SQLite connection.
// It checks if SQLite is enabled and selected as the default database.
// If the connection is successful, it initializes the database and returns nil.
func Init() (err error) {
	cfg := config.App.SqliteConfig
	if !cfg.Enable || config.App.ServerConfig.DB != config.DBSqlite {
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

	zap.S().Infow("successfully connect to sqlite", "path", cfg.Path, "database", cfg.Database, "is_memory", cfg.IsMemory)
	return helper.InitDatabase(Default, dbmap)
}

// New creates and returns a new SQLite database connection with the given configuration.
// Returns (*gorm.DB, error) where error is non-nil if the connection fails.
func New(cfg config.SqliteConfig) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(buildDSN(cfg)), &gorm.Config{Logger: logger.Gorm})
}

func buildDSN(cfg config.SqliteConfig) string {
	dsn := cfg.Path
	if cfg.IsMemory || len(cfg.Path) == 0 {
		if len(cfg.Path) == 0 {
			zap.S().Warn("sqlite path is empty, using in-memory database")
		}
		dsn = "file::memory:?cache=shared" // Ignore file based database if IsMemory is true.
	}
	return dsn
}

func Transaction(fn func(tx *gorm.DB) error) error { return helper.Transaction(Default, fn) }
func Exec(sql string, values any) error            { return helper.Exec(Default, sql, values) }
