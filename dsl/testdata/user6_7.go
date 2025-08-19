package model

import (
	"github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
)

type User6 struct {
	Name string

	model.Empty
}

func (User6) Design() {
	dsl.Migrate(true)
}

type User7 struct {
	Name string
}
