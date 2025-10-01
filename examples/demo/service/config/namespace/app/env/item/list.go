package item

import (
	"demo/model/config/namespace/app/env"

	"github.com/forbearing/gst/service"
	"github.com/forbearing/gst/types"
)

type Lister struct {
	service.Base[*env.Item, *env.Item, *env.ItemRsp]
}

func (i *Lister) List(ctx *types.ServiceContext, req *env.Item) (rsp *env.ItemRsp, err error) {
	log := i.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("item list")
	return rsp, nil
}

func (i *Lister) ListBefore(ctx *types.ServiceContext, items *[]*env.Item) error {
	log := i.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("item list before")
	return nil
}

func (i *Lister) ListAfter(ctx *types.ServiceContext, items *[]*env.Item) error {
	log := i.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("item list after")
	return nil
}
