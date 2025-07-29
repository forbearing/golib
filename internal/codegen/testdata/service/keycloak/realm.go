package service_keycloak

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_keycloak "nebula/model/keycloak"
)

func init() {
	service.Register[*realm]()
}

// realm implements the types.Service[*model_keycloak.Realm] interface.
type realm struct {
	service.Base[*model_keycloak.Realm]
}

func (r *realm) CreateBefore(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) CreateAfter(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) DeleteBefore(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) DeleteAfter(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) UpdateBefore(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) UpdateAfter(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) UpdatePartialBefore(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) UpdatePartialAfter(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) ListBefore(ctx *types.ServiceContext, realms *[]*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) ListAfter(ctx *types.ServiceContext, realms *[]*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) GetBefore(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) GetAfter(ctx *types.ServiceContext, realm *model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchCreateBefore(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchCreateAfter(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchDeleteBefore(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchDeleteAfter(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchUpdateBefore(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchUpdateAfter(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchUpdatePartialBefore(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (r *realm) BatchUpdatePartialAfter(ctx *types.ServiceContext, realms ...*model_keycloak.Realm) error {
	log := r.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("realm batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
