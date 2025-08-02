package codegen

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_imports(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modulePath   string
		modelFileDir string
		modelPkgName string

		want string
	}{
		{
			modulePath:   "codegen",
			modelFileDir: "model",
			modelPkgName: "model",
			want: `import (
	"codegen/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatNode(imports(tt.modulePath, tt.modelFileDir, tt.modelPkgName))
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("imports() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_inits(t *testing.T) {
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
			got, err := formatNode(inits(tt.modelName))
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_types(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modelName    string
		modelPkgname string
		want         string
	}{
		// {
		// 	name:      "user",
		// 	modelName: "User",
		// 	want: `// user implements the types.Service[*model.User] interface.
		// type user struct {
		// 	service.Base[*model.User]
		// }`,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := types(tt.modelName, tt.modelPkgname)
			got, err := formatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				// t.Errorf("types() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_method_1(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		methodName   string
		recvName     string
		modelName    string
		modelPkgName string
		want         string
	}{
		{
			name:         "CreateBefore",
			recvName:     "u",
			modelName:    "User",
			methodName:   "CreateBefore",
			modelPkgName: "model",
			want:         "func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}",
		},
		{
			name:         "UpdateAfter",
			recvName:     "g",
			modelName:    "Group",
			modelPkgName: "model_auth",
			methodName:   "UpdateAfter",
			want:         "func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model_auth.Group) error {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatNode(service_method_1(tt.recvName, tt.modelName, tt.methodName, tt.modelPkgName))
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("service_method_1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_method_2(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		recvName     string
		modelName    string
		methodName   string
		modelPkgName string
		want         string
	}{
		{
			name:         "ListBefore",
			recvName:     "u",
			modelName:    "User",
			methodName:   "ListBefore",
			modelPkgName: "model",
			want:         "func (u *user) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {\n}",
		},
		{
			name:         "ListAfter",
			recvName:     "g",
			modelName:    "Group",
			methodName:   "ListAfter",
			modelPkgName: "model_auth",
			want:         "func (g *group) ListAfter(ctx *types.ServiceContext, groups *[]*model_auth.Group) error {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := service_method_2(tt.recvName, tt.modelName, tt.methodName, tt.modelPkgName)
			got, err := formatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("service_method_2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_method_3(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		recvName     string
		modelName    string
		methodName   string
		modelPkgName string
		want         string
	}{
		{
			name:         "CreateManyBefore",
			recvName:     "u",
			modelName:    "User",
			methodName:   "CreateManyBefore",
			modelPkgName: "model",
			want:         "func (u *user) CreateManyBefore(ctx *types.ServiceContext, users ...*model.User) error {\n}",
		},
		{
			name:         "UpdateManyBefore",
			recvName:     "g",
			modelName:    "Group",
			methodName:   "UpdateManyBefore",
			modelPkgName: "model_auth",
			want:         "func (g *group) UpdateManyBefore(ctx *types.ServiceContext, groups ...*model_auth.Group) error {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := service_method_3(tt.recvName, tt.modelName, tt.methodName, tt.modelPkgName)
			got, err := formatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("service_method_3() = %v, want %v", got, tt.want)
			}
		})
	}
}
