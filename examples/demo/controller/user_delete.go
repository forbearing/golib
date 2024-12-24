package controller

import (
	"context"
	"time"

	"demo/model"

	"github.com/forbearing/golib/database"
	"github.com/gin-gonic/gin"
)

func (*user) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user := new(model.User)
	users := make([]*model.User, 0)
	user.SetID("myid")
	if err := database.Database[*model.User]().WithLimit(-1).List(&users); err != nil {
		panic(err)
	}

	// Delete single User record
	// NOTE: User.ID field must be set - records are delete based on ID
	database.Database[*model.User]().Delete(user)
	database.Database[*model.User](ctx).Delete(user)
	// Batch delete user records in database.
	// NOTE: ID field must be set for each user record.
	database.Database[*model.User]().Delete(users...)
	database.Database[*model.User](ctx).Delete(users...)

	// Delete user records without batch size limit.
	database.Database[*model.User]().WithLimit(-1).Delete(users...)
	// Delete user records with batch size of 100.
	database.Database[*model.User]().WithLimit(100).Delete(users...)

	// Delete users without triggering model's hooks.
	// It's useful to break recursive model hook triggers.
	database.Database[*model.User]().WithoutHook().Delete(user)
	database.Database[*model.User]().WithoutHook().Delete(users...)

	// Permanently delete user records from database
	// Default delete only updates `deleted_at` field to current time.
	database.Database[*model.User]().WithPurge().Delete(users...)

	// Permanently delete the user records that `deleted_at` field is not null.
	//
	// SQL equivalent:
	// DELETE FROM users
	// WHERE deleted_at IS NOT NULL
	database.Database[*model.User]().Cleanup()

	// =====================================================================
	// Why choose database.Database().Create/Delete/Update/List/Get methods:
	// beacause model's hooks only invoke in database.Database.
	// =====================================================================
}
