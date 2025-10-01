package helper

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/gst/database"
	"github.com/forbearing/gst/model"
	"github.com/forbearing/gst/util"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// InitDatabase will create the table and table records that predefined in model package.
// NOTE:The version of gorm.io/driver/postgres lower than v1.5.4 have some issues.
// more details see: https://github.com/go-gorm/gorm/issues/6886
func InitDatabase(db *gorm.DB, dbmap map[string]*gorm.DB) (err error) {
	// Install GORM OpenTelemetry tracing plugin
	if err = db.Use(otelgorm.NewPlugin()); err != nil {
		zap.S().Warnw("failed to install GORM OpenTelemetry tracing plugin", "error", err)
	}

	// Install tracing plugin for custom databases
	for _, customDB := range dbmap {
		if err = customDB.Use(otelgorm.NewPlugin()); err != nil {
			zap.S().Warnw("failed to install GORM OpenTelemetry tracing plugin for custom DB", "error", err)
		}
	}

	begin := time.Now()
	// create table automatically in default database.
	for _, m := range model.Tables {
		typ := reflect.TypeOf(m).Elem()
		if err = db.Table(m.GetTableName()).AutoMigrate(m); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to create table(%s)", typ.String()))
		}
	}
	// create table automatically with custom database.
	for _, v := range model.TablesWithDB {
		handler := db
		if val, exists := dbmap[strings.ToLower(v.DBName)]; exists {
			handler = val
		}
		m := v.Table
		typ := reflect.TypeOf(m).Elem()
		if err = handler.Table(m.GetTableName()).AutoMigrate(m); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to create table(%s)", typ.String()))
		}
	}
	zap.S().Infow("database create table", "cost", util.FormatDurationSmart(time.Since(begin)))

	begin = time.Now()
	// create the table records that must be pre-exists before database curds.
	for _, r := range model.Records {
		handler := db
		if val, exists := dbmap[strings.ToLower(r.DBName)]; exists {
			handler = val
		}
		if err = handler.Table(r.Table.GetTableName()).Save(r.Rows).Error; err != nil {
			return errors.Wrap(err, "failed to create table records")
		}
	}
	zap.S().Infow("database create table records", "cost", util.FormatDurationSmart(time.Since(begin)))

	// set default database to 'Default'.
	database.DB = db

	return nil
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
// Transaction executes an arbitrary number of commands in fc within a transaction.
// On success the changes are committed; if an error occurs they are rolled back.
func Transaction(db *gorm.DB, fn func(tx *gorm.DB) error) error { return db.Transaction(fn) }

// Exec executes raw sql without return rows
func Exec(db *gorm.DB, sql string, values any) error { return db.Exec(sql, values).Error }
