package service_setting

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_setting "nebula/model/setting"
)

func init() {
	service.Register[*vendor]()
}

// vendor implements the types.Service[*model_setting.Vendor] interface.
type vendor struct {
	service.Base[*model_setting.Vendor]
}

func (v *vendor) CreateBefore(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) CreateAfter(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) DeleteBefore(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) DeleteAfter(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) UpdateBefore(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) UpdateAfter(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) UpdatePartialBefore(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) UpdatePartialAfter(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) ListBefore(ctx *types.ServiceContext, vendors *[]*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) ListAfter(ctx *types.ServiceContext, vendors *[]*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) GetBefore(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) GetAfter(ctx *types.ServiceContext, vendor *model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchCreateBefore(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchCreateAfter(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchDeleteBefore(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchDeleteAfter(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchUpdateBefore(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchUpdateAfter(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchUpdatePartialBefore(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (v *vendor) BatchUpdatePartialAfter(ctx *types.ServiceContext, vendors ...*model_setting.Vendor) error {
	log := v.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("vendor batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
