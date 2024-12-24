package controller

import (
	"context"
	"time"

	"demo/model"

	"github.com/forbearing/golib/database"
	"github.com/gin-gonic/gin"
)

func (*user) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user := new(model.User)
	id := "myid"

	// Get single User record that `id is 'myid'.
	//
	// SQL equivalent:
	// SELECT * FROM users
	// WHERE id = 'myid'
	// LIMIT 1`
	database.Database[*model.User]().Get(user, id)
	database.Database[*model.User](ctx).Get(user, id)

	// List users without triggering model's hooks.
	// It's useful to break recursive model hook triggers.
	database.Database[*model.User]().WithoutHook().Get(user, id)

	// Get users from cache if exists, otherwise fetch from database.
	database.Database[*model.User]().WithCache().Get(user, id)

	// =====================================================================
	// Why choose database.Database().Create/Delete/Update/List/Get methods:
	// beacause model's hooks only invoke in database.Database.
	// =====================================================================
}
