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
	// Enabled(true)

	// Default Endpoint is the lower case of the model name.
	Endpoint("user2")

	// Default to true,
	Migrate(true)
	Param("user")

	Route("/iam/users", func() {
		List(func() {
			Enabled(true)
			Service(false)
			Payload[*UserReq]()
			Result[*UserRsp]()
		})
		Get(func() {
			Enabled(true)
			Service(false)
		})
	})
	Route("///tenant/users", func() {
		Create(func() {
			Enabled(true)
			Payload[*UserReq]()
			Result[*User]()
		})
		Update(func() {
			Enabled(true)
		})
		Patch(func() {
			Enabled(true)
		})
		CreateMany(func() {
			Enabled(true)
		})
	})

	// Custom create action request "Payload" and response "Result".
	// Default payload and result is the model name.
	Create(func() {
		Enabled(true)
		Service(false)
		Public(true)
		Payload[User]()
		Result[*User]()
	})

	// Custom update action request "Payload" and response "Result".
	Update(func() {
		Enabled(false)
		Payload[*User]()
		Public(false)
		Result[User]()
	})

	Delete(func() {
		Enabled(true)
	})

	List(func() {
		Enabled(true)
	})
}
