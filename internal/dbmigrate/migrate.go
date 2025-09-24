package dbmigrate

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	// Database drivers
	"github.com/forbearing/golib/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	// Atlas packages
	"ariga.io/atlas-provider-gorm/gormschema"
	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
)

type Mode string

const (
	ModeSync    = "sync"
	ModeSafe    = "safe"
	ModeAddonly = "addonly"
)

// Config migration configuration
type Config struct {
	DryRun      bool   // only show SQL, don't execute
	AutoApprove bool   // auto approve execution
	Mode        Mode   // sync: full sync, safe: only add, addonly: only add new tables and fields
	DSN         string // database connection string (supports MySQL, PostgreSQL, SQLite)
	DBType      config.DBType
}

func Migrate(cfg *Config, allModels ...any) error {
	ctx := context.Background()

	// check
	if cfg == nil {
		return errors.New("config is nil")
	}
	switch cfg.Mode {
	case ModeSync, ModeSafe, ModeAddonly:
	default:
		return fmt.Errorf("invalid migration mode: %s, expected: sync, safe, addonly", cfg.Mode)
	}
	if len(cfg.DSN) == 0 {
		return fmt.Errorf("database connection string is empty")
	}

	switch cfg.DBType {
	case config.DBMySQL, config.DBSqlite, config.DBPostgres:
	default:
		return fmt.Errorf("unsupported database type: %s, expected: (%s|%s|%s)", cfg.DBType, config.DBMySQL, config.DBSqlite, config.DBPostgres)
	}

	// Connect to database
	db, err := sql.Open(string(cfg.DBType), cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Test connection
	if err = db.Ping(); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}

	// Create appropriate Atlas driver
	var drv migrate.Driver
	switch cfg.DBType {
	case config.DBMySQL:
		drv, err = mysql.Open(db)
	case config.DBPostgres:
		drv, err = postgres.Open(db)
	case config.DBSqlite:
		drv, err = sqlite.Open(db)
	default:
		return fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}
	if err != nil {
		return fmt.Errorf("failed to open database driver: %w", err)
	}

	// 2. Create temporary database and connection
	var tempConn *sql.DB
	var tempDBName string
	var tempDriver migrate.Driver

	switch cfg.DBType {
	case config.DBMySQL:
		// Create temporary database for MySQL
		tempDBName = fmt.Sprintf("temp_migrate_%d", time.Now().Unix())
		if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", tempDBName)); err != nil {
			return fmt.Errorf("failed to create temporary database: %w", err)
		}
		defer func() {
			if _, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", tempDBName)); err != nil {
				logWarning(fmt.Sprintf("failed to clean up temporary database: %v", err))
			}
		}()

		// Connect to temporary MySQL database
		tempDSN := fmt.Sprintf("nebula:nebula@tcp(localhost:3307)/%s?charset=utf8mb4&parseTime=True&loc=Local", tempDBName)
		tempConn, err = sql.Open("mysql", tempDSN)
		if err != nil {
			return fmt.Errorf("failed to connect to temporary database: %w", err)
		}
		defer tempConn.Close()

		tempDriver, err = mysql.Open(tempConn)

	case config.DBPostgres:
		// For PostgreSQL, use a temporary schema instead of database
		tempDBName = fmt.Sprintf("temp_migrate_%d", time.Now().Unix())
		if _, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", tempDBName)); err != nil {
			return fmt.Errorf("failed to create temporary schema: %w", err)
		}
		defer func() {
			if _, err = db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", tempDBName)); err != nil {
				logWarning(fmt.Sprintf("failed to clean up temporary schema: %v", err))
			}
		}()

		tempConn = db // Use same connection for PostgreSQL
		tempDriver, err = postgres.Open(tempConn)

	case config.DBSqlite:
		// For SQLite, use in-memory database
		tempConn, err = sql.Open("sqlite3", ":memory:")
		if err != nil {
			return fmt.Errorf("failed to create temporary SQLite database: %w", err)
		}
		defer tempConn.Close()

		tempDBName = "main" // SQLite default schema name
		tempDriver, err = sqlite.Open(tempConn)
	}

	if err != nil {
		return fmt.Errorf("failed to create temporary database driver: %w", err)
	}

	// Get DDL and execute to temporary database
	loader := gormschema.New(string(cfg.DBType))
	ddlStmts, err := loader.Load(allModels...)
	if err != nil {
		return fmt.Errorf("failed to generate DDL from models: %w", err)
	}

	// Split DDL statements and execute one by one
	statements := strings.SplitSeq(ddlStmts, ";")
	for stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// For PostgreSQL, prefix table names with schema
		if cfg.DBType == config.DBPostgres && !strings.Contains(strings.ToLower(stmt), "create schema") {
			stmt = strings.ReplaceAll(stmt, "CREATE TABLE ", fmt.Sprintf("CREATE TABLE %s.", tempDBName))
		}

		_, err = tempConn.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute DDL to temporary database: %w, statement: %s", err, stmt)
		}
	}

	// Inspect temporary database schema as target schema
	desired, err := tempDriver.InspectSchema(ctx, tempDBName, nil)
	if err != nil {
		return fmt.Errorf("failed to inspect temporary database schema: %w", err)
	}

	// Get current database schema name
	var currentSchemaName string
	switch cfg.DBType {
	case config.DBMySQL:
		currentSchemaName = "nebula" // Hardcoded for now, could be extracted from DSN
	case config.DBPostgres:
		currentSchemaName = "public" // Default PostgreSQL schema
	case config.DBSqlite:
		currentSchemaName = "main" // SQLite default schema
	}

	// 修正 schema 名称以匹配目标数据库
	desired.Name = currentSchemaName

	// 4. Get tables to inspect
	var tablesToInspect []string
	// Get all table names from desired schema
	for _, t := range desired.Tables {
		tablesToInspect = append(tablesToInspect, t.Name)
	}

	// 5. Inspect current database schema (only check specified tables)
	inspectOpts := &schema.InspectOptions{
		Tables: tablesToInspect,
	}
	current, err := drv.InspectSchema(ctx, currentSchemaName, inspectOpts)
	if err != nil {
		return fmt.Errorf("failed to inspect current database structure: %w", err)
	}

	// 6. Filter target schema based on configuration
	filteredDesired := filterDesiredSchema(desired, cfg, current)

	// 7. Calculate schema differences
	changes, err := drv.SchemaDiff(current, filteredDesired)
	if err != nil {
		return fmt.Errorf("failed to calculate schema differences: %w", err)
	}

	if len(changes) == 0 {
		return nil // Database is up to date, no migration needed
	}

	// 8. Generate SQL migration plan
	plan, err := drv.PlanChanges(ctx, "sync", changes)
	if err != nil {
		return fmt.Errorf("failed to generate migration plan: %w", err)
	}

	// Display operations to be executed
	fmt.Println()
	logSeparator("=", 60)
	logTitle("Migration operations to be executed", len(plan.Changes))
	logSeparator("=", 60)

	destructive := false
	for i, c := range plan.Changes {
		logOperation(i+1, fmt.Sprintf("%s;", c.Cmd))

		// Detect destructive operations
		cmdLower := strings.ToLower(c.Cmd)
		if strings.Contains(cmdLower, "drop") ||
			strings.Contains(cmdLower, "delete") ||
			strings.Contains(cmdLower, "alter") && strings.Contains(cmdLower, "modify") {
			destructive = true
			// fmt.Printf("    Warning: This operation may cause data loss!\n")
			logWarning("Warning: This operation may cause data loss!")
		}
	}
	logSeparator("=", 60)

	// If dry-run mode, stop here
	if cfg.DryRun {
		fmt.Println("Dry-run mode: SQL above will not be executed")
		if destructive {
			fmt.Println("⚠ Destructive operations detected, please execute with caution!")
		}
		return nil
	}

	// If destructive operations exist, require confirmation
	if destructive && !cfg.AutoApprove {
		if !confirm("Are you sure you want to continue? This may cause data loss!") {
			logWarning("Migration cancelled by user")
			return nil
		}
	} else if !cfg.AutoApprove {
		if !confirm("Are you sure you want to execute the migration operations above?") {
			logWarning("Migration cancelled by user")
			return nil
		}
	}

	// 9. Execute migration
	if err := drv.ApplyChanges(ctx, changes); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

