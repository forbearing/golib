package model_authz

import (
	"github.com/forbearing/golib/model"
	"go.uber.org/zap/zapcore"
)

type Permission struct {
	Resource string `json:"resource,omitempty" schema:"resource"`
	Action   string `json:"action,omitempty" schema:"action"`

	model.Base
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
