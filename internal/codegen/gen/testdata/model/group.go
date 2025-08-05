package model

import "github.com/forbearing/golib/model"

type Group struct {
	Name     string
	NumUsers int

	model.Base
}
