package model

import (
	. "github.com/forbearing/gst/dsl"
	pkgmodel "github.com/forbearing/gst/model"
)

type User8 struct {
	Name string

	pkgmodel.Empty
}

func (*User8) Design() {
	Migrate(true)
}

type User9 struct {
	Name string
}
