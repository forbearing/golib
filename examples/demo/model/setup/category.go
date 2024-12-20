package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Category]()
}

type Category struct {
	pkgmodel.Base

	Name string `json:"name"`
}
