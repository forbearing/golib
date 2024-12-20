package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Area]()
}

type Area struct {
	pkgmodel.Base
	Name string `json:"name"`
}
