package model_keycloak

import (
	"github.com/forbearing/golib/model"
	"gorm.io/datatypes"
)

// type Group struct {
// 	// RealmID 组所属的 Realm
// 	RealmID string `json:"realm_id,omitempty" gorm:"type:varchar(36);index;not null"`
//
// 	// Name 是组的名称
// 	// 在同一父组下必须唯一
// 	// 示例: "administrators", "developers", "sales-team"
// 	Name *string `json:"name,omitempty" gorm:"type:varchar(255);not null"`
//
// 	// Path 是组的完整路径
// 	// 从根组到当前组的完整路径，使用 "/" 分隔
// 	// 示例: "/company/department/team", "/administrators"
// 	Path *string `json:"path,omitempty" gorm:"type:varchar(512);uniqueIndex"`
//
// 	// SubGroups 是子组列表
// 	// 包含当前组下的所有直接子组
// 	// 形成树形结构，支持多级组织架构
// 	SubGroups datatypes.JSONType[[]Group] `json:"sub_groups,omitempty" gorm:"type:json"`
//
// 	// Attributes 是组的自定义属性映射
// 	// 可存储任意额外信息，键是属性名，值是字符串数组（支持多值属性）
// 	// 组成员会继承这些属性
// 	// 示例: {"department": ["IT"], "location": ["Beijing", "Shanghai"], "cost_center": ["CC001"]}
// 	Attributes datatypes.JSONType[map[string][]string] `json:"attributes,omitempty" gorm:"type:json"`
//
// 	// Access 是组的访问权限映射
// 	// 定义对组的管理权限
// 	// 示例: {"view": true, "manage": true, "manageMembership": true, "viewMembers": true}
// 	Access datatypes.JSONType[map[string]bool] `json:"access,omitempty" gorm:"type:json"`
//
// 	// ClientRoles 是组在特定客户端应用中的角色映射
// 	// 键是客户端ID，值是该客户端中的角色列表
// 	// 组成员会自动获得这些客户端角色
// 	// 示例: {"inventory-app": ["viewer", "editor"], "hr-system": ["employee"]}
// 	ClientRoles datatypes.JSONType[map[string][]string] `json:"client_roles,omitempty" gorm:"type:json"`
//
// 	// RealmRoles 是组在 realm 级别拥有的角色列表
// 	// 这些是全局角色，对整个 realm 中的所有客户端都有效
// 	// 组成员会自动获得这些 realm 角色
// 	// 示例: ["user", "employee", "manager"]
// 	RealmRoles datatypes.JSONType[[]string] `json:"realm_roles,omitempty" gorm:"type:json"`
//
// 	// ParentID 父组的ID
// 	// 为空表示是顶级组
// 	ParentID *string `json:"parent_id,omitempty" gorm:"type:varchar(36);index"`
//
// 	model.Base
// }

type Group struct {
	// 常用索引字段
	Name *string `json:"name"`
	Path *string `json:"path"` // 保留这个因为 Path 必须唯一

	// 层级关系字段
	ParentID string `json:"parent_id"`
	Level    int    `json:"level"` // 层级深度，便于查询

	RealmId   *string `json:"realm_id,omitempty"`
	RealmName *string `json:"realm_name,omitempty"`

	// 完整的 Group 配置存储在 JSON 字段中
	Data datatypes.JSON `json:"data,omitempty"`

	model.Base
}

// GroupData 完整的组数据结构
type GroupData struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Path        string              `json:"path"`
	SubGroups   []GroupData         `json:"sub_groups"`
	Attributes  map[string][]string `json:"attributes"`
	Access      map[string]bool     `json:"access"`
	ClientRoles map[string][]string `json:"client_roles"`
	RealmRoles  []string            `json:"realm_roles"`
}

func (*Group) TableName() string {
	return "keycloak_groups"
}
