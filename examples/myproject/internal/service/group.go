package service

import (
	"github.com/forbearing/golib/examples/myproject/internal/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
)

func init() {
	service.Register[*group]()
}

type group struct {
	service.Base[*model.Group]
}

func (g *group) CreateBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_CREATE_BEFORE)
	log.Info("group create before")
	return nil
}

func (g *group) CreateAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_CREATE_AFTER)
	log.Info("group create after")
	return nil
}

func (g *group) DeleteBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_DELETE_BEFORE)
	log.Info("group delete before")
	return nil
}

func (g *group) DeleteAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_DELETE_AFTER)
	log.Info("group delete after")
	return nil
}

func (g *group) UpdateBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_UPDATE_BEFORE)
	log.Info("group update before")
	return nil
}

func (g *group) UpdateAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_UPDATE_AFTER)
	log.Info("group update after")
	return nil
}

func (g *group) UpdatePartialBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_UPDATE_PARTIAL_BEFORE)
	log.Info("group update partial before")
	return nil
}

func (g *group) UpdatePartialAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_UPDATE_PARTIAL_AFTER)
	log.Info("group update partial after")
	return nil
}

func (g *group) ListBefore(ctx *types.ServiceContext, groups *[]*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_LIST_BEFORE)
	log.Info("group list before")
	return nil
}

func (g *group) ListAfter(ctx *types.ServiceContext, groups *[]*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_LIST_AFTER)
	log.Info("group list after")
	return nil
}

func (g *group) GetBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_GET_BEFORE)
	log.Info("group get before")
	return nil
}

func (g *group) GetAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_GET_AFTER)
	log.Info("group get after")
	return nil
}

func (g *group) BatchCreateBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_CREATE_BEFORE)
	log.Info("group batch create before")
	return nil
}

func (g *group) BatchCreateAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_CREATE_AFTER)
	log.Info("group batch create after")
	return nil
}

func (g *group) BatchDeleteBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_DELETE_BEFORE)
	log.Info("group batch delete before")
	return nil
}

func (g *group) BatchDeleteAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_DELETE_AFTER)
	log.Info("group batch delete after")
	return nil
}

func (g *group) BatchUpdateBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_BEFORE)
	log.Info("group batch update before")
	return nil
}

func (g *group) BatchUpdateAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_AFTER)
	log.Info("group batch update after")
	return nil
}

func (g *group) BatchUpdatePartialBefore(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_PARTIAL_BEFORE)
	log.Info("group batch update partial before")
	return nil
}

func (g *group) BatchUpdatePartialAfter(ctx *types.ServiceContext, groups ...*model.Group) error {
	log := g.WithServiceContext(ctx, consts.PHASE_BATCH_UPDATE_PARTIAL_AFTER)
	log.Info("group batch update partial after")
	return nil
}
