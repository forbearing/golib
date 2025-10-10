package service_authz

import (
	"github.com/forbearing/gst/database"
	"github.com/forbearing/gst/logger"
	model_authz "github.com/forbearing/gst/model/modelauthz"
	"github.com/forbearing/gst/service"
	"github.com/forbearing/gst/types"
	"github.com/forbearing/gst/types/consts"
	"go.uber.org/zap"
)

type role struct {
	service.Base[*model_authz.Role, *model_authz.Role, *model_authz.Role]
}

func init() {
	service.Register[*role](consts.PHASE_DELETE)
}

// DeleteAfter support filter and delete multiple roles by query parameter `name`.
func (r *role) DeleteAfter(ctx *types.ServiceContext, role *model_authz.Role) error {
	log := logger.Service.WithServiceContext(ctx, consts.PHASE_DELETE_AFTER)
	name := ctx.URL.Query().Get("name")
	if len(name) == 0 {
		return nil
	}

	roles := make([]*model_authz.Role, 0)
	if err := database.Database[*model_authz.Role](ctx.DatabaseContext()).WithLimit(-1).WithQuery(&model_authz.Role{Name: name}).List(&roles); err != nil {
		log.Error(err)
		return err
	}
	for _, role := range roles {
		log.Infoz("will delete role", zap.Object("role", role))
	}
	if err := database.Database[*model_authz.Role](ctx.DatabaseContext()).WithLimit(-1).WithPurge().Delete(roles...); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
