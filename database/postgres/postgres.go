package postgres

import (
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
	dbmap   = make(map[string]*gorm.DB)
)

func Init() (err error) {
	if !config.App.PostgreConfig.Enable || config.App.ServerConfig.DB != config.DBPostgre {
		return
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.App.PostgreConfig.Host,
		config.App.PostgreConfig.Username,
		config.App.PostgreConfig.Password,
		config.App.PostgreConfig.Database,
		config.App.PostgreConfig.Port,
		config.App.PostgreConfig.SSLMode,
		config.App.PostgreConfig.TimeZone,
	)
	if Default, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Gorm}); err != nil {
		zap.S().Error(err)
		return err
	}
	zap.S().Infow("successfully connect to postgres",
		"host", config.App.PostgreConfig.Host,
		"port", config.App.PostgreConfig.Port,
		"database", config.App.PostgreConfig.Database,
		"username", config.App.PostgreConfig.Username,
		"sslmode", config.App.PostgreConfig.SSLMode,
		"timezone", config.App.PostgreConfig.TimeZone,
	)
	return helper.InitDatabase(Default, dbmap)
}

func Transaction(fn func(tx *gorm.DB) error) error { return helper.Transaction(Default, fn) }
func Exec(sql string, values any) error            { return helper.Exec(Default, sql, values) }
