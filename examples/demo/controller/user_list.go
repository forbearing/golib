package controller

import (
	"time"

	"demo/model"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/database/mysql"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type user struct{}

var User = new(user)

func (*user) List(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	users := make([]*model.User, 0)

	// List all users without limit. -1 disable the limit, returns all users.
	database.Database[*model.User]().WithLimit(-1).List(&users)
	// List 100 users.
	database.Database[*model.User]().WithLimit(100).List(&users)
	// List users filtered by 'name', returns all users where name = 'admin', default limit to 1000.
	database.Database[*model.User]().WithQuery(&model.User{Name: util.ValueOf("admin")}).List(&users)
	// List users filtered by 'name', returns all users where name LIKE '%admin%'.
	database.Database[*model.User]().WithQuery(&model.User{Name: util.ValueOf("admin")}, true).List(&users)

	// // List users from another mysql instance and specifiy the table name.
	// // - WithDB: specifies alternate MySQL instance (mysql.MyDB)
	// // - WithTable: uses custom table name "myproject_users"
	// database.Database[*model.User]().WithDB(mysql.MyDB).WithTable("myproject_users").List(&users)

	// List users with raw SQL query condition
	// WithQueryRaw: Uses raw SQL condition "name = ?"
	database.Database[*model.User]().WithQueryRaw("name = ?", "admin").List(&users)

	// List users with multiple field conditions
	//
	// Query conditions (AND logic):
	// - Name = "admin"
	// - Email = "admin@admin.admin"
	// Only users that `name` is 'admin' and `email` is 'admin@admin.admin' will be matched.
	database.Database[*model.User]().WithQuery(&model.User{
		Name:  util.ValueOf("admin"),
		Email: util.ValueOf("admin@admin.admin"),
	})
	// default WithAnd.
	database.Database[*model.User]().WithAnd().WithQuery(&model.User{
		Name:  util.ValueOf("admin"),
		Email: util.ValueOf("admin@admin.admin"),
	})

	// List users with multiple field conditions
	//
	// Query conditions (OR logic):
	// - Name = "admin" OR
	// - Email = "admin@admin.admin"
	//
	// All users that `name` is 'admin' or `email` is 'admin@admin.admin' will be matched.
	database.Database[*model.User]().WithOr().WithQuery(&model.User{
		Name:  util.ValueOf("admin"),
		Email: util.ValueOf("admin@admin.admin"),
	})

	// List users with specific field section.
	//
	// Default fields (always included):
	// - id
	// - created_by
	// - update_by
	// - created_at
	// - updated_at
	//
	// Additional selected fields:
	// - name
	// - email
	//
	// SQL equivalent:
	// SELECT id, created_by, update_by, created_at, updated_at,
	//        name, email
	// FROM users
	// LIMIT 1000
	database.Database[*model.User]().WithSelect("name", "email").List(&users)

	// List users with creation time range condition
	//
	// Time range configuration:
	// - Field: created_at
	// - Begin: next week (now + 7 days)
	// - End: now
	//
	// SQL equivalent:
	// SELECT * FROM users
	// WHERE created_at BETWEEN [now] AND [next week]
	// LIMIT 1000
	now := time.Now()
	begin := time.Now().Add(1 * 24 * 7 * time.Hour)
	database.Database[*model.User]().WithTimeRange("created_at", begin, now).List(&users)

	// List users using specific index
	// It's your responsibility to ensure index `idx_composite_name_email_createdby` already exists.
	database.Database[*model.User]().WithIndex("idx_composite_name_email_createdby").List(&users)

	// List users order by 'created_at' in descending order.
	database.Database[*model.User]().WithOrder("created_at desc").List(&users)

	// List users with multiple field ordering
	//
	// Order configuration (comma separated):
	// 1. created_at DESC - Sort by creation time descending
	// 2. name ASC - Then sort by name ascending
	//
	// SQL equivalent:
	// SELECT * FROM users
	// ORDER BY created_at DESC, name ASC
	// LIMIT 1000
	database.Database[*model.User]().WithOrder("created_at desc, name asc").List(&users)

	// Execute operations in a database transaction
	mysql.Default.Transaction(func(tx *gorm.DB) (err error) {
		all := make([]*model.User, 0)
		if err = database.Database[*model.User]().WithTransaction(tx).WithLock().WithLimit(-1).List(&all); err != nil {
			return err
		}
		if err = database.Database[*model.User]().WithTransaction(tx).WithLock().WithLimit(-1).WithPurge().Delete(all...); err != nil {
			return err
		}
		if err = database.Database[*model.User]().WithTransaction(tx).WithLock().WithLimit(-1).Update(users...); err != nil {
			return err
		}
		return nil
	})

	// List users without triggering model's hooks
	// its usefull to break recurive trigger model's hooks.
	database.Database[*model.User]().WithoutHook().List(&users)

	// =====================================================================
	// Why choose database.Database().Create/Delete/Update/List/Get methods:
	// beacause model's hooks only invoke in database.Database.
	// =====================================================================
}

func (*user) Create(c *gin.Context) {
	database.Database[*model.User]().WithTryRun()
}

func (*user) Delete(c *gin.Context) {
	database.Database[*model.User]().WithPurge()
}
