package service_setting

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_setting "nebula/model/setting"
)

func init() {
	service.Register[*tenant]()
}

// tenant implements the types.Service[*model_setting.Tenant] interface.
type tenant struct {
	service.Base[*model_setting.Tenant]
}

func (t *tenant) CreateBefore(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) CreateAfter(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) DeleteBefore(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) DeleteAfter(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) UpdateBefore(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) UpdateAfter(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) UpdatePartialBefore(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) UpdatePartialAfter(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) ListBefore(ctx *types.ServiceContext, tenants *[]*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) ListAfter(ctx *types.ServiceContext, tenants *[]*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) GetBefore(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) GetAfter(ctx *types.ServiceContext, tenant *model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchCreateBefore(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchCreateAfter(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchDeleteBefore(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchDeleteAfter(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchUpdateBefore(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchUpdateAfter(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchUpdatePartialBefore(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (t *tenant) BatchUpdatePartialAfter(ctx *types.ServiceContext, tenants ...*model_setting.Tenant) error {
	log := t.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("tenant batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
