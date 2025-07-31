package service

import (
	"nebula/model"

	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

func init() {
	service.Register[*user2]()
}

// user2 implements the types.Service[*model.User2] interface.
type user2 struct {
	service.Base[*model.User2]
}

func (u *user2) CreateBefore(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) CreateAfter(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) DeleteBefore(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) DeleteAfter(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) UpdateBefore(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) UpdateAfter(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) UpdatePartialBefore(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) UpdatePartialAfter(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) ListBefore(ctx *types.ServiceContext, user2s *[]*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) ListAfter(ctx *types.ServiceContext, user2s *[]*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) GetBefore(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) GetAfter(ctx *types.ServiceContext, user2 *model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchCreateBefore(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchCreateAfter(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchDeleteBefore(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchDeleteAfter(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchUpdateBefore(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchUpdateAfter(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchUpdatePartialBefore(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}

func (u *user2) BatchUpdatePartialAfter(ctx *types.ServiceContext, user2s ...*model.User2) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user2 batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
