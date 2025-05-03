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

func (g *Group) Request(GroupRequest, GroupResponse) {}

type GroupRequest struct {
	Name string `json:"name,omitempty" schema:"name"`
}

type GroupResponse struct {
	CustomName string  `json:"custom_name,omitempty" schema:"custom_name"`
	CustomDesc *string `json:"custom_desc,omitempty" schema:"custom_desc"`
}
