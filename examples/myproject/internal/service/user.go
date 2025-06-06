package service

import (
	pkgmodel "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

func init() {
	service.Register[*user]()
}

type user struct {
	service.Base[*pkgmodel.User]
}

func (u *user) CreateBefore(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create before")
	return nil
}

func (u *user) CreateAfter(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create after")
	return nil
}

func (u *user) DeleteBefore(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user delete before")
	return nil
}

func (u *user) DeleteAfter(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user delete after")
	return nil
}

func (u *user) UpdateBefore(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update before")
	return nil
}

func (u *user) UpdateAfter(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update after")
	return nil
}

func (u *user) UpdatePartialBefore(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update partial before")
	return nil
}

func (u *user) UpdatePartialAfter(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user update partial after")
	return nil
}

func (u *user) ListBefore(ctx *types.ServiceContext, users *[]*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user list before")
	return nil
}

func (u *user) ListAfter(ctx *types.ServiceContext, users *[]*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user list after")
	return nil
}

func (u *user) GetBefore(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user get before")
	return nil
}

func (u *user) GetAfter(ctx *types.ServiceContext, user *pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user get after")
	return nil
}

func (u *user) BatchCreateBefore(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch create before")
	return nil
}

func (u *user) BatchCreateAfter(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch create after")
	return nil
}

func (u *user) BatchDeleteBefore(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch delete before")
	return nil
}

func (u *user) BatchDeleteAfter(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch delete after")
	return nil
}

func (u *user) BatchUpdateBefore(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update before")
	return nil
}

func (u *user) BatchUpdateAfter(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update after")
	return nil
}

func (u *user) BatchUpdatePartialBefore(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update partial before")
	return nil
}

func (u *user) BatchUpdatePartialAfter(ctx *types.ServiceContext, users ...*pkgmodel.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user batch update partial after")
	return nil
}
