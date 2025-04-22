package model_authz

import (
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"go.uber.org/zap/zapcore"
)

type RolePermission struct {
	RoleId       string `json:"role_id" schema:"role_id"`
	PermissionId string `json:"permission_id" schema:"permission_id"`

	Role string `json:"role" schema:"role"` // 角色名,只是为了方便查询

	model.Base
}

func (r *RolePermission) DeleteAfter() error { return database.Database[*RolePermission]().Cleanup() }

func (r *RolePermission) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if r == nil {
		return nil
	}
	enc.AddString("role_id", r.RoleId)
	enc.AddString("permission_id", r.PermissionId)
	enc.AddString("role", r.Role)
	return nil
}
