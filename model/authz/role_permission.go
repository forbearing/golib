package model_authz

import (
	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap/zapcore"
)

func init() {
	model.Register[*RolePermission]()
}

type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny  Effect = "deny"
)

// TODO: remove RoleId and only keep Role(role name).
type RolePermission struct {
	Role string `json:"role" schema:"role"`

	Resource string `json:"resource" schema:"resource"`
	Action   string `json:"action" schema:"action"`
	Effect   Effect `json:"effect" schema:"effect"`

	model.Base
}

func (r *RolePermission) CreateBefore() error {
	if len(r.Role) == 0 {
		return errors.New("role_id is required")
	}
	if len(r.Resource) == 0 {
		return errors.New("resource is required")
	}
	if len(r.Action) == 0 {
		return errors.New("action is required")
	}

	// default effect is allow.
	switch r.Effect {
	case EffectAllow, EffectDeny:
	default:
		r.Effect = EffectAllow
	}
	// If the role already has the permission(Resource+Action), set same id to just update it.
	r.SetID(util.HashID(r.Role, r.Resource, r.Action))

	return nil
}

func (r *RolePermission) CreateAfter() error {
	// grant the permission: (role, resource, action)
	return rbac.RBAC().GrantPermission(r.Role, r.Resource, r.Action)
}

func (r *RolePermission) DeleteBefore() error {
	// The request always only contains id, so we should get the RolePermission from database.
	if err := database.Database[*RolePermission](nil).Get(r, r.ID); err != nil {
		return err
	}
	// revoke the role's permission
	return rbac.RBAC().RevokePermission(r.Role, r.Resource, r.Action)
}

func (r *RolePermission) DeleteAfter() error {
	return database.Database[*RolePermission](nil).Cleanup()
}

func (r *RolePermission) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if r == nil {
		return nil
	}
	enc.AddString("role", r.Role)
	enc.AddString("resource", r.Resource)
	enc.AddString("action", r.Action)
	enc.AddString("effect", string(r.Effect))
	return nil
}
