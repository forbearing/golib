package setting

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
)

type Tenant struct {
	Name string `json:"name,omitempty" schema:"name"`

	model.Base
}

func (Tenant) Design() {
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
