package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Group]()
}

type Group struct {
	Name string `json:"name,omitempty" schema:"name" gorm:"unique" binding:"required"`

	pkgmodel.Base
}
