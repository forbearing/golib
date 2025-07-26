package main

import (
	"bytes"
	"go/ast"
	"go/format"
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
// filename: model/user.go
package model

import pkgmodel "github.com/forbearing/golib/model"

type User struct {
	Name  string
	Age   int
	Email string

	pkgmodel.Base
}

type Group struct {
	Name    string
	Members []User

	pkgmodel.Base
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
			want:    "codegen",
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
			want: "pkgmodel",
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
		filename   string
		modulePath string
		want       []ModelInfo
		wantErr    bool
	}{
		{
			name:       "default",
			filename:   filename1,
			modulePath: modulePath,
			want: []ModelInfo{
				{PackageName: "model", ModelName: "User", ModelVarName: "u", ModulePath: "codegen", ModelFilePath: "/tmp/model"},
				{PackageName: "model", ModelName: "Group", ModelVarName: "g", ModulePath: "codegen", ModelFilePath: "/tmp/model"},
			},
			wantErr: false,
		},
		{
			name:       "named",
			filename:   filename2,
			modulePath: modulePath,
			want: []ModelInfo{
				{PackageName: "pkgmodel", ModelName: "User", ModelVarName: "u", ModulePath: "codegen", ModelFilePath: "/tmp/model"},
				{PackageName: "pkgmodel", ModelName: "Group", ModelVarName: "g", ModulePath: "codegen", ModelFilePath: "/tmp/model"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := findModels(tt.filename, tt.modulePath)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("findModels() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("findModels() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findModels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateServiceMethod(t *testing.T) {
	fset := token.NewFileSet()
	var buf bytes.Buffer

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		model      ModelInfo
		methodName string
		want       string
	}{
		// TODO: Add test cases.
		{
			name:       "user",
			methodName: "CreateBefore",
			model: ModelInfo{
				PackageName:   "model",
				ModelName:     "User",
				ModelVarName:  "u",
				ModulePath:    "codegen",
				ModelFilePath: "/tmp/model",
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
			res := generateServiceMethod(tt.model, tt.methodName)
			buf.Reset()
			if err := format.Node(&buf, fset, res); err != nil {
				t.Error(err)
				return
			}
			got := buf.String()
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateServiceMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}
