package service

import (
	"codegen/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)

func init() {
	service.Register[*user]()
}
// user implements the types.Service[*model.User] interface.
type user struct {
	service.Base[*model.User]
}

func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {
	log := u.WithServiceContext(ctx, ctx.GetPhase())
	log.Info("user create before")
	// =============================
	// Add your business logic here.
	// =============================

	return nil
}
