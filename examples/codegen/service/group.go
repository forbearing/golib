package service

import (
	"codegen/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

func init() {
	service.Register[*group]()
}
// group implements the types.Service[*model.Group] interface.
type group struct {
	service.Base[*model.Group]
}

func (g *group) CreateBefore(ctx *types.ServiceContext, group *model.Group) error {
	log := g.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("group create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
