package model

import "github.com/forbearing/golib/authz/rbac"

type Permission struct {
	RoleId   string `json:"role_id,omitempty" schema:"role_id"`
	Resource string `json:"resource,omitempty" schema:"resource"`
	Action   string `json:"action,omitempty" schema:"action"`

	Base
}

func (p *Permission) CreateAfter() error {
	return rbac.RBAC().GrantPermission(p.RoleId, p.Resource, p.Action)
}

func (p *Permission) DeleteAfter() error {
	return rbac.RBAC().RevokePermission(p.RoleId, p.Resource, p.Action)
}
