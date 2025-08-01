package codegen

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var src1 = `
package model

import "github.com/forbearing/golib/model"

type User struct {
	Name  string
	Age   int
	Email string

	model.Base
}

type Group struct {
	Name    string
	Members []User

	model.Base
}

type GroupUser struct {
	GroupId int
	UserId  int
}
	`

var src2 = `
package model

import model_auth "github.com/forbearing/golib/model"

type User struct {
	Name  string
	Age   int
	Email string

	model_auth.Base
}

type Group struct {
	Name    string
	Members []User

	model_auth.Base
}

type GroupUser struct {
	GroupId int
	UserId  int
}
	`

func Test_getModulePath(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			want:    "github.com/forbearing/golib",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := getModulePath()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getModulePath() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getModulePath() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("getModulePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findModelPackageName(t *testing.T) {
	fset := token.NewFileSet()
	file1, err := parser.ParseFile(fset, "user.go", src1, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	file2, err := parser.ParseFile(fset, "user.go", src2, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		file *ast.File
		want string
	}{
		{
			name: "default",
			file: file1,
			want: "model",
		},
		{
			name: "named",
			file: file2,
			want: "model",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findModelPackageName(tt.file)
			if got != tt.want {
				t.Errorf("findModelPackageName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findModels(t *testing.T) {
	modulePath, err := getModulePath()
	if err != nil {
		t.Fatal(err)
	}

	tmpdir := "/tmp/model"
	if err = os.MkdirAll(tmpdir, 0o750); err != nil {
		t.Fatal(err)
	}

	filename1 := filepath.Join(tmpdir, "user.go")
	filename2 := filepath.Join(tmpdir, "user2.go")
	if err = os.WriteFile(filename1, []byte(src1), 0o644); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filename2, []byte(src2), 0o644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		modulePath string
		filename   string
		want       []*ModelInfo
		wantErr    bool
	}{
		{
			name:       "default",
			modulePath: modulePath,
			filename:   filename1,
			want: []*ModelInfo{
				{ModulePath: "github.com/forbearing/golib", PackageName: "model", ModelName: "User", ModelVarName: "u", ModelFileDir: "/tmp/model"},
				{ModulePath: "github.com/forbearing/golib", PackageName: "model", ModelName: "Group", ModelVarName: "g", ModelFileDir: "/tmp/model"},
			},
			wantErr: false,
		},
		{
			name:       "named",
			modulePath: modulePath,
			filename:   filename2,
			want: []*ModelInfo{
				{ModulePath: "github.com/forbearing/golib", PackageName: "model", ModelName: "User", ModelVarName: "u", ModelFileDir: "/tmp/model"},
				{ModulePath: "github.com/forbearing/golib", PackageName: "model", ModelName: "Group", ModelVarName: "g", ModelFileDir: "/tmp/model"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := findModels(tt.modulePath, tt.filename)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("findModels() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("findModels() succeeded unexpectedly")
			}
			var got2 []ModelInfo
			var want2 []ModelInfo
			for _, v := range got {
				got2 = append(got2, *v)
			}
			for _, v := range tt.want {
				want2 = append(want2, *v)
			}
			if !reflect.DeepEqual(got2, want2) {
				t.Errorf("findModels() = %v, want %v", got2, want2)
			}
		})
	}
}

func Test_modelPkg2ServicePkg(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		pkgName string
		want    string
	}{
		{
			name:    "test1",
			pkgName: "model",
			want:    "service",
		},
		{
			name:    "test2",
			pkgName: "model2",
			want:    "service2",
		},
		{
			name:    "test3",
			pkgName: "model_system",
			want:    "service_system",
		},
		{
			name:    "test4",
			pkgName: "modelAuth",
			want:    "serviceAuth",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modelPkg2ServicePkg(tt.pkgName)
			if got != tt.want {
				t.Errorf("modelPkg2ServicePkg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateServiceMethod(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		model      *ModelInfo
		methodName string
		want       string
	}{
		// TODO: Add test cases.
		{
			name:       "user",
			methodName: "CreateBefore",
			model: &ModelInfo{
				PackageName:  "model",
				ModelName:    "User",
				ModelVarName: "u",
				ModulePath:   "codegen",
				ModelFileDir: "/tmp/model",
			},
			want: `func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {
        log := u.WithServiceContext(ctx, ctx.GetPhase())
        log.Info("user create before")
        return nil
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatNode(generateServiceMethod1(tt.model, tt.methodName))
			if err != nil {
				t.Error(err)
				return
			}
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateServiceMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateServiceFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		info *ModelInfo
		want string
	}{
		// 		{
		// 			info: ModelInfo{PackageName: "model", ModelName: "User", ModelVarName: "u", ModulePath: "codegen", ModelFilePath: "model"},
		// 			want: `package service
		//
		// import (
		// 	"codegen/model"
		// 	"github.com/forbearing/golib/service"
		// 	"github.com/forbearing/golib/types"
		// )
		//
		// func init() {
		// 	service.Register[*user]()
		// }
		//
		// // user implements the types.Service[*model.User] interface.
		// type user struct {
		// 	service.Base[*model.User]
		// }
		//
		// func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {
		// 	log := u.WithServiceContext(ctx, ctx.GetPhase())
		// 	log.Info("user create before")
		// 	return nil
		// }
		// `,
		// 		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatNode(generateServiceFile(tt.info))
			if err != nil {
				t.Error(err)
				return
			}
			if got != tt.want {
				t.Errorf("generateServiceFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
