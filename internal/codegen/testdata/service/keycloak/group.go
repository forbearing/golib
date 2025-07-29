package service_keycloak

import (
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	model_keycloak "nebula/model/keycloak"
)

func init() {
	service.Register[*group]()
}

// group implements the types.Service[*model_keycloak.Group] interface.
type group struct {
	service.Base[*model_keycloak.Group]
}

func (g *group) CreateBefore(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) CreateAfter(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) DeleteBefore(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) DeleteAfter(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) UpdateBefore(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) UpdatePartialBefore(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) UpdatePartialAfter(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) ListBefore(ctx *types.ServiceContext, groups *[]*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group list before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) ListAfter(ctx *types.ServiceContext, groups *[]*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group list after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) GetBefore(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group get before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) GetAfter(ctx *types.ServiceContext, group *model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group get after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchCreateBefore(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchCreateAfter(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch create after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchDeleteBefore(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch delete before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchDeleteAfter(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch delete after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchUpdateBefore(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchUpdateAfter(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchUpdatePartialBefore(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update partial before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
func (g *group) BatchUpdatePartialAfter(ctx *types.ServiceContext, groups ...*model_keycloak.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update partial after")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
