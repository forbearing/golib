package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Cluster]()
}

type Cluster struct {
	pkgmodel.Base

	Name string `json:"name"`
}
