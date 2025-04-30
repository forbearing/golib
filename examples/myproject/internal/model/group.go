package model

import (
	pkgmodel "github.com/forbearing/golib/model"
	"go.uber.org/zap/zapcore"
)

func init() {
	pkgmodel.Register[*Group]()
}

type Group struct {
	Name        string `json:"name,omitempty" schema:"name"`
	Desc        string `json:"desc,omitempty" schema:"desc"`
	MemberCount int    `json:"member_count" gorm:"default:0"`

	pkgmodel.Base
}

func (g *Group) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if g == nil {
		return nil
	}
	enc.AddString("name", g.Name)
	enc.AddString("desc", g.Desc)
	enc.AddInt("member_count", g.MemberCount)
	enc.AddObject("base", &g.Base)
	return nil
}

// func (g *Group) Request(GroupRequest, GroupResponse) {}

type GroupRequest struct {
	Name string
}

type GroupResponse struct {
	FieldC string
	FieldD *int
}
