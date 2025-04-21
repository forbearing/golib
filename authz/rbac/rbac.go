package rbac

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/forbearing/golib/types"
)

var (
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
)

type rbac struct {
	enforcer *casbin.Enforcer
	addapter *gormadapter.Adapter
}

func RBAC() types.RBAC {
	return &rbac{
		enforcer: Enforcer,
		addapter: Adapter,
	}
}

func (r *rbac) AddRole(name string) error    { return nil }
func (r *rbac) RemoveRole(name string) error { return nil }

func (r *rbac) GrantPermission(role string, resource string, action string) error  { return nil }
func (r *rbac) RevokePermission(role string, resource string, action string) error { return nil }

func (r *rbac) AssignRole(subject string, role string) error   { return nil }
func (r *rbac) UnassignRole(subject string, role string) error { return nil }
