package model

import "github.com/forbearing/golib/model"

type User struct {
	Name string `json:"name,omitempty"`
	Age  uint8  `json:"age,omitempty"`

	model.Base
}
