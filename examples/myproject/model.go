package main

import "github.com/forbearing/golib/model"

func init() {
	model.Register[*User]()
	model.Register[*Group]()
	model.Register[*model.SysInfo]()
}

type User struct {
	Name   string `json:"name,omitempty" schema:"name" gorm:"unique"`
	Email  string `json:"email,omitempty" schema:"email" gorm:"unique"`
	Avatar string `json:"avatar,omitempty" schema:"avatar"`

	model.Base
}

type Group struct {
	Name string `json:"name,omitempty" schema:"name" gorm:"unique" binding:"required"`

	model.Base
}
