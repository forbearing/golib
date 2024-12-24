package controller

import (
	"context"
	"time"

	"demo/model"

	"github.com/forbearing/golib/database"
	"github.com/gin-gonic/gin"
)

func (*user) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user := new(model.User)
	users := make([]*model.User, 0)

	// Create single User record
	// NOTE: ID field can be empty, database will auto-generate unique ID
	database.Database[*model.User]().Create(user)
	database.Database[*model.User](ctx).Create(user)

	// Batch create user records in database.
	// NOTE: ID fields can be empty, database will auto-generate unique IDs for each record
	database.Database[*model.User]().Create(users...)
	database.Database[*model.User](ctx).Create(users...)

	// Create user records without batch size limit.
	database.Database[*model.User]().WithLimit(-1).Create(users...)
	// Create user records with batch size of 100.
	// It's useful when the User model contains so many fields.
	database.Database[*model.User]().WithLimit(100).Create(users...)

	// Create users without triggering model's hooks.
	// It's useful to break recursive model hook triggers.
	database.Database[*model.User]().WithoutHook().Create(user)
	database.Database[*model.User]().WithoutHook().Create(users...)

	// Create users with all fields but except 'email' and 'avatar'.
	// result: created users will have empty values for 'email' and 'avatar', other fields are saved as povided.
	database.Database[*model.User]().WithOmit("email", "avatar").Create(user)
	database.Database[*model.User]().WithOmit("email", "avatar").Create(users...)

	// =====================================================================
	// Why choose database.Database().Create/Delete/Update/List/Get methods:
	// beacause model's hooks only invoke in database.Database.
	// =====================================================================
}
