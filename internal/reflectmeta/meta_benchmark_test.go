package reflectmeta_test

import (
	"reflect"
	"testing"

	"github.com/forbearing/gst/internal/reflectmeta"
)

func BenchmarkReflectmeta_GetStructMeta(b *testing.B) {
	typ := reflect.TypeOf(User{})

	for b.Loop() {
		meta := reflectmeta.GetStructMeta(typ)
		for i := range meta.NumField() {
			_ = meta.JsonTag(i)
			_ = meta.SchemaTag(i)
			_ = meta.GormTag(i)
			_ = meta.QueryTag(i)
			_ = meta.UrlTag(i)
		}
	}
}

func BenchmarkNativeReflect(b *testing.B) {
	typ := reflect.TypeOf(User{})

	for b.Loop() {
		for i := range typ.NumField() {
			field := typ.Field(i)
			_ = field.Tag.Get("json")
			_ = field.Tag.Get("schema")
			_ = field.Tag.Get("gorm")
			_ = field.Tag.Get("query")
			_ = field.Tag.Get("url")
		}
	}
}
