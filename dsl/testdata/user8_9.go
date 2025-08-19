package model

import (
	. "github.com/forbearing/golib/dsl"
	pkgmodel "github.com/forbearing/golib/model"
)

type User8 struct {
	Name string
}

type User8 struct {
	Name string

	pkgmodel.Empty
}

func (*User8) Design() {
	Migrate(true)
}
