// filename: model/user.go
package model

import "github.com/forbearing/golib/model"

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`

	model.Base
}

type Group struct {
	Name    string `json:"name"`
	Members []User `json:"members"`

	model.Base
}

type GroupUser struct {
	GroupId int `json:"group_id"`
	UserId  int `json:"user_id"`
}
