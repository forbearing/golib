package model

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/database"
)

type Role struct {
	Name string `json:"name,omitempty" schema:"name"`

	Base
}

func (r *Role) check() error {
	if len(strings.TrimSpace(r.Name)) == 0 {
		return errors.New("name is required")
	}
	roles := make([]*Role, 0)
	if err := database.Database[*Role]().WithLimit(1).WithQuery(&Role{Name: r.Name}).List(&roles); err != nil {
		return err
	}
	if len(roles) > 0 {
		return fmt.Errorf("role(%s) already exists", r.Name)
	}

	return nil
}

func (r *Role) CreateBefore() error { return r.check() }
func (r *Role) CreateAfter() error  { return rbac.RBAC().AddRole(r.Name) }
func (r *Role) DeleteAfter() error  { return rbac.RBAC().RemoveRole(r.Name) }
