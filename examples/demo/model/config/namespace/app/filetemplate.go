package app

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
)

type FileTemplate struct {
	model.Base
}

func (FileTemplate) Design() {
	Endpoint("file-templates") // 改成复数

	Create(func() {
		Enabled(true)
		Service(false)
	})

	Delete(func() {
		Enabled(true)
		Service(false)
	})

	Update(func() {
		Enabled(true)
		Service(false)
	})

	Patch(func() {
		Enabled(true)
		Service(false)
	})

	List(func() {
		Enabled(true)
		Service(false)
	})

	Get(func() {
		Enabled(true)
		Service(false)
	})
}
