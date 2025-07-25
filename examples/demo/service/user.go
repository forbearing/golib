package service

import (
	"demo/model"

	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

// init registers user service layer.
// NOTE: you should always ensure current service package import directly or indirectly in `main.go`.
func init() {
	service.Register[*user]()
}

// user implements the types.Service[*model.User] interface.
// service.Base[*model.User] is a service layer associated with *model.User.
// It's strongly recommended to set user unexported.
type user struct {
	service.Base[*model.User]
}

func (u *user) CreateBefore(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) CreateAfter(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) DeleteBefore(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) DeleteAfter(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) UpdateBefore(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) UpdateAfter(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) UpdatePartialBefore(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) UpdatePartialAfter(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) ListAfter(ctx *types.ServiceContext, users *[]*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) GetBefore(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) GetAfter(ctx *types.ServiceContext, users *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchCreateBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchCreateAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchDeleteBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchDeleteAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchUpdateBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchUpdateAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchUpdatePartialBefore(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user) BatchUpdatePartialAfter(ctx *types.ServiceContext, users ...*model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
