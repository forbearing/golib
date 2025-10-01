package app

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
)

type Env struct {
	model.Base
}

func (Env) Design() {
	Endpoint("envs") // 改成复数
	Param("env")

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
