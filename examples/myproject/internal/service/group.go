package service

import (
	"github.com/forbearing/golib/examples/myproject/internal/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
)

func init() {
	service.Register[*group](consts.PHASE_CREATE)
}

type group struct {
	service.Base[*model.Group, *model.Group, *model.Group]
}

func (g *group) CreateBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group create before")
	return nil
}

func (g *group) CreateAfter(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group create after")
	return nil
}

func (g *group) DeleteBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group delete before")
	return nil
}

func (g *group) DeleteAfter(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group delete after")
	return nil
}

func (g *group) UpdateBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update before")
	return nil
}

func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update after")
	return nil
}

func (g *group) PatchBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update partial before")
	return nil
}

func (g *group) PatchAfter(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update partial after")
	return nil
}

func (g *group) ListBefore(ctx *types.ServiceContext, groups *[]*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group list before")
	return nil
}

func (g *group) ListAfter(ctx *types.ServiceContext, groups *[]*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group list after")
	return nil
}

func (g *group) GetBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group get before")
	return nil
}

func (g *group) GetAfter(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group get after")
	return nil
}

func (g *group) CreateManyBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch create before")
	return nil
}

func (g *group) CreateManyAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch create after")
	return nil
}

func (g *group) DeleteManyBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch delete before")
	return nil
}

func (g *group) DeleteManyAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch delete after")
	return nil
}

func (g *group) UpdateManyBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update before")
	return nil
}

func (g *group) UpdateManyAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update after")
	return nil
}

func (g *group) PatchManyBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update partial before")
	return nil
}

func (g *group) PatchManyAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update partial after")
	return nil
}
