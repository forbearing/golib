package service

import (
	"demo/model"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
)

// init registers user service layer.
// NOTE: you should always ensure current service package import directly or indirectly in `main.go`.
func init() {
	// // Register user service layer, generic type *model.User can be inferred.
	// service.Register[*user]()

	// // Alternatively, you can explicitly specify both types.
	// service.Register[*user, *model.User]()

	// Register user service with custom fields initialization
	service.Register(&user{Field1: "value1", Field2: "value2"})
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
	log := u.WithServiceContext(ctx, consts.PHASE_CREATE_BEFORE)
	log.Info("user create before")
	// =============================
	// Add your business logic here.
	// =============================

	// example1: you can operate database in service layer.
	for i := range users {
		_ = users[i]
		logs := make([]*model.Log, 0)
		database.Database[*model.Log]().List(&logs)
	}
	return nil
}

func (u *user) CreateAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_CREATE_AFTER)
	log.Info("user create after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) DeleteBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_DELETE_BEFORE)
	log.Info("user delete before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) DeleteAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_DELETE_AFTER)
	log.Info("user delete after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) UpdateBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_UPDATE_BEFORE)
	log.Info("user update before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) UpdateAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_UPDATE_AFTER)
	log.Info("user update after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) UpdatePartialBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_UPDATE_PARTIAL_BEFORE)
	log.Info("user update partial before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) UpdatePartialAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_UPDATE_PARTIAL_AFTER)
	log.Info("user update partial after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_LIST_BEFORE)
	log.Info("user list before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) ListAfter(ctx *types.ServiceContext, users *[]*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_LIST_AFTER)
	log.Info("user list after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) GetBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_GET_BEFORE)
	log.Info("user get before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) GetAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_GET_AFTER)
	log.Info("user get after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchCreateBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_CREATE_BEFORE)
	log.Info("user batch create before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchCreateAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_CREATE_AFTER)
	log.Info("user batch create after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchDeleteBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_DELETE_BEFORE)
	log.Info("user batch delete before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchDeleteAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_DELETE_AFTER)
	log.Info("user batch delete after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchUpdateBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_BEFORE)
	log.Info("user batch update before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchUpdateAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_AFTER)
	log.Info("user batch update after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchUpdatePartialBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_PARTIAL_BEFORE)
	log.Info("user batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}

func (u *user) BatchUpdatePartialAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_PARTIAL_AFTER)
	log.Info("user batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================
	return nil
}
