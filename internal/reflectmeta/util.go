package reflectmeta

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stoewer/go-strcase"
	"go.uber.org/zap"
)

// StructFieldToMap2 extracts the field tags from a struct and writes them into a map.
// This map can then be used to build SQL query conditions.
// FIXME: if the field type is boolean or ineger, disable the fuzzy matching.
func StructFieldToMap2(_typ reflect.Type, val reflect.Value, q map[string]string) {
	if q == nil {
		q = make(map[string]string)
	}
	meta := GetStructMeta(_typ)
	for i := range meta.NumField() {
		if val.Field(i).IsZero() {
			continue
		}

		switch meta.Field(i).Type.Kind() {
		case reflect.Chan, reflect.Map, reflect.Func:
			continue
		case reflect.Struct:
			// All `model.XXX` extends the basic model named `Base`,
			if meta.Field(i).Name == "Base" {
				if !val.Field(i).FieldByName("CreatedBy").IsZero() {
					// Not overwrite the "CreatedBy" value set in types.Model.
					// The "CreatedBy" value set in types.Model has higher priority than base model.
					if _, loaded := q["created_by"]; !loaded {
						q["created_by"] = val.Field(i).FieldByName("CreatedBy").Interface().(string)
					}
				}
				if !val.Field(i).FieldByName("UpdatedBy").IsZero() {
					// Not overwrite the "UpdatedBy" value set in types.Model.
					// The "UpdatedBy" value set in types.Model has higher priority than base model.
					if _, loaded := q["updated_by"]; !loaded {
						q["updated_by"] = val.Field(i).FieldByName("UpdatedBy").Interface().(string)
					}
				}
				if !val.Field(i).FieldByName("ID").IsZero() {
					// Not overwrite the "ID" value set in types.Model.
					// The "ID" value set in types.Model has higher priority than base model.
					if _, loaded := q["id"]; !loaded {
						q["id"] = val.Field(i).FieldByName("ID").Interface().(string)
					}
				}
			} else {
				StructFieldToMap2(meta.Field(i).Type, val.Field(i), q)
			}
			continue
		}
		// "json" tag priority is higher than typ.Field(i).Name
		jsonTagStr := strings.TrimSpace(meta.JsonTag(i))
		jsonTagItems := strings.Split(jsonTagStr, ",")
		// NOTE: strings.Split always returns at least one element(empty string)
		// We should not use len(jsonTagItems) to check the json tags exists.
		jsonTag := ""
		if len(jsonTagItems) == 0 {
			// the structure lowercase field name as the query condition.
			jsonTagItems[0] = meta.Field(i).Name
		}
		jsonTag = jsonTagItems[0]
		if len(jsonTag) == 0 {
			// the structure lowercase field name as the query condition.
			jsonTag = meta.Field(i).Name
		}
		// "schema" tag have higher priority than "json" tag
		schemaTagStr := strings.TrimSpace(meta.SchemaTag(i))
		schemaTagItems := strings.Split(schemaTagStr, ",")
		schemaTag := ""
		if len(schemaTagItems) > 0 {
			schemaTag = schemaTagItems[0]
		}
		if len(schemaTag) > 0 && schemaTag != jsonTag {
			fmt.Printf("------------------ json tag replace by schema tag: %s -> %s\n", jsonTag, schemaTag)
			jsonTag = schemaTag
		}

		if !val.Field(i).CanInterface() {
			continue
		}
		v := val.Field(i).Interface()
		var _v string
		switch val.Field(i).Kind() {
		case reflect.Bool:
			// 由于 WHERE IN 语句会自动加上单引号,比如 WHERE `default` IN ('true')
			// 但是我们想要的是 WHERE `default` IN (true),
			// 所以没办法就只能直接转成 int 了.
			_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			_v = fmt.Sprintf("%d", v)
		case reflect.Float32, reflect.Float64:
			_v = fmt.Sprintf("%g", v)
		case reflect.String:
			_v = fmt.Sprintf("%s", v)
		case reflect.Pointer:
			v = val.Field(i).Elem().Interface()
			// switch typ.Elem().Kind() {
			switch val.Field(i).Elem().Kind() {
			case reflect.Bool:
				_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				_v = fmt.Sprintf("%d", v)
			case reflect.Float32, reflect.Float64:
				_v = fmt.Sprintf("%g", v)
			case reflect.String:
				_v = fmt.Sprintf("%s", v)
			case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
			default:
				_v = fmt.Sprintf("%v", v)
			}
		case reflect.Slice:
			_len := val.Field(i).Len()
			if _len == 0 {
				zap.S().Warn("reflect.Slice length is 0")
				_len = 1
			}
			slice := reflect.MakeSlice(val.Field(i).Type(), _len, _len)
			// fmt.Println("--------------- slice element", slice.Index(0), slice.Index(0).Kind(), slice.Index(0).Type())
			switch slice.Index(0).Kind() {
			case reflect.String: // handle string slice.
				// WARN: val.Field(i).Type() is model.GormStrings not []string,
				// execute statement `slice.Interface().([]string)` directly will case panic.
				// _v = strings.Join(slice.Interface().([]string), ",") // the slice type is GormStrings not []string.
				// We should make the slice of []string again.
				slice = reflect.MakeSlice(reflect.TypeOf([]string{}), _len, _len)
				reflect.Copy(slice, val.Field(i))
				_v = strings.Join(slice.Interface().([]string), ",")
			default:
				_v = fmt.Sprintf("%v", v)
			}
		case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
		default:
			_v = fmt.Sprintf("%v", v)
		}

		// json tag name naming format must be same as gorm table columns,
		// both should be "snake case" or "camel case".
		// gorm table columns naming format default to 'snake case',
		// so the json tag name is converted to "snake case here"
		// q[strcase.SnakeCase(jsonTag)] = val.Field(i).Interface()
		q[strcase.SnakeCase(jsonTag)] = _v
	}
}

