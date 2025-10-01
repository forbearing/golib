package iam

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
	"gorm.io/datatypes"
)

type User struct {
	Name      string                      `json:"name" schema:"name"`
	Enabled   *bool                       `json:"enabled" schema:"enabled"`
	FirstName *string                     `json:"first_name" schema:"first_name"`
	LastName  *string                     `json:"last_name" schema:"last_name"`
	Email     *string                     `json:"email" schema:"email"`
	Groups    datatypes.JSONSlice[string] `json:"groups" schema:"groups"`

	TenantID          *string `json:"tenant_id" schema:"tenant_id"`
	TenantName        *string `json:"tenant_name" schema:"tenant_name"`
	TenantDisplayName *string `json:"tenant_display_name" schema:"tenant_display_name"`

	Data datatypes.JSON `json:"data"`

	model.Base
}

func (User) Design() {
	Migrate(true)

	List(func() {
		Enabled(true)
		Service(false)
	})
	Get(func() {
		Enabled(true)
		Service(false)
	})
}
