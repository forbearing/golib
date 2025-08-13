package dsl

import (
	"slices"
	"strings"

	"github.com/forbearing/golib/types/consts"
)

func Enabled(bool)    {}
func Endpoint(string) {}
func Migrate(string)  {}
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

func Import(func()) {}
func Export(func()) {}

type Design struct {
	Enabled  bool   // default enabled
	Endpoint string // Endpoint defaults to the lower case of the model name, its used by router.
	Migrate  bool   // Migrate to database or not, default to true.

	// default payload and result is the model name
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

	Import *Action
	Export *Action
}

type Action struct {
	Enabled bool   // defaults to false.
	Payload string // current Action Payload
	Result  string // current Action Result
}

var methodList = []string{
	"Enabled",
	"Endpoint",
	"Migrate",
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

	consts.PHASE_IMPORT.MethodName(),
	consts.PHASE_EXPORT.MethodName(),
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
