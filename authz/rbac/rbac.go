package rbac

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/forbearing/gst/types"
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

// AddRole is a no-op in Casbin, roles are created implicitly when used.
func (r *rbac) AddRole(name string) error {
	return nil
}

func (r *rbac) RemoveRole(name string) error {
	if _, err := r.enforcer.DeleteRole(name); err != nil {
		return err
	}
	return r.enforcer.SavePolicy()
}

func (r *rbac) GrantPermission(role string, resource string, action string) error {
	if _, err := r.enforcer.AddPermissionForUser(role, resource, action, "allow"); err != nil {
		return err
	}
	return r.enforcer.SavePolicy()
}

func (r *rbac) RevokePermission(role string, resource string, action string) error {
	if _, err := r.enforcer.DeletePermissionForUser(role, resource, action, "allow"); err != nil {
		return err
	}
	return r.enforcer.SavePolicy()
}

func (r *rbac) AssignRole(subject string, role string) error {
	if _, err := r.enforcer.AddRoleForUser(subject, role); err != nil {
		return err
	}
	return r.enforcer.SavePolicy()
}

func (r *rbac) UnassignRole(subject string, role string) error {
	if _, err := r.enforcer.DeleteRoleForUser(subject, role); err != nil {
		return err
	}
	return r.enforcer.SavePolicy()
}

// | 操作             | 函数                                  |
// | ---------------- | ------------------------------------- |
// | 添加角色权限     | `AddPolicy(role, obj, act)`           |
// | 删除角色权限     | `RemovePolicy(...)`                   |
// | 给用户授权角色   | `AddGroupingPolicy(user, role)`       |
// | 删除用户授权     | `RemoveGroupingPolicy(user, role)`    |
// | 查询用户角色     | `GetRolesForUser(user)`               |
// | 查询角色权限     | `GetPermissionsForUser(role)`         |
// | 查询用户所有权限 | `GetImplicitPermissionsForUser(user)` |

// // 查询用户拥有的角色
// RBAC.enforcer.GetRolesForUser("root")
// // 查询角色拥有的权限
// RBAC.enforcer.GetFilteredPolicy(0, "admin")
// // 查询用户拥有的权限（继承）
// RBAC.enforcer.GetImplicitPermissionsForUser("root")
