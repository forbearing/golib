package sqlite

import (
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database/helper"
	"github.com/forbearing/golib/logger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Default *gorm.DB
	dbmap   = make(map[string]*gorm.DB)
)

func Init() (err error) {
	if !config.App.SqliteConfig.Enable {
		return
	}

	dsn := config.App.SqliteConfig.Path
	if config.App.SqliteConfig.IsMemory {
		dsn = "file::memory:?cache=shared" // Ignore file based database if IsMemory is true.
	}
	zap.S().Infow("database info", "path", config.App.SqliteConfig.Path, "database", config.App.SqliteConfig.Database, "is_memory", config.App.SqliteConfig.IsMemory)
	if Default, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{Logger: logger.Gorm}); err != nil {
		zap.S().Error(err)
		return err
	}
	return helper.InitDatabase(Default, dbmap)
}

func Transaction(fn func(tx *gorm.DB) error) error { return helper.Transaction(Default, fn) }
func Exec(sql string, values any) error            { return helper.Exec(Default, sql, values) }
