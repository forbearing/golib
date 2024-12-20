package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Certificate]()
}

type Certificate struct {
	pkgmodel.Base

	Name string `json:"name"`
}
