package model_authz

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap/zapcore"
)

func init() {
	model.Register[*UserRole]()
}

type UserRole struct {
	UserId string `json:"user_id,omitempty" schema:"user_id"`
	RoleId string `json:"role_id,omitempty" schema:"role_id"`

	User string `json:"user,omitempty" schema:"user"` // 用户名,只是为了方便查询
	Role string `json:"role,omitempty" schema:"role"` // 角色名,只是为了方便查询

	model.Base
}

func (r *UserRole) CreateBefore() error {
	if len(r.UserId) == 0 {
		return errors.New("user_id is required")
	}
	if len(r.RoleId) == 0 {
		return errors.New("role_id is required")
	}
	// expands field: user and role
	user, role := new(model.User), new(Role)
	if err := database.Database[*model.User](nil).Get(user, r.UserId); err != nil {
		return err
	}
	if err := database.Database[*Role](nil).Get(role, r.RoleId); err != nil {
		return err
	}
	r.User, r.Role = user.Name, role.Name

	// If the user already has the role, set same id to just update it.
	r.SetID(util.HashID(r.UserId, r.RoleId))

	return nil
}

func (r *UserRole) CreateAfter() error {
	if err := database.Database[*UserRole](nil).Update(r); err != nil {
		return err
	}
	// NOTE: must be role name not role id.
	if err := rbac.RBAC().AssignRole(r.UserId, r.Role); err != nil {
		return err
	}

	// update casbin_rule field: `user`, `role`, `remark`
	user := new(model.User)
	if err := database.Database[*model.User](nil).Get(user, r.UserId); err != nil {
		return err
	}
	casbinRules := make([]*CasbinRule, 0)
	if err := database.Database[*CasbinRule](nil).WithLimit(1).WithQuery(&CasbinRule{V0: r.UserId, V1: r.Role}).List(&casbinRules); err != nil {
		return err
	}
	if len(casbinRules) > 0 {
		casbinRules[0].User = user.Name
		casbinRules[0].Role = r.Role
		casbinRules[0].Remark = util.ValueOf(fmt.Sprintf("%s -> %s", r.User, r.Role))
		return database.Database[*CasbinRule](nil).Update(casbinRules[0])
	}
	return nil
}

func (r *UserRole) DeleteBefore() error {
	// The delete request always don't have user_id and role_id, so we should get the role from database.
	if err := database.Database[*UserRole](nil).Get(r, r.ID); err != nil {
		return err
	}
	// NOTE: must be role name not role id.
	return rbac.RBAC().UnassignRole(r.UserId, r.Role)
}

func (r *UserRole) DeleteAfter() error { return database.Database[*UserRole](nil).Cleanup() }

func (r *UserRole) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if r == nil {
		return nil
	}
	enc.AddString("user_id", r.UserId)
	enc.AddString("role_id", r.RoleId)
	enc.AddString("user", r.User)
	enc.AddString("role", r.Role)
	enc.AddObject("base", &r.Base)
	return nil
}
