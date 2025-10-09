package model

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/gst/types"
	"github.com/forbearing/gst/util"
	"github.com/gertd/go-pluralize"
	"github.com/stoewer/go-strcase"
)

var pluralizeCli = pluralize.NewClient()

func (gs GormStrings) Value() (driver.Value, error) {
	// It will return "", if gs is nil or empty string.
	return strings.Trim(strings.Join(gs, ","), ","), nil
}

// GormScannerWrapper converts object to GormScanner that can be used in GORM.
// WARN: you must pass pointer to object.
func GormScannerWrapper(object any) *GormScanner {
	return &GormScanner{Object: object}
}

type GormScanner struct {
	Object any
}

func (g *GormScanner) Scan(value any) (err error) {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case string:
		err = json.Unmarshal(util.StringToBytes(v), g.Object)
	case []byte:
		err = json.Unmarshal(v, g.Object)
	default:
		err = errors.New("unsupported type, expected string or []byte")
	}
	return err
}

func (g *GormScanner) Value() (driver.Value, error) {
	data, err := json.Marshal(g.Object)
	if err != nil {
		return nil, err
	}
	return util.BytesToString(data), nil
}

func GetTableName[M types.Model]() string {
	return strcase.SnakeCase(pluralizeCli.Plural(reflect.TypeOf(*new(M)).Elem().Name()))
}

// AreTypesEqual checks if the types of M, REQ and RSP are equal
// If the M is a struct only has field model.Empty, always return false.
func AreTypesEqual[M types.Model, REQ types.Request, RSP types.Response]() bool {
	if IsModelEmpty[M]() {
		return false
	}
	typ1 := reflect.TypeFor[M]()
	typ2 := reflect.TypeFor[REQ]()
	typ3 := reflect.TypeFor[RSP]()
	return typ1 == typ2 && typ2 == typ3
}

// IsModelEmpty check the REQ is struct only has anonymous field model.Empty or has no fields, eg:
//
//	type Login struct {
//		model.Empty
//	}
//
//	type Logout struct{
//	}
func IsModelEmpty[T any]() bool {
	typ := reflect.TypeFor[T]()

	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return false
	}
	if typ.NumField() == 0 {
		return true
	}
	if typ.NumField() == 1 {
		field := typ.Field(0)
		target := reflect.TypeFor[Empty]()
		return field.Anonymous && field.Type == target
	}

	return false
}
