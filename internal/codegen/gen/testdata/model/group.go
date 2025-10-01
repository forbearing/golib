package model

import "github.com/forbearing/gst/model"

type Group struct {
	Name     string
	NumUsers int

	model.Base
}
