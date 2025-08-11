package gen

import (
	"go/ast"
	"testing"

	"github.com/forbearing/golib/types/consts"
	"github.com/kr/pretty"
)

func logPrintHelloworld() ast.Stmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("Println"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{Value: `"hello world"`},
			},
		},
	}
}

func TestBuildRouterFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		pkgName      string
		modelImports []string
		stmts        []ast.Stmt
		want         string
		wantErr      bool
	}{
		{
			name:         "log_println_hello_world",
			pkgName:      "router",
			modelImports: []string{"helloworld/model"},
			stmts:        []ast.Stmt{logPrintHelloworld()},
			want: `package router

import "log"

func Init() error {
	log.Println("hello world")
	return nil
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := BuildRouterFile(tt.pkgName, tt.modelImports, tt.stmts...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("BuildRouterFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("BuildRouterFile() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("BuildRouterFile() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestBuildServiceFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		pkgName     string
		modelImport []string
		types       []*ast.GenDecl
		stmts       []ast.Stmt
		phase       consts.Phase
		want        string
		wantErr     bool
	}{
		{
			name:        "user",
			pkgName:     "service",
			modelImport: []string{"helloworld/model"},
			types:       []*ast.GenDecl{types("model", "User", "User", "User", consts.PHASE_CREATE, false)},
			stmts:       []ast.Stmt{StmtServiceRegister("user")},
			want: `package service

import (
	"helloworld/model"

	"github.com/forbearing/golib/service"
)

func Init() error {
	service.Register[*user]()
	return nil
}

type userCreator struct {
	service.Base[*model.User, *model.User, *model.User]
}
`,
			wantErr: false,
		},
		{
			name:        "user_group",
			modelImport: []string{"helloworld/model"},
			pkgName:     "service",
			types: []*ast.GenDecl{
				types("model", "User", "User", "User", consts.PHASE_CREATE, false),
				types("model", "Group", "Group", "Group", consts.PHASE_UPDATE, false),
			},
			stmts: []ast.Stmt{StmtServiceRegister("user"), StmtServiceRegister("group")},
			want: `package service

import (
	"helloworld/model"

	"github.com/forbearing/golib/service"
)

func Init() error {
	service.Register[*user]()
	service.Register[*group]()
	return nil
}

type userCreator struct {
	service.Base[*model.User, *model.User, *model.User]
}

type groupUpdater struct {
	service.Base[*model.Group, *model.Group, *model.Group]
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := BuildServiceFile(tt.pkgName, tt.modelImport, tt.types, tt.stmts...)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("BuildServiceFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("BuildServiceFile() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("BuildServiceFile() = \n%v\n, want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}

func TestBuildMainFile(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		projectName string
		want        string
		wantErr     bool
	}{
		{
			name:        "test1",
			projectName: "helloworld",
			want: `package main

import (
	"helloworld/router"
	"helloworld/service"

	"github.com/forbearing/golib/bootstrap"
	. "github.com/forbearing/golib/util"
)

func main() {
	RunOrDie(bootstrap.Bootstrap)
	RunOrDie(service.Init)
	RunOrDie(router.Init)
	RunOrDie(bootstrap.Run)
}
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := BuildMainFile(tt.projectName)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("BuildMainFile() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("BuildMainFile() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("BuildMainFile() = \n%v,\n want \n%v\n", pretty.Sprintf("% #v", got), pretty.Sprintf("% #v", tt.want))
			}
		})
	}
}
