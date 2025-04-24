package reflectmeta

import (
	"reflect"
	"sync"
)

var metaCache sync.Map // map[string]*StructMeta

type StructMeta struct {
	Type     reflect.Type
	Fields   []reflect.StructField
	numField int
	FieldMap map[string]int // 字段名 -> 索引

	JSONNames   []string
	SchemaNames []string
	GormNames   []string
	QueryNames  []string
	UrlNames    []string
}

func GetStructMeta(t reflect.Type) *StructMeta {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	key := t.PkgPath() + "." + t.Name()
	if meta, ok := metaCache.Load(key); ok {
		return meta.(*StructMeta)
	}

	fieldCount := t.NumField()
	fields := make([]reflect.StructField, 0, fieldCount)
	fieldMap := make(map[string]int)
	jsonNames := make([]string, 0, fieldCount)
	schemaNames := make([]string, 0, fieldCount)
	gormNames := make([]string, 0, fieldCount)
	queryNames := make([]string, 0, fieldCount)
	urlNames := make([]string, 0, fieldCount)

	var parseFields func(reflect.Type)
	parseFields = func(rt reflect.Type) {
		for i := range rt.NumField() {
			field := rt.Field(i)

			// 如果是匿名字段且是结构体，递归展开
			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				parseFields(field.Type)
				continue
			}

			fields = append(fields, field)

			// // NOTE: strings.Split always returns at least one element(empty string)
			// // We should not use len(jsonTagItems) to check the json tags exists.
			// jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
			// schemaTag := strings.Split(field.Tag.Get("schema"), ",")[0]
			// jsonNames = append(jsonNames, jsonTag)
			// schemaNames = append(schemaNames, schemaTag)

			jsonNames = append(jsonNames, field.Tag.Get("json"))
			schemaNames = append(schemaNames, field.Tag.Get("schema"))
			gormNames = append(gormNames, field.Tag.Get("gorm"))
			queryNames = append(queryNames, field.Tag.Get("query"))
			urlNames = append(urlNames, field.Tag.Get("url"))

			fieldMap[field.Name] = len(fields) - 1
		}
	}
	parseFields(t)

	meta := &StructMeta{
		Type:     t,
		Fields:   fields,
		numField: fieldCount,
		FieldMap: fieldMap,

		JSONNames:   jsonNames,
		SchemaNames: schemaNames,
		GormNames:   gormNames,
		QueryNames:  queryNames,
		UrlNames:    urlNames,
	}
	metaCache.Store(key, meta)
	return meta
}

func (m *StructMeta) String() string                  { return m.Type.String() }
func (m *StructMeta) Field(i int) reflect.StructField { return m.Fields[i] }
func (m *StructMeta) NumField() int                   { return m.numField }
func (m *StructMeta) JsonTag(i int) string            { return m.JSONNames[i] }
func (m *StructMeta) SchemaTag(i int) string          { return m.SchemaNames[i] }
func (m *StructMeta) GormTag(i int) string            { return m.GormNames[i] }
func (m *StructMeta) QueryTag(i int) string           { return m.QueryNames[i] }
func (m *StructMeta) UrlTag(i int) string             { return m.UrlNames[i] }
