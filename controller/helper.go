package controller

import (
	"fmt"
	"reflect"

	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

func CreateSession(c *gin.Context) *model.Session {
	ua := useragent.New(c.Request.UserAgent())
	engineName, engineVersion := ua.Engine()
	browserName, browserVersion := ua.Browser()
	return &model.Session{
		UserId:         c.GetString(consts.CTX_USER_ID),
		Username:       c.GetString(consts.CTX_USERNAME),
		Platform:       ua.Platform(),
		OS:             ua.OS(),
		EngineName:     engineName,
		EngineVersion:  engineVersion,
		BrowserName:    browserName,
		BrowserVersion: browserVersion,
	}
}

func partialUpdateValue(log types.Logger, typ reflect.Type, oldVal reflect.Value, newVal reflect.Value) {
	for i := range typ.NumField() {
		// fmt.Println(typ.Field(i).Name, typ.Field(i).Type, typ.Field(i).Type.Kind(), newVal.Field(i).IsValid(), newVal.Field(i).CanSet())
		switch typ.Field(i).Type.Kind() {
		case reflect.Struct: // skip update base model.
			switch typ.Field(i).Type.Name() {
			case "GormTime": // The underlying type of model.GormTime(type of time.Time) is struct, we should continue handle.

			case "Base":
				// 有些结构体会匿名继承其他的结构体，例如 AssetChecking 匿名继承 Asset, 所以要可以额外检查是不是某个匿名结构体.
				// 可以自动深度查找,不需要链式查找, 例如
				// newVal.FieldByName("Asset").FieldByName("Remark").IsValid() 可以简化为
				// newVal.FieldByName("Remark").IsValid()

				// Make sure the type of "Remark" is pointer to golang base type.
				fieldRemark := "Remark"
				if oldVal.FieldByName(fieldRemark).CanSet() {
					if newVal.FieldByName(fieldRemark).IsValid() { // WARN: oldVal.FieldByName(fieldRemark) maybe <invalid reflect.Value>
						if !newVal.FieldByName(fieldRemark).IsZero() {
							// output log must before set value.
							if newVal.FieldByName(fieldRemark).Kind() == reflect.Pointer {
								log.Info(fmt.Sprintf("[UpdatePartial %s] field: %q: %v --> %v", fieldRemark, typ.Name(),
									oldVal.FieldByName(fieldRemark).Elem(), newVal.FieldByName(fieldRemark).Elem())) // WARN: you shouldn't call oldVal.FieldByName(fieldRemark).Elem().Interface()
							} else {
								log.Info(fmt.Sprintf("[UpdatePartial %s] field: %q: %v --> %v", fieldRemark, typ.Name(),
									oldVal.FieldByName(fieldRemark).Interface(), newVal.FieldByName(fieldRemark).Interface()))
							}
							oldVal.FieldByName(fieldRemark).Set(newVal.FieldByName(fieldRemark)) // set old value by new value
						}
					}
				}
				// Make sure the type of "Order" is pointer to golang base type.
				fieldOrder := "Order"
				if oldVal.FieldByName(fieldOrder).CanSet() {
					if newVal.FieldByName(fieldOrder).IsValid() { // WARN: oldVal.FieldByName(fieldOrder) maybe <invalid reflect.Value>
						if !newVal.FieldByName(fieldOrder).IsZero() {
							// output log must before set value.
							if newVal.FieldByName(fieldOrder).Kind() == reflect.Pointer {
								log.Info(fmt.Sprintf("[UpdatePartial %s] field: %q: %v --> %v", fieldOrder, typ.Name(),
									oldVal.FieldByName(fieldOrder).Elem(), newVal.FieldByName(fieldOrder).Elem())) // WARN: you shouldn't call oldVal.FieldByName(fieldOrder).Elem().Interface()
							} else {
								log.Info(fmt.Sprintf("[UpdatePartial %s] field: %q: %v --> %v", fieldOrder, typ.Name(),
									oldVal.FieldByName(fieldOrder).Interface(), newVal.FieldByName(fieldOrder).Interface()))
							}
							oldVal.FieldByName(fieldOrder).Set(newVal.FieldByName(fieldOrder)) // set old value by new value.
						}
					}
				}
				continue

			default:
				continue
			}
		}
		if !oldVal.Field(i).CanSet() {
			log.Warnf("field %q is cannot set, skip", typ.Field(i).Name)
			continue
		}
		if !newVal.Field(i).IsValid() {
			// log.Warnf("field %s is invalid, skip", typ.Field(i).Name)
			continue
		}
		// base type such like int and string have default value(zero value).
		// If the struct field(the field type is golang base type) supported by patch update,
		// the field type must be pointer to base type, such like *string, *int.
		if newVal.Field(i).IsZero() {
			// log.Warnf("field %s is zero value, skip", typ.Field(i).Name)
			// log.Warnf("DeepEqual: %v : %v : %v : %v", typ.Field(i).Name, newVal.Field(i).Interface(), oldVal.Field(i).Interface(), reflect.DeepEqual(newVal.Field(i), oldVal.Field(i)))
			continue
		}
		// output log must before set value.
		if newVal.Field(i).Kind() == reflect.Pointer {
			log.Info(fmt.Sprintf("[UpdatePartial %s] field: %q: %v --> %v", typ.Name(), typ.Field(i).Name, oldVal.Field(i).Elem().Interface(), newVal.Field(i).Elem().Interface()))
		} else {
			log.Info(fmt.Sprintf("[UpdatePartial %s] field: %q: %v --> %v", typ.Name(), typ.Field(i).Name, oldVal.Field(i).Interface(), newVal.Field(i).Interface()))
		}
		oldVal.Field(i).Set(newVal.Field(i)) // set old value by new value
	}
}
