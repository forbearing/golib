package main

import (
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/internal/dbmigrate"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	dryRun      bool
	autoApprove bool
	backup      bool
	mode        string
	tables      string
	dsn         string
	db          string
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "database migration based on GORM models",
	Long: `Database migration tool that synchronizes database schema with GORM models.

Supported modes:
  sync     - Full synchronization (add/modify/delete)
  safe     - Only add new tables and columns (default)
  addonly  - Only add new tables, skip existing tables

Examples:
  gg migrate --dry-run                    # Show SQL without executing
  gg migrate --mode=sync --auto-approve   # Full sync with auto approval
  gg migrate --tables=users,groups        # Migrate specific tables only
  gg migrate --dsn="user:pass@tcp(host:port)/db"  # Custom database connection`,
	Run: func(cmd *cobra.Command, args []string) {
		migrateRun()
	},
}

func init() {
	migrateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show SQL without executing")
	migrateCmd.Flags().BoolVar(&autoApprove, "auto-approve", false, "auto approve without confirmation")
	migrateCmd.Flags().BoolVar(&backup, "backup", true, "backup before destructive operations")
	migrateCmd.Flags().StringVar(&mode, "mode", "safe", "migration mode: sync, safe, addonly")
	migrateCmd.Flags().StringVar(&tables, "tables", "", "specific tables to migrate (comma separated)")
	migrateCmd.Flags().StringVar(&dsn, "dsn", "", "database connection string")
	migrateCmd.Flags().StringVar(&db, "db", "", "database type")
}

type TestUser struct {
	Name string
	Age  string

	gorm.Model
}
type TestGroup struct {
	Name string

	gorm.Model
}

func migrateRun() {
	const dsn = "nebula:nebula@tcp(localhost:3307)/nebula?charset=utf8mb4&parseTime=True&loc=Local"

	logSection("Starting database migration...")

	err := dbmigrate.Migrate(&dbmigrate.Config{
		DryRun:      dryRun,
		AutoApprove: autoApprove,
		Mode:        dbmigrate.Mode(mode),
		DSN:         dsn,
		DBType:      config.DBType(db),
	}, &TestUser{}, &TestGroup{})
	if err != nil {
		checkErr(err)
	}

	logSuccess("Migration completed successfully")
}