// StructFieldToMap extracts the field tags from a struct and writes them into a map.
// This map can then be used to build SQL query conditions.
// FIXME: if the field type is boolean or ineger, disable the fuzzy matching.
func StructFieldToMap(typ reflect.Type, val reflect.Value, q map[string]string) {
	if q == nil {
		q = make(map[string]string)
	}
	for i := range typ.NumField() {
		if val.Field(i).IsZero() {
			continue
		}

		switch typ.Field(i).Type.Kind() {
		case reflect.Chan, reflect.Map, reflect.Func:
			continue
		case reflect.Struct:
			// All `model.XXX` extends the basic model named `Base`,
			if typ.Field(i).Name == "Base" {
				if !val.Field(i).FieldByName("CreatedBy").IsZero() {
					// Not overwrite the "CreatedBy" value set in types.Model.
					// The "CreatedBy" value set in types.Model has higher priority than base model.
					if _, loaded := q["created_by"]; !loaded {
						q["created_by"] = val.Field(i).FieldByName("CreatedBy").Interface().(string)
					}
				}
				if !val.Field(i).FieldByName("UpdatedBy").IsZero() {
					// Not overwrite the "UpdatedBy" value set in types.Model.
					// The "UpdatedBy" value set in types.Model has higher priority than base model.
					if _, loaded := q["updated_by"]; !loaded {
						q["updated_by"] = val.Field(i).FieldByName("UpdatedBy").Interface().(string)
					}
				}
				if !val.Field(i).FieldByName("ID").IsZero() {
					// Not overwrite the "ID" value set in types.Model.
					// The "ID" value set in types.Model has higher priority than base model.
					if _, loaded := q["id"]; !loaded {
						q["id"] = val.Field(i).FieldByName("ID").Interface().(string)
					}
				}
			} else {
				StructFieldToMap(typ.Field(i).Type, val.Field(i), q)
			}
			continue
		}
		// "json" tag priority is higher than typ.Field(i).Name
		jsonTagStr := strings.TrimSpace(typ.Field(i).Tag.Get("json"))
		jsonTagItems := strings.Split(jsonTagStr, ",")
		// NOTE: strings.Split always returns at least one element(empty string)
		// We should not use len(jsonTagItems) to check the json tags exists.
		jsonTag := ""
		if len(jsonTagItems) == 0 {
			// the structure lowercase field name as the query condition.
			jsonTagItems[0] = typ.Field(i).Name
		}
		jsonTag = jsonTagItems[0]
		if len(jsonTag) == 0 {
			// the structure lowercase field name as the query condition.
			jsonTag = typ.Field(i).Name
		}
		// "schema" tag have higher priority than "json" tag
		schemaTagStr := strings.TrimSpace(typ.Field(i).Tag.Get("schema"))
		schemaTagItems := strings.Split(schemaTagStr, ",")
		schemaTag := ""
		if len(schemaTagItems) > 0 {
			schemaTag = schemaTagItems[0]
		}
		if len(schemaTag) > 0 && schemaTag != jsonTag {
			fmt.Printf("------------------ json tag replace by schema tag: %s -> %s\n", jsonTag, schemaTag)
			jsonTag = schemaTag
		}

		if !val.Field(i).CanInterface() {
			continue
		}
		v := val.Field(i).Interface()
		var _v string
		switch val.Field(i).Kind() {
		case reflect.Bool:
			// 由于 WHERE IN 语句会自动加上单引号,比如 WHERE `default` IN ('true')
			// 但是我们想要的是 WHERE `default` IN (true),
			// 所以没办法就只能直接转成 int 了.
			_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			_v = fmt.Sprintf("%d", v)
		case reflect.Float32, reflect.Float64:
			_v = fmt.Sprintf("%g", v)
		case reflect.String:
			_v = fmt.Sprintf("%s", v)
		case reflect.Pointer:
			v = val.Field(i).Elem().Interface()
			// switch typ.Elem().Kind() {
			switch val.Field(i).Elem().Kind() {
			case reflect.Bool:
				_v = fmt.Sprintf("%d", boolToInt(v.(bool)))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				_v = fmt.Sprintf("%d", v)
			case reflect.Float32, reflect.Float64:
				_v = fmt.Sprintf("%g", v)
			case reflect.String:
				_v = fmt.Sprintf("%s", v)
			case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
			default:
				_v = fmt.Sprintf("%v", v)
			}
		case reflect.Slice:
			_len := val.Field(i).Len()
			if _len == 0 {
				zap.S().Warn("reflect.Slice length is 0")
				_len = 1
			}
			slice := reflect.MakeSlice(val.Field(i).Type(), _len, _len)
			// fmt.Println("--------------- slice element", slice.Index(0), slice.Index(0).Kind(), slice.Index(0).Type())
			switch slice.Index(0).Kind() {
			case reflect.String: // handle string slice.
				// WARN: val.Field(i).Type() is model.GormStrings not []string,
				// execute statement `slice.Interface().([]string)` directly will case panic.
				// _v = strings.Join(slice.Interface().([]string), ",") // the slice type is GormStrings not []string.
				// We should make the slice of []string again.
				slice = reflect.MakeSlice(reflect.TypeOf([]string{}), _len, _len)
				reflect.Copy(slice, val.Field(i))
				_v = strings.Join(slice.Interface().([]string), ",")
			default:
				_v = fmt.Sprintf("%v", v)
			}
		case reflect.Struct, reflect.Map, reflect.Chan, reflect.Func: // ignore the struct, map, chan, func
		default:
			_v = fmt.Sprintf("%v", v)
		}

		// json tag name naming format must be same as gorm table columns,
		// both should be "snake case" or "camel case".
		// gorm table columns naming format default to 'snake case',
		// so the json tag name is converted to "snake case here"
		// q[strcase.SnakeCase(jsonTag)] = val.Field(i).Interface()
		q[strcase.SnakeCase(jsonTag)] = _v
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
