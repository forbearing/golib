package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Database]()
}

type Database struct {
	pkgmodel.Base

	Name string `json:"name"`
}
