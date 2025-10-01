package app

import (
	"demo/model/config/namespace"

	"github.com/forbearing/gst/service"
	"github.com/forbearing/gst/types"
)

type Lister struct {
	service.Base[*namespace.App, *namespace.App, *namespace.AppRsp]
}

func (a *Lister) List(ctx *types.ServiceContext, req *namespace.App) (rsp *namespace.AppRsp, err error) {
	log := a.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("app list")
	return rsp, nil
}

func (a *Lister) ListBefore(ctx *types.ServiceContext, apps *[]*namespace.App) error {
	log := a.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("app list before")
	return nil
}

func (a *Lister) ListAfter(ctx *types.ServiceContext, apps *[]*namespace.App) error {
	log := a.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("app list after")
	return nil
}
