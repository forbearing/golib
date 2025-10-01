package model

import (
	"github.com/forbearing/gst/dsl"
	pkgmodel "github.com/forbearing/gst/model"
)

type User2 struct {
	Name string
	Addr string

	pkgmodel.Base
}

func (User2) Design() {
	// Default to true.
	dsl.Enabled(false)
	dsl.Param("{user}")

	// Default Endpoint is the lower case of the model name.
	// dsl.Endpoint("user")

	// Custom create action request "Payload" and response "Result".
	dsl.Create(func() {
		dsl.Payload[User2]()
		dsl.Result[*User3]()
	})

	// Custom update partial action request "Payload" and response "Result".
	dsl.Patch(func() {
		dsl.Enabled(true)
		dsl.Payload[*User]()
		dsl.Result[User]()
	})

	// Invalid design.
	dsl.Patch2(func() {
		dsl.Enabled(false)
		dsl.Payload[*User]()
		dsl.Result[User]()
	})
}
