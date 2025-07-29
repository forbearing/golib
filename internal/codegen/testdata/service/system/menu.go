package service_system

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_system "nebula/model/system"
)

func init() {
	service.Register[*menu]()
}

// menu implements the types.Service[*model_system.Menu] interface.
type menu struct {
	service.Base[*model_system.Menu]
}

func (m *menu) CreateBefore(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) CreateAfter(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) DeleteBefore(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) DeleteAfter(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) UpdateBefore(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) UpdateAfter(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) UpdatePartialBefore(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) UpdatePartialAfter(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) ListBefore(ctx *types.ServiceContext, menus *[]*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) ListAfter(ctx *types.ServiceContext, menus *[]*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) GetBefore(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) GetAfter(ctx *types.ServiceContext, menu *model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchCreateBefore(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchCreateAfter(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchDeleteBefore(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchDeleteAfter(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchUpdateBefore(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchUpdateAfter(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchUpdatePartialBefore(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (m *menu) BatchUpdatePartialAfter(ctx *types.ServiceContext, menus ...*model_system.Menu) error {
	log := m.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("menu batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
