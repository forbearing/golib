package controller

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user := new(model.User)
	users := make([]*model.User, 0)
	categories := make([]*model.Category, 0)

	database.Database[*model.User]().List(&users)    // without context
	database.Database[*model.User](ctx).List(&users) // with custom context.

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
	database.Database[*model.User]().WithAnd().WithQuery(&model.User{
		Name:  util.ValueOf("admin"),
		Email: util.ValueOf("admin@admin.admin"),
	}) // default WithAnd.

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

	// List users without triggering model's hooks.
	// It's useful to break recursive model hook triggers.
	database.Database[*model.User]().WithoutHook().List(&users)

	// List users from cache if exists, otherwise fetch from database.
	database.Database[*model.User]().WithCache().List(&users)

	// List all users but excluding those with names 'root' or 'admin'.
	//
	// SQL equivalent:
	// SELECT * FROM users
	// WHERE name NOT IN ('root', 'admin')
	// LIMIT 1000
	expands := map[string][]any{"name": {"root", "admin"}}
	database.Database[*model.User]().WithExclude(expands).List(&users)

	// List users with pagination, returns users from offset 1 with limit 20.//
	//
	// SQL equivalent:
	// SELECT * FROM users
	// LIMIT 20 OFFSET 1
	database.Database[*model.User]().WithScope(1, 20).List(&users)

	// List categories with expanded relationships for 'Children' and 'Parent'.
	//
	// Expands configuration:
	// - Children: expands child categories
	// - Parent: expands parent category
	//
	// Example JSON response:
	// {
	//   "id": "electronics",
	//   "name": "Electronics",
	//   "parent": {
	//     "id": "root",
	//     "name": "Root"
	//   },
	//   "children": [
	//     {
	//       "id": "phone",
	//       "name": "Phones"
	//     },
	//     {
	//       "id": "computer",
	//       "name": "Computers"
	//     }
	//   ]
	// }// The frontend will get the nested json data from backend.
	database.Database[*model.Category]().WithExpand(new(model.Category).Expands()).List(&categories)

	// Get the first user from database.
	database.Database[*model.User]().First(user)
	// Get the last user from database.
	database.Database[*model.User]().Last(user)
	// Get a random user from database.
	database.Database[*model.User]().Take(user)

	// =====================================================================
	// Why choose database.Database().Create/Delete/Update/List/Get methods:
	// beacause model's hooks only invoke in database.Database.
	// =====================================================================
}
