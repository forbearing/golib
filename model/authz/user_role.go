package model_authz

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
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
	bindings := make([]*UserRole, 0)
	if err := database.Database[*UserRole]().WithLimit(1).WithQuery(&UserRole{UserId: r.UserId, RoleId: r.RoleId}).List(&bindings); err != nil {
		return err
	}
	if len(bindings) > 0 {
		return fmt.Errorf("user_role(%s) already exists", bindings[0].ID)
	}

	// expands field: user and role
	user, role := new(model.User), new(Role)
	if err := database.Database[*model.User]().Get(user, r.UserId); err != nil {
		return err
	}
	if err := database.Database[*Role]().Get(role, r.RoleId); err != nil {
		return err
	}
	r.User, r.Role = user.Name, role.Name

	return nil
}

func (r *UserRole) CreateAfter() error {
	if err := database.Database[*UserRole]().Update(r); err != nil {
		return err
	}
	// TODO: add remark for casbin_rule.
	return rbac.RBAC().AssignRole(r.UserId, r.RoleId)
}

func (r *UserRole) DeleteBefore() error {
	// The delete request always don't have user_id and role_id, so we should get the role from database.
	if err := database.Database[*UserRole]().Get(r, r.ID); err != nil {
		return err
	}
	return rbac.RBAC().UnassignRole(r.UserId, r.RoleId)
}

func (r *UserRole) DeleteAfter() error { return database.Database[*UserRole]().Cleanup() }

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
