package model

import (
	pkgmodel "github.com/forbearing/golib/model"
	"go.uber.org/zap/zapcore"
)

func init() {
	pkgmodel.Register[*Department]()
}

type Department struct {
	Name string `json:"name,omitempty" schema:"name" gorm:"unique" binding:"required"`
	Desc string `json:"desc,omitempty" schema:"desc"`

	pkgmodel.Base
}

func (g *Department) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if g == nil {
		return nil
	}
	enc.AddString("name", g.Name)
	enc.AddString("desc", g.Desc)
	enc.AddObject("base", &g.Base)
	return nil
}
