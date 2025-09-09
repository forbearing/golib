package database_test

import (
	"os"
	"testing"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/stretchr/testify/require"
)

// TestUser test user model
type TestUser struct {
	model.Base
	Name     string `json:"name" gorm:"size:100;not null"`
	Email    string `json:"email" gorm:"size:255;uniqueIndex"`
	Age      int    `json:"age"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}

func (u *TestUser) GetTableName() string {
	return "test_users"
}

// TestProduct test product model
type TestProduct struct {
	model.Base
	Name        string  `json:"name" gorm:"size:200;not null"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2)"`
	CategoryID  string  `json:"category_id" gorm:"size:36"`
}

func (p *TestProduct) GetTableName() string {
	return "test_products"
}

// TestCategory test category model
type TestCategory struct {
	model.Base
	Name     string `json:"name" gorm:"size:100;not null;uniqueIndex"`
	ParentID string `json:"parent_id" gorm:"size:36"`
}

func (c *TestCategory) GetTableName() string {
	return "test_categories"
}

func init() {
	os.Setenv(config.LOGGER_DIR, "/tmp/test_database")
	os.Setenv(config.DATABASE_TYPE, string(config.DBSqlite))
	os.Setenv(config.SQLITE_IS_MEMORY, "false")
	os.Setenv(config.SQLITE_PATH, "/tmp/test.db")

	// os.Setenv(config.DATABASE_TYPE, string(config.DBMySQL))
	// os.Setenv(config.MYSQL_DATABASE, "test")
	// os.Setenv(config.MYSQL_USERNAME, "test")
	// os.Setenv(config.MYSQL_PASSWORD, "test")

	model.Register[*TestUser]()
	model.Register[*TestProduct]()
	model.Register[*TestCategory]()

	if err := bootstrap.Bootstrap(); err != nil {
		panic(err)
	}
}

// createTestUsers creates test user data
func createTestUsers(t *testing.T) []*TestUser {
	now := time.Now()
	users := []*TestUser{
		{
			Base: model.Base{
				ID:        "user-001",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Zhang San",
			Email:    "zhangsan@example.com",
			Age:      25,
			IsActive: true,
		},
		{
			Base: model.Base{
				ID:        "user-002",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Li Si",
			Email:    "lisi@example.com",
			Age:      30,
			IsActive: true,
		},
		{
			Base: model.Base{
				ID:        "user-003",
				CreatedAt: &now,
				UpdatedAt: &now,
			},
			Name:     "Wang Wu",
			Email:    "wangwu@example.com",
			Age:      28,
			IsActive: false,
		},
	}

	err := database.Database[*TestUser](nil).Create(users...)
	require.NoError(t, err)

	return users
}
