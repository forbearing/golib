# Database Migration Package

This package provides database migration functionality with support for multiple database types.

## Supported Databases

- **MySQL**: Full support with temporary database creation
- **PostgreSQL**: Full support with temporary schema creation  
- **SQLite**: Full support with in-memory temporary database

## Usage

```go
import "github.com/forbearing/golib/internal/dbmigrate"

// Configure migration
cfg := &dbmigrate.Config{
    DSN:         "your-database-connection-string",
    Mode:        dbmigrate.ModeSync, // or ModeSafe, ModeAddonly
    DryRun:      false,
    AutoApprove: false,
}

// Run migration with your GORM models
err := dbmigrate.Migrate(cfg, &User{}, &Product{}, &Order{})
if err != nil {
    log.Fatal(err)
}
```

## DSN Examples

### MySQL
```
user:password@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local
```

### PostgreSQL
```
postgres://user:password@localhost:5432/database?sslmode=disable
host=localhost user=username password=password dbname=database port=5432 sslmode=disable
```

### SQLite
```
/path/to/database.db
./database.sqlite3
file:database.db
```

## Migration Modes

- **sync**: Full synchronization - adds, modifies, and removes tables/columns
- **safe**: Safe mode - only adds new tables and columns, never removes
- **addonly**: Add-only mode - only adds new tables, skips existing ones

## Features

- Automatic database type detection from DSN
- Temporary database/schema creation for safe migration planning
- Destructive operation detection and confirmation
- Dry-run mode for preview
- Auto-approval option for CI/CD