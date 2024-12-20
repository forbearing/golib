package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Datacenter]()
}

type Datacenter struct {
	pkgmodel.Base

	Name string `json:"name"`
}
