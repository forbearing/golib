package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*User]()
}

type User struct {
	pkgmodel.Base

	Name string `json:"name"`
}
