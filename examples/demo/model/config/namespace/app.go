package namespace

import (
	. "github.com/forbearing/gst/dsl"
	"github.com/forbearing/gst/model"
)

type App struct {
	model.Base
}

type AppRsp struct {
	Param string `json:"param"`
}

func (App) Design() {
	Endpoint("apps") // 改成复数
	Param("app")

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
		Result[*AppRsp]()
	})

	Get(func() {
		Enabled(true)
		Service(false)
	})
}
