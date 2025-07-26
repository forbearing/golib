package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"
)

// service_method_1 creates a function declaration for a method in a service interface.
// For example:
//
//	func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}"
//	"func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {\n}",
//
// More details see test case.
func service_method_1(recvName, modelName, methodName string, body ...ast.Stmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(recvName)},
					Type: &ast.StarExpr{
						X: ast.NewIdent(strings.ToLower(modelName)),
					},
				},
			},
		},
		Name: ast.NewIdent(methodName),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("types"),
								Sel: ast.NewIdent("ServiceContext"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent(strings.ToLower(modelName))},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("model"),
								Sel: ast.NewIdent(modelName),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

// imports
/*
import (
	"codegen/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)
*/
func imports(modulePath, modelFilePath string) *ast.GenDecl {
	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", filepath.Join(modulePath, modelFilePath)),
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/forbearing/golib/service"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/forbearing/golib/types"`,
				},
			},
		},
	}
}

// inits
/*
func init() {
	service.Register[*user]()
}
*/
func inits(modelName string) *ast.FuncDecl {
	return &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.IndexExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("service"),
								Sel: ast.NewIdent("Register"),
							},
							Index: &ast.StarExpr{
								X: ast.NewIdent(strings.ToLower(modelName)),
							},
						},
					},
				},
			},
		},
	}
}

// types
/*
// user implements the types.Service[*model.User] interface.
type user struct {
	service.Base[*model.User]
}
*/
func types(modelName string) *ast.GenDecl {
	return &ast.GenDecl{
		Doc: nil,
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Type: &ast.StructType{},
			},
		},
	}
}
