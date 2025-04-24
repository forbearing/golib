package model_authz

import (
	"fmt"

	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap/zapcore"
)

func init() {
	model.Register[*Permission]()
}

type Permission struct {
	Resource string `json:"resource,omitempty" schema:"resource"`
	Action   string `json:"action,omitempty" schema:"action"`

	model.Base
}

func (p *Permission) CreateBefore() error {
	p.SetID(util.HashID(p.Resource, p.Action))
	p.Remark = util.ValueOf(fmt.Sprintf("%s %s", p.Action, p.Resource))
	return nil
}

func (p *Permission) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if p == nil {
		return nil
	}
	enc.AddString("resource", p.Resource)
	enc.AddString("action", p.Action)
	enc.AddObject("base", &p.Base)
	return nil
}
