package main

import (
	"bytes"
	"go/format"
	"go/token"
	"reflect"
	"testing"
)

func Test_service_method_1(t *testing.T) {
	fset := token.NewFileSet()
	var buf bytes.Buffer

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		methodName string
		recvName   string
		modelName  string
		want       string
	}{
		{
			name:       "CreateBefore",
			recvName:   "u",
			modelName:  "User",
			methodName: "CreateBefore",
			want:       "func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}",
		},
		{
			name:       "UpdateAfter",
			recvName:   "g",
			modelName:  "Group",
			methodName: "UpdateAfter",
			want:       "func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := service_method_1(tt.recvName, tt.modelName, tt.methodName)
			buf.Reset()
			if err := format.Node(&buf, fset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("service_method_1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_imports(t *testing.T) {
	filset := token.NewFileSet()
	var buf bytes.Buffer

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modulePath    string
		modelFilePath string
		want          string
	}{
		{
			modulePath:    "codegen",
			modelFilePath: "model",
			want: `import (
	"codegen/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := imports(tt.modulePath, tt.modelFilePath)
			buf.Reset()
			if err := format.Node(&buf, filset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inits(t *testing.T) {
	fileset := token.NewFileSet()
	var buf bytes.Buffer

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modelName string
		want      string
	}{
		{
			name:      "user",
			modelName: "User",
			want: `func init() {
	service.Register[*user]()
}`,
		},
		{
			name:      "group",
			modelName: "Group",
			want: `func init() {
	service.Register[*group]()
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := inits(tt.modelName)
			buf.Reset()
			if err := format.Node(&buf, fileset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inits() = %v, want %v", got, tt.want)
			}
		})
	}
}
