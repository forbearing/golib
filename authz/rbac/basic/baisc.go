package basic

import (
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/model"
)

const (
	Root = "root"
)

var adminRole = "admin"

var defaultAdmins = []string{
	"admin",
	"root",
}

var modelData = []byte(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
#e = priority(p.eft) || some(where (p.eft == allow))
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch3(r.obj, p.obj) && r.act == p.act)
`)

func Init() (err error) {
	if !config.App.Auth.RBACEnable {
		return nil
	}

	filename := filepath.Join(config.Tempdir(), "casbin_model.conf")
	if err = os.WriteFile(filename, modelData, 0o644); err != nil {
		return errors.Wrapf(err, "failed to write model file %s", filename)
	}
	// NOTE: gormadapter.NewAdapterByDBWithCustomTable creates the Casbin policy table with an auto-incrementing primary key.
	if rbac.Adapter, err = gormadapter.NewAdapterByDBWithCustomTable(database.DB, new(model.CasbinRule)); err != nil {
		return errors.Wrap(err, "failed to create casbin adapter")
	}
	if rbac.Enforcer, err = casbin.NewEnforcer(filename, rbac.Adapter); err != nil {
		return errors.Wrap(err, "failed to create casbin enforcer")
	}

	rbac.Enforcer.SetLogger(logger.Casbin)
	rbac.Enforcer.EnableLog(true)

	for _, user := range defaultAdmins {
		rbac.Enforcer.AddGroupingPolicy(user, adminRole)
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

	return rbac.Enforcer.LoadPolicy()
}
