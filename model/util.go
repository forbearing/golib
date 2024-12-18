package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/util"
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

func GetTablename[M types.Model]() string {
	return strcase.LowerCamelCase(pluralizeCli.Plural(reflect.TypeOf(*new(M)).Elem().Name()))
}
