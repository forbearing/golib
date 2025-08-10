package service

import (
	"helloworld/model"

	"github.com/forbearing/golib/types"
)

func (u *user) Create(ctx *types.ServiceContext, req *model.UserReq) (*model.UserRsp, error) {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create")
	return &model.UserRsp{}, nil
}

func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create before")
	return nil
}

func (u *user) CreateAfter(ctx *types.ServiceContext, user *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create after")
	return nil
}
