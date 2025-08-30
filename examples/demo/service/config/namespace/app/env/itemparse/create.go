package itemparse

import (
	"demo/model/config/namespace/app/env"

	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

type Creator struct {
	service.Base[*env.ItemParse, *env.ItemParse, *env.ItemParseRsp]
}

func (i *Creator) Create(ctx *types.ServiceContext, req *env.ItemParse) (rsp *env.ItemParseRsp, err error) {
	log := i.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("itemparse create")
	return rsp, nil
}

func (i *Creator) CreateBefore(ctx *types.ServiceContext, itemparse *env.ItemParse) error {
	log := i.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("itemparse create before")
	return nil
}

func (i *Creator) CreateAfter(ctx *types.ServiceContext, itemparse *env.ItemParse) error {
	log := i.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("itemparse create after")
	return nil
}
