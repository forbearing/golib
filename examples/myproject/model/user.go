package model

import pkgmodel "github.com/forbearing/golib/model"

func init() {
	pkgmodel.Register(users...)
}

var users []*User = []*User{
	{Name: "user01", Email: "user01@example.com", Base: pkgmodel.Base{ID: "user01"}},
	{Name: "user02", Email: "user02@example.com", Base: pkgmodel.Base{ID: "user02"}},
	{Name: "user03", Email: "user03@example.com", Base: pkgmodel.Base{ID: "user03"}},
}

type User struct {
	Name   string `json:"name,omitempty" schema:"name" gorm:"unique"`
	Email  string `json:"email,omitempty" schema:"email" gorm:"unique"`
	Avatar string `json:"avatar,omitempty" schema:"avatar"`

	pkgmodel.Base
}