// filterDesiredSchema 根据配置过滤目标schema
func filterDesiredSchema(desired *schema.Schema, config *Config, current *schema.Schema) *schema.Schema {
	if config.Mode == "sync" {
		return desired // 完全同步模式，不过滤
	}

	filtered := &schema.Schema{
		Name:   desired.Name,
		Tables: []*schema.Table{},
		Realm:  desired.Realm,
	}

	currentTableMap := make(map[string]*schema.Table)
	for _, t := range current.Tables {
		currentTableMap[t.Name] = t
	}

	for _, desiredTable := range desired.Tables {
		currentTable, exists := currentTableMap[desiredTable.Name]

		if config.Mode == "addonly" && exists {
			// addonly mode: skip existing tables
			filtered.Tables = append(filtered.Tables, currentTable)
			continue
		}

		if !exists {
			// New table, add it
			filtered.Tables = append(filtered.Tables, desiredTable)
			continue
		}

		// safe mode or existing table case
		if config.Mode == "safe" {
			// Create a new table with only new fields added
			filteredTable := &schema.Table{
				Name:    desiredTable.Name,
				Schema:  desiredTable.Schema,
				Columns: []*schema.Column{},
				Indexes: []*schema.Index{},
			}

			// Keep all existing fields
			currentColMap := make(map[string]bool)
			for _, col := range currentTable.Columns {
				currentColMap[col.Name] = true
				filteredTable.Columns = append(filteredTable.Columns, col)
			}

			// Only add new fields
			for _, col := range desiredTable.Columns {
				if !currentColMap[col.Name] {
					filteredTable.Columns = append(filteredTable.Columns, col)
				}
			}

			// Handle indexes (only add new indexes)
			currentIdxMap := make(map[string]bool)
			for _, idx := range currentTable.Indexes {
				currentIdxMap[idx.Name] = true
				filteredTable.Indexes = append(filteredTable.Indexes, idx)
			}

			for _, idx := range desiredTable.Indexes {
				if !currentIdxMap[idx.Name] {
					filteredTable.Indexes = append(filteredTable.Indexes, idx)
				}
			}

			// Keep primary key settings
			if desiredTable.PrimaryKey != nil {
				filteredTable.PrimaryKey = desiredTable.PrimaryKey
			} else if currentTable.PrimaryKey != nil {
				filteredTable.PrimaryKey = currentTable.PrimaryKey
			}

			filtered.Tables = append(filtered.Tables, filteredTable)
		}
	}

	return filtered
}

// confirm requests user confirmation
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	logPrompt(prompt)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
