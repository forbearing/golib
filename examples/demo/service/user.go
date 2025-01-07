package service

import (
	"demo/model"

	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
)

func init() {
	// // Register user service layer, generic type *model.User can be inferred.
	// service.Register[*user]()

	// // Alternatively, you can explicitly specify both types.
	// service.Register[*user, *model.User]()

	// Register user service with custom fields initialization
	service.Register[*user](&user{Field1: "value1", Field2: "value2"})
}

// user implements the types.Service[*model.User] interface.
// service.Base[*model.User] is a service layer associated with *model.User.
// It's strongly recommended to set user unexported.
type user struct {
	Field1 string
	Field2 string

	service.Base[*model.User]
}

func (u *user) CreateBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := types.LoggerWithContext(ctx, u, consts.PHASE_CREATE_BEFORE)
	log.Info("user create before")
	// =============================
	// Add your business logic here.
	// =============================
	for i := range users {
		_ = users[i]
	}
	return nil
}
