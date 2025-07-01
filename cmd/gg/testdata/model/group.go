package model

import "github.com/forbearing/golib/model"

type Group struct {
	Name        string `json:"name,omitempty" schema:"name"`
	Desc        string `json:"desc,omitempty" schema:"desc"`
	MemberCount int    `json:"member_count" gorm:"default:0"`

	model.Base
}
