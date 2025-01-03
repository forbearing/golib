package rbac

import (
	"strings"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/types/consts"
)

const (
	root = "root"
)

// role 其实也是 group, 只需要通过 AddGroupingPolicy 将用户加进组来就行.
// 加进组其实就是赋予权限或角色.
var roles = []string{
	"root",
	"admin",
	"user",
	"guest",
}

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
	_, err := RBAC.AddPolicy(r.name, url, method)
	return err
}

// DelPolicy 角色组删除策略
func (r *role) DelPolicy(url, method string) error {
	_, err := RBAC.RemovePolicy(r.name, url, method)
	return err
}

// BindTo 给指定用户绑定该角色
func (r *role) BindTo(user string) error {
	_, err := RBAC.AddGroupingPolicy(user, r.name)
	return err
}

// UnbindTo 给指定用户解除绑定该角色
func (r *role) UnbindTo(user string) error {
	_, err := RBAC.RemoveGroupingPolicy(user, r.name)
	return err
}

var RBAC = new(rbac)

// https://blog.csdn.net/LeoForBest/article/details/133610889
// https://juejin.cn/post/7269563694676819968
func Init() (err error) {
	if !config.App.ServerConfig.EnableRBAC {
		return nil
	}

	// casbin.NewEnforcer("model.conf", "policy.csv") 从文件中加载策略
	if RBAC.adapter, err = gormadapter.NewAdapterByDB(database.DB); err != nil {
		return
	}
	if RBAC.enforcer, err = casbin.NewEnforcer(consts.FileRbacConf, RBAC.adapter); err != nil {
		return
	}
	// RBAC.enforcer.AddFunction("isAdmin", func(args ...any) (any, error) {
	// 	username := args[0].(string)
	// 	return RBAC.enforcer.HasRoleForUser(username, "root")
	// })
	if strings.ToLower(config.App.LogLevel) == "debug" {
		RBAC.enforcer.EnableLog(true)
	}

	// 将 root 用户加入所有的所有角色组
	for _, role := range roles {
		if _, err = RBAC.AddGroupingPolicy(root, role); err != nil {
			return err
		}
	}

	// for _, route := range model.Routes {
	// 	path := route.Path
	// 	path = strings.TrimPrefix(path, `/api`)
	// 	path = strings.TrimPrefix(path, `/`)
	// 	path = strings.TrimSuffix(path, `/`)
	// 	path = filepath.Join("/api", path)
	// 	// deduplicate and translate 'VerbMost' and 'VerboAll'
	// 	verbMap := make(map[model.Verb]struct{})
	// 	for _, verb := range route.Verbs {
	// 		switch verb {
	// 		case model.VerbCreate, model.VerbDelete, model.VerbUpdate,
	// 			model.VerbUpdatePartial, model.VerbList, model.VerbGet,
	// 			model.VerbImport, model.VerbExport:
	// 			verbMap[verb] = struct{}{}
	// 		case model.VerbMost:
	// 			verbMap[model.VerbCreate] = struct{}{}
	// 			verbMap[model.VerbDelete] = struct{}{}
	// 			verbMap[model.VerbUpdate] = struct{}{}
	// 			verbMap[model.VerbUpdatePartial] = struct{}{}
	// 			verbMap[model.VerbList] = struct{}{}
	// 			verbMap[model.VerbGet] = struct{}{}
	// 		case model.VerbAll:
	// 			verbMap[model.VerbCreate] = struct{}{}
	// 			verbMap[model.VerbDelete] = struct{}{}
	// 			verbMap[model.VerbUpdate] = struct{}{}
	// 			verbMap[model.VerbUpdatePartial] = struct{}{}
	// 			verbMap[model.VerbList] = struct{}{}
	// 			verbMap[model.VerbGet] = struct{}{}
	// 			verbMap[model.VerbImport] = struct{}{}
	// 			verbMap[model.VerbExport] = struct{}{}
	// 		}
	// 	}
	// 	for verb := range verbMap {
	// 		switch verb {
	// 		case model.VerbCreate:
	// 			RBAC.AddPolicy(root, path, http.MethodPost)
	// 		case model.VerbDelete:
	// 			RBAC.AddPolicy(root, path, http.MethodDelete)
	// 		case model.VerbUpdate:
	// 			RBAC.AddPolicy(root, path, http.MethodPut)
	// 		case model.VerbUpdatePartial:
	// 			RBAC.AddPolicy(root, path, http.MethodPatch)
	// 		case model.VerbList:
	// 			RBAC.AddPolicy(root, path, http.MethodGet)
	// 		case model.VerbGet:
	// 			RBAC.AddPolicy(root, path+"/:id", http.MethodGet)
	// 		case model.VerbExport:
	// 			RBAC.AddPolicy(root, path+"/export", http.MethodGet)
	// 		case model.VerbImport:
	// 			RBAC.AddPolicy(root, path+"/import", http.MethodPost)
	// 		}
	// 	}
	// }

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
func (r *rbac) AddPolicy(role, url, method string) (affected bool, err error) {
	if err = r.enforcer.LoadPolicy(); err != nil {
		return
	}
	if affected, err = r.enforcer.AddPolicy(role, url, method); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}

// RemovePolicy 给角色组 role 删除策略
func (r *rbac) RemovePolicy(role, url, method string) (affected bool, err error) {
	if err = r.enforcer.LoadPolicy(); err != nil {
		return
	}
	if affected, err = r.enforcer.RemovePolicy(role, url, method); err != nil {
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
func (r *rbac) AddGroupingPolicy(user, role string) (affected bool, err error) {
	if affected, err = r.enforcer.AddGroupingPolicy(user, role); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}

// RemoveGroupingPolicy 将用户user从角色组role中移除.
func (r *rbac) RemoveGroupingPolicy(user, role string) (affected bool, err error) {
	if affected, err = r.enforcer.RemoveGroupingPolicy(user, role); err != nil {
		return
	}
	return affected, r.enforcer.SavePolicy()
}

// // GetAllRoles
// func (r *rbac) GetAllRoles() []string {
// 	return r.enforcer.GetAllRoles()
// }
//
// func (r *rbac) GetGroupingPolicy() [][]string {
// 	return r.enforcer.GetGroupingPolicy()
// }
