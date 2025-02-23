package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*pkgmodel.SysInfo]()
	pkgmodel.Register[*pkgmodel.Session]()
}
