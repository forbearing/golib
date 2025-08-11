package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"go/token"
	"reflect"
	"testing"

	"github.com/forbearing/golib/types/consts"
	"github.com/kr/pretty"
)

func TestImports(t *testing.T) {
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
			got, err := FormatNode(Imports(tt.modulePath, tt.modelFileDir, tt.modelPkgName))
			if err != nil {
				t.Error(err)
				return
			}
			fmt.Println(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Imports() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestInits(t *testing.T) {
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
			got, err := FormatNode(Inits(tt.modelName))
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Inits() = \n%v\n, want \n%v\n", fmt.Sprintf("% #v", got), fmt.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestTypes(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modelPkgname string
		modelName    string
		reqName      string
		rspName      string
		phase        consts.Phase
		withComments bool
		want         string
	}{
		{
			name:         "user",
			modelPkgname: "model",
			modelName:    "User",
			reqName:      "User",
			rspName:      "User",
			phase:        consts.PHASE_CREATE,
			want: `type userCreator struct {
	service.Base[*model.User, *model.User, *model.User]
}`,
		},
		{
			name:         "user",
			modelPkgname: "model",
			modelName:    "User",
			reqName:      "UserReq",
			rspName:      "UserRsp",
			phase:        consts.PHASE_UPDATE,
			want: `type userUpdater struct {
	service.Base[*model.User, *model.UserReq, *model.UserRsp]
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := Types(tt.modelPkgname, tt.modelName, tt.reqName, tt.rspName, tt.phase, tt.withComments)
			var buf bytes.Buffer
			fset := token.NewFileSet()
			if err := format.Node(&buf, fset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("Types() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestServiceMethod1(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		recvName     string
		modelName    string
		modelPkgName string
		phase        consts.Phase
		want         string
	}{
		{
			name:         "CreateBefore",
			recvName:     "u",
			modelName:    "User",
			modelPkgName: "model",
			phase:        consts.PHASE_CREATE_BEFORE,
			want:         "func (u *userCreator) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}",
		},
		{
			name:         "UpdateAfter",
			recvName:     "g",
			modelName:    "Group",
			modelPkgName: "model_auth",
			phase:        consts.PHASE_UPDATE_AFTER,
			want:         "func (g *groupUpdater) UpdateAfter(ctx *types.ServiceContext, group *model_auth.Group) error {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatNode(ServiceMethod1(tt.recvName, tt.modelName, tt.modelPkgName, tt.phase))
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceMethod1() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestServiceMethod2(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		recvName     string
		modelName    string
		modelPkgName string
		phase        consts.Phase
		want         string
	}{
		{
			name:         "ListBefore",
			recvName:     "u",
			modelName:    "User",
			modelPkgName: "model",
			phase:        consts.PHASE_LIST_BEFORE,
			want:         "func (u *userLister) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {\n}",
		},
		{
			name:         "ListAfter",
			recvName:     "g",
			modelName:    "Group",
			modelPkgName: "model_auth",
			phase:        consts.PHASE_LIST_AFTER,
			want:         "func (g *groupLister) ListAfter(ctx *types.ServiceContext, groups *[]*model_auth.Group) error {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ServiceMethod2(tt.recvName, tt.modelName, tt.modelPkgName, tt.phase)
			got, err := FormatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceMethod2() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestServiceMethod3(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		recvName     string
		modelName    string
		phase        consts.Phase
		modelPkgName string
		want         string
	}{
		{
			name:         "CreateManyBefore",
			recvName:     "u",
			modelName:    "User",
			modelPkgName: "model",
			phase:        consts.PHASE_CREATE_MANY_BEFORE,
			want:         "func (u *userManyCreator) CreateManyBefore(ctx *types.ServiceContext, users ...*model.User) error {\n}",
		},
		{
			name:         "UpdateManyBefore",
			recvName:     "g",
			modelName:    "Group",
			modelPkgName: "model_auth",
			phase:        consts.PHASE_UPDATE_MANY_BEFORE,
			want:         "func (g *groupManyUpdater) UpdateManyBefore(ctx *types.ServiceContext, groups ...*model_auth.Group) error {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ServiceMethod3(tt.recvName, tt.modelName, tt.modelPkgName, tt.phase)
			got, err := FormatNode(res)
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceMethod3() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestServiceMethod4(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		recvName     string
		modelName    string
		modelPkgName string
		reqName      string
		rspName      string
		phase        consts.Phase
		want         string
	}{
		{
			name:         "Create",
			recvName:     "u",
			modelName:    "User",
			modelPkgName: "model",
			reqName:      "User",
			rspName:      "User",
			phase:        consts.PHASE_CREATE,
			want:         "func (u *userCreator) Create(ctx *types.ServiceContext, req *model.User) (rsp *model.User, err error) {\n}",
		},
		{
			name:         "Update",
			recvName:     "g",
			modelName:    "Group",
			modelPkgName: "model",
			reqName:      "GroupRequest",
			rspName:      "GroupResponse",
			phase:        consts.PHASE_UPDATE,
			want:         "func (g *groupUpdater) Update(ctx *types.ServiceContext, req *model.GroupRequest) (rsp *model.GroupResponse, err error) {\n}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ServiceMethod4(tt.recvName, tt.modelName, tt.modelPkgName, tt.reqName, tt.rspName, tt.phase)
			got, err := FormatNode(res)
			if err != nil {
				t.Error(err)
				return
			}

			if got != tt.want {
				t.Errorf("ServiceMethod4() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}
