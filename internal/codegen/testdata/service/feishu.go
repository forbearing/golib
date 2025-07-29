package service

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"nebula/model"
)

func init() {
	service.Register[*feishu]()
}

// feishu implements the types.Service[*model.Feishu] interface.
type feishu struct {
	service.Base[*model.Feishu]
}

func (f *feishu) CreateBefore(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) CreateAfter(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) DeleteBefore(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) DeleteAfter(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) UpdateBefore(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) UpdateAfter(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) UpdatePartialBefore(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) UpdatePartialAfter(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) ListBefore(ctx *types.ServiceContext, feishus *[]*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) ListAfter(ctx *types.ServiceContext, feishus *[]*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) GetBefore(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) GetAfter(ctx *types.ServiceContext, feishu *model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchCreateBefore(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchCreateAfter(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchDeleteBefore(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchDeleteAfter(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchUpdateBefore(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchUpdateAfter(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchUpdatePartialBefore(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (f *feishu) BatchUpdatePartialAfter(ctx *types.ServiceContext, feishus ...*model.Feishu) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("feishu batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
