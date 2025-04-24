package reflectmeta_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/forbearing/golib/internal/reflectmeta"
)

type Base struct {
	ID        string     `json:"id" gorm:"primaryKey" schema:"id" url:"-"`
	CreatedBy string     `json:"created_by,omitempty" gorm:"index" schema:"created_by" url:"-"`
	UpdatedBy string     `json:"updated_by,omitempty" gorm:"index" schema:"updated_by" url:"-"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"index" schema:"-" url:"-"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" gorm:"index" schema:"-" url:"-"`
	Remark    *string    `json:"remark,omitempty" gorm:"size:10240" schema:"-" url:"-"`
	Order     *uint      `json:"order,omitempty" schema:"-" url:"-"`
}

type User struct {
	Base
	Name  string `json:"name"`
	Email string `schema:"email_address"`
	Age   int
}

func TestGetStructMeta(t *testing.T) {
	meta := reflectmeta.GetStructMeta(reflect.TypeOf(User{}))

	if meta.Type.Name() != "User" {
		t.Errorf("expected type name 'User', got: %s", meta.Type.Name())
	}

	expectedFieldCount := 10 // 7 Base + 3 User
	if len(meta.Fields) != expectedFieldCount {
		t.Errorf("expected %d fields, got: %d", expectedFieldCount, len(meta.Fields))
	}

	expectJSON := []string{
		"id", "created_by,omitempty", "updated_by,omitempty", "created_at,omitempty", "updated_at,omitempty", "remark,omitempty", "order,omitempty",
		"name", "", // Email field no json tag
	}
	expectSchema := []string{
		"id", "created_by", "updated_by", "-", "-", "-", "-",
		"", "email_address", // Name field no schema tag
	}

	// fmt.Println(len(meta.JSONNames), len(expectJSON))
	// fmt.Println(meta.JSONNames)
	// fmt.Println(expectJSON)
	for i := range expectJSON {
		if meta.JSONNames[i] != expectJSON[i] {
			t.Errorf("json tag mismatch at index %d: expected %q, got %q", i, expectJSON[i], meta.JSONNames[i])
		}
		if meta.SchemaNames[i] != expectSchema[i] {
			t.Errorf("schema tag mismatch at index %d: expected %q, got %q", i, expectSchema[i], meta.SchemaNames[i])
		}
	}

	if idx, ok := meta.FieldMap["Email"]; !ok || idx != 8 {
		t.Errorf("expected Email at index 8 in FieldMap, got index: %d, exists: %v", idx, ok)
	}
}
