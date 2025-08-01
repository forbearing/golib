package dsl

const input1 = `
package model

import (
	. "github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
)

type User struct {
	Name string
	Addr string

	model.Base
}

func (User) Design() {
	// Default to true.
	//Enabled(true)

	// Default Endpoint is the lower case of the model name.
	Endpoint("user2")

	// Custom create action request "Payload" and response "Result".
	Create(func() {
		Payload[User]()
		Result[*User]()
	})

	// Custom update action request "Payload" and response "Result".
	Update(func() {
		Payload[*User]()
		Result[User]()
	})
}
	`

const input2 = `
package model

import (
	"github.com/forbearing/golib/dsl"
	pkgmodel "github.com/forbearing/golib/model"
)

type User2 struct {
	Name string
	Addr string

	pkgmodel.Base
}

func (User2) Design() {
	// Default to true.
	dsl.Enabled(false)

	// Default Endpoint is the lower case of the model name.
	// dsl.Endpoint("user")

	// Custom create action request "Payload" and response "Result".
	dsl.Create(func() {
		dsl.Payload[User2]()
		dsl.Result[*User3]()
	})

	// Custom update partial action request "Payload" and response "Result".
	dsl.UpdatePartial(func() {
		dsl.Payload[*User]()
		dsl.Result[User]()
	})

	// Invalid design.
	dsl.UpdatePartial2(func() {
		dsl.Payload[*User]()
		dsl.Result[User]()
	})
}
`

const input3 = `
package model

import (
	"github.com/forbearing/golib/dsl"
	. "github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
	pkgmodel "github.com/forbearing/golib/model"
)

type User3 struct {
	Name string
	Addr string

	model.Base
}

func (User3) Design() {
	// Default to true.
	Enabled(true)

	// Default Endpoint is the lower case of the model name.
	Endpoint("user")

	// Custom create action request "Payload" and response "Result".
	Create(func() {
		Payload[User]()
		Result[*User]()
	})

	// Custom update action request "Payload" and response "Result".
	Update(func() {
		Payload[*User]()
		Result[User]()
	})
}

type User4 struct {
	Name string
	Addr string

	pkgmodel.Base
}

func (User4) Design() {
	// Default to true.
	dsl.Enabled(true)

	// Default Endpoint is the lower case of the model name.
	dsl.Endpoint("user")

	// Custom create action request "Payload" and response "Result".
	dsl.Create(func() {
		dsl.Payload[User]()
		dsl.Result[*User]()
	})

	// Custom update action request "Payload" and response "Result".
	dsl.Update(func() {
		dsl.Payload[*User]()
		dsl.Result[User]()
	})
}

	`

const input4 = `
package model

type User2 struct {
	Name string
	Addr string
}`
