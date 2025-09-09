package file

import (
	"demo/model/config/namespace/app/env"

	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

type Lister struct {
	service.Base[*env.File, *env.File, *env.File]
}

func (f *Lister) List(ctx *types.ServiceContext, req *env.File) (rsp *env.File, err error) {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("file list")
	return rsp, nil
}

func (f *Lister) ListBefore(ctx *types.ServiceContext, files *[]*env.File) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("file list before")
	return nil
}

func (f *Lister) ListAfter(ctx *types.ServiceContext, files *[]*env.File) error {
	log := f.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("file list after")
	return nil
}

func (f *Lister) Filter(ctx *types.ServiceContext, file *env.File) *env.File {
	return file
}

func (f *Lister) FilterRaw(ctx *types.ServiceContext) string {
	return ""
}
