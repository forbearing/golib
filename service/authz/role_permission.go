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

type rolePermission struct {
	service.Base[*model_authz.RolePermission, *model_authz.RolePermission, *model_authz.RolePermission]
}

func init() {
	service.Register[*rolePermission](consts.PHASE_DELETE)
}

// DeleteAfter support delete multiple role_permissions by query parameters `role`, `resource`, `action`
func (*rolePermission) DeleteAfter(ctx *types.ServiceContext, rolePermission *model_authz.RolePermission) error {
	log := logger.Service.WithServiceContext(ctx, consts.PHASE_DELETE_AFTER)
	role := ctx.URL.Query().Get("role")
	resource := ctx.URL.Query().Get("resource")
	action := ctx.URL.Query().Get("action")

	rolePermissions := make([]*model_authz.RolePermission, 0)
	if err := database.Database[*model_authz.RolePermission]().WithLimit(-1).WithQuery(&model_authz.RolePermission{
		Role:     role,
		Resource: resource,
		Action:   action,
	}).List(&rolePermissions); err != nil {
		log.Error(err)
		return err
	}
	for _, rp := range rolePermissions {
		log.Infoz("will delete role permission", zap.Object("role_permission", rp))
	}
	if err := database.Database[*model_authz.RolePermission]().WithLimit(-1).WithPurge().Delete(rolePermissions...); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
