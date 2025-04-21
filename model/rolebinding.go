package model

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/database"
)

type Rolebinding struct {
	UserId string `json:"user_id,omitempty" schema:"user_id"`
	RolId  string `json:"role_id,omitempty" schema:"role_id"`

	Base
}

func (r *Rolebinding) check() error {
	if len(r.UserId) == 0 {
		return errors.New("user_id is required")
	}
	if len(r.RolId) == 0 {
		return errors.New("role_id is required")
	}
	bindings := make([]*Rolebinding, 0)
	if err := database.Database[*Rolebinding]().WithLimit(1).WithQuery(&Rolebinding{UserId: r.UserId, RolId: r.RolId}).List(&bindings); err != nil {
		return err
	}
	if len(bindings) > 0 {
		return fmt.Errorf("rolebinding(%s) already exists", bindings[0].ID)
	}

	return nil
}

func (r *Rolebinding) CreateBefore() error { return r.check() }
func (r *Rolebinding) CreateAfter() error  { return rbac.RBAC().AssignRole(r.UserId, r.RolId) }
func (r *Rolebinding) DeleteAfter() error  { return rbac.RBAC().UnassignRole(r.UserId, r.RolId) }
