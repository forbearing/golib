package model

import "github.com/forbearing/golib/model"

func init() {
	// create table `groups` automatically.
	// Ensure the package `model` will imported directly or indirectly in `main.go`.
	model.Register[*Group]()
}

type Group struct {
	model.Base

	Name *string `json:"name,omitempty"`
}
