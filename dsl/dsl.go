package dsl

import (
	"slices"
	"strings"

	"github.com/forbearing/golib/types/consts"
)

func Enabled(bool)    {}
func Endpoint(string) {}
func Payload[T any]() {}
func Result[T any]()  {}

func Create(func()) {}
func Delete(func()) {}
func Update(func()) {}
func Patch(func())  {}
func List(func())   {}
func Get(func())    {}

func CreateMany(func()) {}
func DeleteMany(func()) {}
func UpdateMany(func()) {}
func PatchMany(func())  {}

type Design struct {
	Enabled  bool
	Endpoint string

	Create *Action
	Delete *Action
	Update *Action
	Patch  *Action
	List   *Action
	Get    *Action

	CreateMany *Action
	DeleteMany *Action
	UpdateMany *Action
	PatchMany  *Action
}

type Action struct {
	Payload string
	Result  string
}

var methodList = []string{
	"Enabled",
	"Endpoint",
	"Payload",
	"Result",

	consts.PHASE_CREATE.MethodName(),
	consts.PHASE_DELETE.MethodName(),
	consts.PHASE_UPDATE.MethodName(),
	consts.PHASE_PATCH.MethodName(),
	consts.PHASE_LIST.MethodName(),
	consts.PHASE_GET.MethodName(),

	consts.PHASE_CREATE_MANY.MethodName(),
	consts.PHASE_DELETE_MANY.MethodName(),
	consts.PHASE_UPDATE_MANY.MethodName(),
	consts.PHASE_PATCH_MANY.MethodName(),
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
