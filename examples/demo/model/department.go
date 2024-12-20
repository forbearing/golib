package model

import "github.com/forbearing/golib/model"

func init() {
	// create table `departments` automatically.
	// Ensure the package `model` will imported directly or indirectly in `main.go`.
	model.Register[*Department]()
}

type Department struct {
	model.Base

	Name   *string `json:"name,omitempty"`
	Leader *string `json:"leader,omitempty"`
}
