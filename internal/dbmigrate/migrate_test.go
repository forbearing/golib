package dbmigrate_test

import (
	"os"
	"testing"

	"github.com/forbearing/golib/internal/dbmigrate"
	"gorm.io/gorm"
)

type TestUser struct {
	Name string
	Age  string

	gorm.Model
}
type TestGroup struct {
	Name string

	gorm.Model
}

func TestDBMigrate(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SKIP_DB_TEST", "true")
	defer os.Unsetenv("SKIP_DB_TEST")

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		config    *dbmigrate.Config
		allModels []any
	}{
		{
			name:   "test1",
			config: dbmigrate.NewConfig(true, false, true, "safe", ""),
			allModels: []any{
				&TestUser{},
				&TestGroup{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbmigrate.Migrate(tt.config, tt.allModels...)
		})
	}
}
