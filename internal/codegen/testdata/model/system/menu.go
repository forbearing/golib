package model_system

import "github.com/forbearing/golib/model"

var menuRoot = &Menu{ParentId: "root", Base: model.Base{ID: "root"}}

func init() {
	model.Register[*Menu]([]*Menu{menuRoot}...)
}

type Menu struct {
	Path    string `json:"path,omitempty" schema:"path"`
	Element string `json:"element,omitempty" schema:"element"`
	Label   string `json:"label,omitempty" schema:"label"`
	Icon    string `json:"icon,omitempty" schema:"icon"`

	Visiable *bool  `json:"visiable" schema:"visiable"`
	Default  string `json:"default,omitempty" schema:"default"`
	Status   *uint  `json:"status" gorm:"type:smallint;default:1;comment:status(0: disabled, 1: enabled)" schema:"status"`

	ParentId string  `json:"parent_id,omitempty" gorm:"size:191" schema:"parent_id"`
	Parent   *Menu   `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"` // 父路由
	Children []*Menu `json:"children,omitempty" gorm:"foreignKey:ParentId"`             // 子路由

	model.Base
}
