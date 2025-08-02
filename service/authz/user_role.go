package service_authz

import (
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/logger"
	model_authz "github.com/forbearing/golib/model/authz"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"go.uber.org/zap"
)

type userRole struct {
	service.Base[*model_authz.UserRole, *model_authz.UserRole, *model_authz.UserRole]
}

func init() {
	service.Register[*userRole]()
}

// DeleteAfter support filter and delete multiple user_roles by query parameter `user` and `role`.
func (r *userRole) DeleteAfter(ctx *types.ServiceContext, userRole *model_authz.UserRole) error {
	log := logger.Service.WithServiceContext(ctx, consts.PHASE_DELETE_AFTER)
	user := ctx.URL.Query().Get("user")
	role := ctx.URL.Query().Get("role")

	userRoles := make([]*model_authz.UserRole, 0)
	if err := database.Database[*model_authz.UserRole]().WithQuery(&model_authz.UserRole{User: user, Role: role}).WithLimit(-1).List(&userRoles); err != nil {
		log.Error(err)
		return err
	}
	for _, rb := range userRoles {
		log.Infoz("will delete user role", zap.Object("user_role", rb))
	}
	if err := database.Database[*model_authz.UserRole]().WithLimit(-1).WithPurge().Delete(userRoles...); err != nil {
		return err
	}

	return nil
}
