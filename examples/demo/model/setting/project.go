package setting

import (
	. "github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
)

type Project struct {
	Name string `json:"name,omitempty" schema:"name"`

	model.Base
}

func (Project) Design() {
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
