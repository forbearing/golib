package dsl

import (
	"slices"
	"strings"
)

func Create(func())             {}
func Delete(func())             {}
func Update(func())             {}
func UpdatePartial(func())      {}
func Get(func())                {}
func List(func())               {}
func BatchCreate(func())        {}
func BatchDelete(func())        {}
func BatchUpdate(func())        {}
func BatchUpdatePartial(func()) {}

func Payload[T any]() {}
func Result[T any]()  {}

func Endpoint(string) {}
func Enabled(bool)    {}

type Design struct {
	Enabled  bool
	Endpoint string

	Create             *Action
	Update             *Action
	UpdatePartial      *Action
	BatchCreate        *Action
	BatchUpdate        *Action
	BatchUpdatePartial *Action
}

type Action struct {
	Payload string
	Result  string
}

var methodList = []string{
	"Create",
	"Delete",
	"Update",
	"UpdatePartial",
	"Get",
	"List",
	"BatchCreate",
	"BatchDelete",
	"BatchUpdate",
	"BatchUpdatePartial",

	"Enabled",
	"Endpoint",
	"Payload",
	"Result",
}

func is(name string) bool {
	return slices.Contains(methodList, name)
}

// trimQuote trim "str" two side quote: -"-, -'-, -`-
func trimQuote(str string) string {
	return strings.TrimFunc(str, func(r rune) bool {
		return r == '`' || r == '"' || r == '\''
	})
}
