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
