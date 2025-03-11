package controller

import (
	"demo/model"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/types"
	"github.com/gin-gonic/gin"
)

func (u *user) Update(c *gin.Context) {
	ctx := new(types.DatabaseContext)

	user := new(model.User)
	users := make([]*model.User, 0)

	// Update single User record
	// NOTE: User.ID field must be set - records are update based on ID
	// All fields will be updated.
	database.Database[*model.User]().Update(user)
	database.Database[*model.User](ctx).Update(user)
	// Batch update user records in database.
	// NOTE: ID field must be set for each user record.
	// All fields will be updated.
	database.Database[*model.User]().Update(users...)
	database.Database[*model.User](ctx).Update(users...)

	// Update user records without batch size limit.
	database.Database[*model.User]().WithLimit(-1).Update(users...)
	// Update user records with batch size of 100.
	// It's useful when the User model contains so many fields.
	database.Database[*model.User]().WithLimit(100).Update(users...)

	// Update users without triggering model's hooks.
	// It's useful to break recursive model hook triggers.
	database.Database[*model.User]().WithoutHook().Update(user)
	database.Database[*model.User]().WithoutHook().Update(users...)

	// Update users with all fields but except 'email' and 'avatar'.
	database.Database[*model.User]().WithOmit("email", "avatar").Update(user)
	database.Database[*model.User]().WithOmit("email", "avatar").Update(users...)

	// Update Single User record.
	// Only update `name` field.
	// NOTE: UpdateById will not trigger model's hooks.
	//
	// SQL equivalent:
	// UPDATE users
	// SET name = 'new_name'
	// WHERE id = 'my_user_id'
	// AND name = 'old_name'
	database.Database[*model.User]().UpdateById("my_user_id", "old_name", "new_name")
	database.Database[*model.User](ctx).UpdateById("my_user_id", "old_name", "new_name")

	// =====================================================================
	// Why choose database.Database().Create/Delete/Update/List/Get methods:
	// beacause model's hooks only invoke in database.Database.
	// =====================================================================
}
