package basic

import (
	"os"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/model"
)

const (
	root = "root"
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
p = priority, sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
#e = priority(p.eft) || some(where (p.eft == allow))
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, "admin") || (g(r.sub, p.sub) && keyMatch3(r.obj, p.obj) && r.act == p.act)
`)

// role 角色组
type role struct {
	name string
}

// Role 创建一个 role 对象,
// 后续直接调用 role 对象的方法来给 role 增删策略, 绑定/解除绑定用户即可.
// 返回的 role 对象拥有这些方法:
// AddPolicy: 给 role 增加策略
// DelPolicy: 给 role 删除策略
// BindTo: 给指定用户绑定该角色
// UnbindTo: 给指定用户解除绑定该角色
func Role(name string) *role { return &role{name: name} }

// AddPolicy 角色组添加策略
func (r *role) AddPolicy(url, method string) error {
	_, err := RBAC.AddPolicy(0, r.name, url, method)
	return err
}

// DelPolicy 角色组删除策略
func (r *role) DelPolicy(url, method string) error {
	_, err := RBAC.RemovePolicy(0, r.name, url, method)
	return err
}

// Bind 给指定用户绑定该角色
func (r *role) Bind(user string) error {
	_, err := RBAC.AddGroupingPolicy(0, user, r.name)
	return err
}

// Unbind 给指定用户解除绑定该角色
func (r *role) Unbind(user string) error {
	_, err := RBAC.RemoveGroupingPolicy(0, user, r.name)
	return err
}

var RBAC = new(rbac)

// https://blog.csdn.net/LeoForBest/article/details/133610889
// https://juejin.cn/post/7269563694676819968
func Init() (err error) {
	if !config.App.Auth.RBACEnable {
		return nil
	}

	filename := filepath.Join(config.Tempdir(), "rbac.conf")
	if err = os.WriteFile(filename, modelData, 0o644); err != nil {
		return errors.Wrapf(err, "failed to write model file %s", filename)
	}
	if RBAC.adapter, err = gormadapter.NewAdapterByDBWithCustomTable(database.DB, new(model.CasbinRule)); err != nil {
		return errors.Wrap(err, "failed to create casbin adapter")
	}
	if RBAC.enforcer, err = casbin.NewEnforcer(filename, RBAC.adapter); err != nil {
		return errors.Wrap(err, "failed to create casbin enforcer")
	}

	RBAC.enforcer.SetLogger(logger.Casbin)
	RBAC.enforcer.EnableLog(true)

	for _, user := range defaultAdmins {
		RBAC.enforcer.AddGroupingPolicy(user, adminRole)
	}

	return RBAC.enforcer.LoadPolicy()
}

type rbac struct {
	enforcer *casbin.Enforcer
	adapter  *gormadapter.Adapter
}

func (r *rbac) Enforcer(sub, obj, act string) (bool, error) {
	return r.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 给角色组 role 添加策略.
func (r *rbac) AddPolicy(priority int, role, url, method string) (affected bool, err error) {
	if err = r.enforcer.LoadPolicy(); err != nil {
		return
	}
	if affected, err = r.enforcer.AddPolicy(priority, role, url, method); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}

// RemovePolicy 给角色组 role 删除策略
func (r *rbac) RemovePolicy(priority int, role, url, method string) (affected bool, err error) {
	if err = r.enforcer.LoadPolicy(); err != nil {
		return
	}
	if affected, err = r.enforcer.RemovePolicy(priority, role, url, method); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}

// UpdateRole 给角色组 role 更新策略
func (r *rbac) UpdateRole(oldRole, oldUrl, oldMethod, newRole, newUrl, newMethod string) (affected bool, err error) {
	if err = r.enforcer.LoadPolicy(); err != nil {
		return
	}
	if affected, err = r.enforcer.UpdatePolicy(
		[]string{oldRole, oldUrl, oldMethod},
		[]string{newRole, newUrl, newMethod}); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}

// AddGroupingPolicy 将用户user加入进角色组role.
func (r *rbac) AddGroupingPolicy(priority int, user, role string) (affected bool, err error) {
	if affected, err = r.enforcer.AddGroupingPolicy(user, role); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}

// RemoveGroupingPolicy 将用户user从角色组role中移除.
func (r *rbac) RemoveGroupingPolicy(priority int, user, role string) (affected bool, err error) {
	if affected, err = r.enforcer.RemoveGroupingPolicy(priority, user, role); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}
