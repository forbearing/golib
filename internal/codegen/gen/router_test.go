package gen_test

import (
	"go/ast"
	"testing"

	"github.com/forbearing/golib/internal/codegen/gen"
)

func log_print_helloworld() ast.Stmt {
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
		pkgName string
		stmts   []ast.Stmt
		want    string
		wantErr bool
	}{
		{
			name:    "log_println_hello_world",
			pkgName: "router",
			stmts:   []ast.Stmt{log_print_helloworld()},
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
			got, gotErr := gen.BuildRouterFile(tt.pkgName, tt.stmts...)
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
				t.Errorf("BuildRouterFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
