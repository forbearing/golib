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

	FieldIndexes [][]int // 每个字段的 Index 路径（支持嵌套匿名字段）
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

	fieldIndexes := make([][]int, 0, fieldCount)

	var parseFields func(reflect.Type, []int)
	parseFields = func(rt reflect.Type, parentIndex []int) {
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			indexPath := append(parentIndex, i)

			if field.Anonymous && field.Type.Kind() == reflect.Struct {
				parseFields(field.Type, indexPath)
				continue
			}

			fields = append(fields, field)
			fieldIndexes = append(fieldIndexes, indexPath)

			jsonNames = append(jsonNames, field.Tag.Get("json"))
			schemaNames = append(schemaNames, field.Tag.Get("schema"))
			gormNames = append(gormNames, field.Tag.Get("gorm"))
			queryNames = append(queryNames, field.Tag.Get("query"))
			urlNames = append(urlNames, field.Tag.Get("url"))
			fieldMap[field.Name] = len(fields) - 1
		}
	}

	parseFields(t, []int{})

	meta := &StructMeta{
		Type:         t,
		Fields:       fields,
		numField:     fieldCount,
		FieldMap:     fieldMap,
		FieldIndexes: fieldIndexes,

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
