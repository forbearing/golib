package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Contact]()
}

type Contact struct {
	pkgmodel.Base

	Name string `json:"name"`
}
