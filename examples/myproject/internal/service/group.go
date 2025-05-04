package service

import (
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/examples/myproject/internal/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
)

func init() {
	service.Register[*group]()
}

type group struct {
	service.Base[*model.Group]
}

func (g *group) CreateBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group create before")

	// Has Custom Request and Response
	req, ok := ctx.GetRequest().(*model.GroupRequest)
	if ok {
		if err := database.Database[*model.Group]().Create(&model.Group{Name: req.Name}); err != nil {
			log.Error(err)
		}
	}
	resp, ok := ctx.GetResponse().(*model.GroupResponse)
	if ok {
		resp.CustomName = "custom name in create before"
		resp.CustomDesc = util.ValueOf("desc in create before")
	}
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
	resp, ok := ctx.GetResponse().(*model.GroupResponse)
	if ok {
		resp.CustomName = "custom name in update before"
		resp.CustomDesc = util.ValueOf("desc in update before")
	}
	log.Info("group update before")
	return nil
}

func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group update after")
	return nil
}

func (g *group) UpdatePartialBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	resp, ok := ctx.GetResponse().(*model.GroupResponse)
	if ok {
		resp.CustomName = "custom name in update partial before"
		resp.CustomDesc = util.ValueOf("custom desc in update partial before")
	}
	log.Info("group update partial before")
	return nil
}

func (g *group) UpdatePartialAfter(ctx *types.ServiceContext, group *model.Group) error {
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

func (g *group) BatchCreateBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	resp, ok := ctx.GetResponse().(*model.GroupResponse)
	if ok {
		resp.CustomName = "custom name in batch create before"
		resp.CustomDesc = util.ValueOf("custom desc in batch create before")
	}
	log.Info("group batch create before")
	return nil
}

func (g *group) BatchCreateAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch create after")
	return nil
}

func (g *group) BatchDeleteBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch delete before")
	return nil
}

func (g *group) BatchDeleteAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch delete after")
	return nil
}

func (g *group) BatchUpdateBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	resp, ok := ctx.GetResponse().(*model.GroupResponse)
	if ok {
		resp.CustomName = "custom name in batch update before"
		resp.CustomDesc = util.ValueOf("custom desc in batch update before")
	}
	log.Info("group batch update before")
	return nil
}

func (g *group) BatchUpdateAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update after")
	return nil
}

func (g *group) BatchUpdatePartialBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	resp, ok := ctx.GetResponse().(*model.GroupResponse)
	if ok {
		resp.CustomName = "custom name in batch updte partial before"
		resp.CustomDesc = util.ValueOf("custom desc in batch update partial before")
	}
	log.Info("group batch update partial before")
	return nil
}

func (g *group) BatchUpdatePartialAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group batch update partial after")
	return nil
}
