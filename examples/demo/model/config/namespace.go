package config

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
)

type Namespace struct {
	model.Base
}

func (Namespace) Design() {
	Endpoint("namespaces") // 改成复数
	Param("ns")

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
