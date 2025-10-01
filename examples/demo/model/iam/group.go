package iam

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
	"gorm.io/datatypes"
)

type Group struct {
	Name string

	Path *string `json:"path"` // 保留这个因为 Path 必须唯一

	// 层级关系字段
	ParentID string `json:"parent_id"`
	Level    int    `json:"level"` // 层级深度，便于查询

	TenantId   *string `json:"tenant_id,omitempty"`
	TenantName *string `json:"tenant_name,omitempty"`

	// 完整的 Group 配置存储在 JSON 字段中
	Data datatypes.JSON `json:"data,omitempty"`

	model.Base
}

func (Group) Design() {
	Migrate(true)

	List(func() {
		Enabled(true)
		Service(false)
	})
	Get(func() {
		Enabled(true)
		Service(false)
	})
}
