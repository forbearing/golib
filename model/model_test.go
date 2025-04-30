package model_test

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/forbearing/golib/model"
	"github.com/stretchr/testify/assert"
)

type (
	User  struct{ model.Base }
	User2 struct{ model.Base }
	User3 struct{ model.Base }
	User4 struct{ model.Base }
	User5 struct{ model.Base }
	User6 struct{ model.Base }
	User7 struct{ model.Base }
	User8 struct{ model.Base }
)

func (*User2) Request(Request1)             {}
func (*User3) Request(Request2, struct{})   {}
func (*User4) Request1(struct{}, struct{})  {}
func (*User5) Request(int, struct{})        {}
func (*User6) Request(struct{}, int)        {}
func (*User7) Request(*Request1)            {}
func (*User8) Request(*Request1, *Response) {}

type Request1 struct {
	Field1 string
	Field2 *string
}

type Request2 struct {
	Field3 *int
	Field4 struct{}
}

type Response struct {
	Field5 *int
	Field6 struct{}
}

func TestHasRequest(t *testing.T) {
	assert.Equal(t, false, model.HasRequest[*User]())
	assert.Equal(t, false, model.HasResponse[*User]())

	assert.Equal(t, true, model.HasRequest[*User2]())
	assert.Equal(t, false, model.HasResponse[*User2]())

	assert.Equal(t, true, model.HasRequest[*User3]())
	assert.Equal(t, true, model.HasResponse[*User3]())

	assert.Equal(t, false, model.HasRequest[*User4]())
	assert.Equal(t, false, model.HasResponse[*User4]())

	assert.Equal(t, false, model.HasRequest[*User5]())
	assert.Equal(t, true, model.HasResponse[*User5]())

	assert.Equal(t, true, model.HasRequest[*User6]())
	assert.Equal(t, false, model.HasResponse[*User6]())

	assert.Equal(t, true, model.HasRequest[*User7]())
	assert.Equal(t, false, model.HasResponse[*User7]())

	assert.Equal(t, true, model.HasRequest[*User8]())
	assert.Equal(t, true, model.HasResponse[*User8]())
}

func TestNewRequest(t *testing.T) {
	value := model.NewRequest[*User]()
	spew.Dump(value)
}
