package env

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
)

type File struct {
	model.Base
}

func (File) Design() {
	Endpoint("files") // 改成复数
	Param("file")

	Route("/config/files", func() {
		Create(func() {
			Enabled(true)
		})
		Delete(func() {
			Enabled(true)
		})
		Update(func() {
			Enabled(true)
		})
		Patch(func() {
			Enabled(true)
		})
		List(func() {
			Enabled(true)
			Service(true)
		})
		Get(func() {
			Enabled(true)
		})
	})

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
