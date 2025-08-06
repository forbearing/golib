package gen

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/imports"
)

// BuildRouterFile generates a router.go file, the content like below:
/*
package router

import (
	"helloworld/model"

	"github.com/forbearing/golib/router"
)

func Init() error {
	router.Register[*model.Group, *model.Group, *model.Group](router.API(), "group")
	router.Register[*model.User, *model.User, *model.User](router.API(), "user")
	return nil
}
*/
// FIXME: process imports automatically problem.
func BuildRouterFile(pkgName string, stmts ...ast.Stmt) (string, error) {
	body := make([]ast.Stmt, 0)
	body = append(body, stmts...)
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})

	f := &ast.File{
		Name: ast.NewIdent(pkgName),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/forbearing/golib/router"`,
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("Init"),
				Type: &ast.FuncType{
					TypeParams: nil,
					Params:     nil,
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: ast.NewIdent("error")},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: body,
				},
			},
		},
	}

	return formatAndImports(f)
}

// BuildServiceFile generates a service.go file, the content like below:
/*
package service

import "github.com/forbearing/golib/service"

func Init() error {
	service.Register[*group]()
	service.Register[*user]()
	return nil
}
*/
// FIXME: process imports automatically problem.
func BuildServiceFile(pkgName string, stmts ...ast.Stmt) (string, error) {
	body := make([]ast.Stmt, 0)
	body = append(body, stmts...)
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})

	f := &ast.File{
		Name: ast.NewIdent(pkgName),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/forbearing/golib/service"`,
						},
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("Init"),
				Type: &ast.FuncType{
					TypeParams: nil,
					Params:     nil,
					Results: &ast.FieldList{
						List: []*ast.Field{
							{Type: ast.NewIdent("error")},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: body,
				},
			},
		},
	}

	return formatAndImports(f)
}

// BuildMainFile generates a main.go file, the content like below:
/*
package main

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
*/
func BuildMainFile(projectName string) (string, error) {
	f := &ast.File{
		Name: ast.NewIdent("main"),
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", projectName+"/service")}},
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", projectName+"/router")}},
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", "github.com/forbearing/golib/bootstrap")}},
					&ast.ImportSpec{
						Path: &ast.BasicLit{Value: fmt.Sprintf("%q", "github.com/forbearing/golib/util")},
						Name: ast.NewIdent("."),
					},
				},
			},
			&ast.FuncDecl{
				Name: ast.NewIdent("main"),
				Type: &ast.FuncType{},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: ast.NewIdent("RunOrDie"),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("bootstrap"),
										Sel: ast.NewIdent("Bootstrap"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: ast.NewIdent("RunOrDie"),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("service"),
										Sel: ast.NewIdent("Init"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: ast.NewIdent("RunOrDie"),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("router"),
										Sel: ast.NewIdent("Init"),
									},
								},
							},
						},
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: ast.NewIdent("RunOrDie"),
								Args: []ast.Expr{
									&ast.SelectorExpr{
										X:   ast.NewIdent("bootstrap"),
										Sel: ast.NewIdent("Run"),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return formatAndImports(f)
}

// formatAndImports formats code use gofumpt and processes imports.
func formatAndImports(f *ast.File) (string, error) {
	formatted, err := FormatNodeExtra(f)
	if err != nil {
		return "", err
	}

	result, err := imports.Process("", []byte(formatted), nil)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
