package model

import (
	. "github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
)

type User5 struct {
	Name string
	Addr string

	model.Base
}
