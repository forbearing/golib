package env

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
)

type Item struct {
	model.Base
}

type ItemRsp struct {
	Ns  string
	App string
	Env string
}

func (Item) Design() {
	Endpoint("items") // 改成复数
	Param("key")

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
		Service(true)
		Result[*ItemRsp]()
	})

	Get(func() {
		Enabled(true)
		Service(false)
	})
}
