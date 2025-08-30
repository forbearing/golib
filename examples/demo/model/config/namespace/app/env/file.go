package env

import (
	. "github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
)

type File struct {
	model.Base
}

func (File) Design() {
	Endpoint("files") // 改成复数
	Param("filename")

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
