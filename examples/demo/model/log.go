package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register[*Log]()
}

type Log struct {
	pkgmodel.Base

	Level   string `json:"level"`
	Message string `json:"message"`
}
