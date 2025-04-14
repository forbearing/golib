package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Star]()
}

type Star struct {
	pkgmodel.Base
}
