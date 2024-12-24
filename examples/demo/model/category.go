package model

import pkgmodel "github.com/forbearing/golib/model"

type Category struct {
	pkgmodel.Base

	Name     string     `json:"name,omitempty" gorm:"unique" schema:"name"`
	Status   *uint      `json:"status,omitempty" gorm:"type:smallint;comment:status(0: disabled, 1: enabled)" schema:"status"`
	ParentId string     `json:"parent_id,omitempty" gorm:"size:191" schema:"parent_id"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentId"`
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentId;references:ID"`
}

func (*Category) Expands() []string {
	return []string{"Children", "Parent"}
}

func (*Category) Excludes() map[string][]any {
	return map[string][]any{pkgmodel.KeyName: {pkgmodel.RootName, pkgmodel.UnknownName, pkgmodel.NoneName}}
}

func (c *Category) CreateBefore() error {
	if len(c.ParentId) == 0 {
		c.ParentId = pkgmodel.RootId
	}
	return nil
}
