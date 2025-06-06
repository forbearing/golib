package model_authz

import (
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authz/rbac"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"go.uber.org/zap/zapcore"
)

func init() {
	model.Register[*Role]()
}

type Role struct {
	Name string `json:"name,omitempty" schema:"name"`

	model.Base
}

func (r *Role) CreateBefore() error {
	if len(strings.TrimSpace(r.Name)) == 0 {
		return errors.New("name is required")
	}
	// Ensure the role with the same name share the same ID.
	// If the role already exists, set same id to just update it.
	r.SetID(util.HashID(r.Name))

	return nil
}
func (r *Role) UpdateBefore() error { return r.CreateAfter() }

func (r *Role) CreateAfter() error { return rbac.RBAC().AddRole(r.Name) }
func (r *Role) DeleteBefore() error {
	// The delete request always don't have role id, so we should get the role from database.
	if err := database.Database[*Role]().Get(r, r.ID); err != nil {
		return err
	}
	if len(r.Name) > 0 {
		return rbac.RBAC().RemoveRole(r.Name)
	}
	return nil
}
func (r *Role) DeleteAfter() error { return database.Database[*Role]().Cleanup() }

func (r *Role) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if r == nil {
		return nil
	}
	enc.AddString("name", r.Name)
	enc.AddObject("base", &r.Base)
	return nil
}
