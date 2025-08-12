package gen

import (
	"fmt"
	"go/ast"
	"go/token"
)

// BuildModelFile generates a model.go file, the content like below:
/*
package model

import "github.com/forbearing/golib/model"

func Init() error {
	model.Register[*Group]()
	model.Register[*User]()

	return nil
}
*/
func BuildModelFile(pkgName string, modelImports []string, stmts ...ast.Stmt) (string, error) {
	// Create init function body
	body := make([]ast.Stmt, 0)
	body = append(body, stmts...)
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})

	// Create Init function declaration
	initDecl := &ast.FuncDecl{
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
	}

	// Create import declaration
	importDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/forbearing/golib/model"`,
				},
			},
		},
	}
	for _, modelImport := range modelImports {
		importDecl.Specs = append(importDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf(`"%s"`, modelImport),
			},
		})
	}

	// Create file AST
	f := &ast.File{
		Name:  ast.NewIdent(pkgName),
		Decls: []ast.Decl{
			// NOTE: imports must appear before other declarations
		},
	}

	// Add imports
	f.Decls = append(f.Decls, importDecl)
	// Add init function
	f.Decls = append(f.Decls, initDecl)

	return formatAndImports(f, false)
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
func BuildServiceFile(pkgName string, modelImports []string, types []*ast.GenDecl, stmts ...ast.Stmt) (string, error) {
	body := make([]ast.Stmt, 0)
	body = append(body, stmts...)
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})

	initDecl := &ast.FuncDecl{
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
	}

	// imports service
	imports := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/forbearing/golib/service"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/forbearing/golib/types/consts"`,
				},
			},
		},
	}
	// imports, such like: "helloworld/model"
	for _, name := range modelImports {
		imports.Specs = append(imports.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", name),
			},
		})
	}

	f := &ast.File{
		Name:  ast.NewIdent(pkgName),
		Decls: []ast.Decl{
			// NOTE: imports must appear before other declarations
		},
	}
	// imports
	f.Decls = append(f.Decls, imports)
	// Init() declarations.
	f.Decls = append(f.Decls, initDecl)
	// type declarations.
	for _, typ := range types {
		f.Decls = append(f.Decls, typ)
	}

	return formatAndImports(f, false)
}

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
func BuildRouterFile(pkgName string, modelImports []string, stmts ...ast.Stmt) (string, error) {
	body := make([]ast.Stmt, 0)
	body = append(body, stmts...)
	body = append(body, &ast.ReturnStmt{
		Results: []ast.Expr{
			ast.NewIdent("nil"),
		},
	})

	initDecl := &ast.FuncDecl{
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
	}

	importDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/forbearing/golib/router"`,
				},
			},
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: `"github.com/forbearing/golib/types/consts"`,
				},
			},
		},
	}
	for _, imp := range modelImports {
		importDecl.Specs = append(importDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", imp),
			},
		})
	}

	f := &ast.File{
		Name:  ast.NewIdent(pkgName),
		Decls: []ast.Decl{
			// NOTE: imports must appear before other declarations
		},
	}

	f.Decls = append(f.Decls, importDecl)
	f.Decls = append(f.Decls, initDecl)

	return formatAndImports(f, false)
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

	return formatAndImports(f, false)
}
