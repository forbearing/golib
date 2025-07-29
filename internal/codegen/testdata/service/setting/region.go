package service_setting

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_setting "nebula/model/setting"
)

func init() {
	service.Register[*region]()
}

// region implements the types.Service[*model_setting.Region] interface.
type region struct {
	service.Base[*model_setting.Region]
}

func (r *region) CreateBefore(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) CreateAfter(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) DeleteBefore(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) DeleteAfter(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) UpdateBefore(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) UpdateAfter(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) UpdatePartialBefore(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) UpdatePartialAfter(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) ListBefore(ctx *types.ServiceContext, regions *[]*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) ListAfter(ctx *types.ServiceContext, regions *[]*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) GetBefore(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) GetAfter(ctx *types.ServiceContext, region *model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchCreateBefore(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchCreateAfter(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchDeleteBefore(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchDeleteAfter(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchUpdateBefore(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchUpdateAfter(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchUpdatePartialBefore(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *region) BatchUpdatePartialAfter(ctx *types.ServiceContext, regions ...*model_setting.Region) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("region batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
