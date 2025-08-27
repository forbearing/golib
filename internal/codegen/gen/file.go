package gen

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/forbearing/golib/types/consts"
)

// BuildModelFile generates a model.go file, the content like below:
/*
package model

import "github.com/forbearing/golib/model"

func init() {
	model.Register[*Group]()
	model.Register[*User]()
}
*/
func BuildModelFile(pkgName string, modelImports []string, stmts ...ast.Stmt) (string, error) {
	// Create init function body
	body := make([]ast.Stmt, 0)
	body = append(body, stmts...)

	// Create Init function declaration
	initDecl := &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{
			TypeParams: nil,
			Params:     nil,
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

	// Create generated code comment at the top of the file
	generatedComment := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text:  consts.CodeGeneratedComment(),
				Slash: token.Pos(1),
			},
		},
	}
	f.Comments = []*ast.CommentGroup{generatedComment}
	// Set package name position to ensure comment appears before it
	f.Name.NamePos = token.Pos(2)

	// If the caller does not pass stmts or stmts is empty, then the Init function body is empty,
	// So we should not imports any external package.
	if len(stmts) != 0 {
		// Add imports
		f.Decls = append(f.Decls, importDecl)
	}
	// Add init function
	f.Decls = append(f.Decls, initDecl)

	return FormatNodeExtra(f, false)
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
func BuildServiceFile(pkgName string, modelImports []string, stmts ...ast.Stmt) (string, error) {
	// Handle import conflicts when modelImports contain packages with same base name
	// For example: ["nebula/service/pkg1/user", "nebula/service/pkg2/user"]
	// Should be renamed to:
	// import (
	//     pkg1_user "nebula/service/pkg1/user"
	//     pkg2_user "nebula/service/pkg2/user"
	// )
	importAliases := ResolveImportConflicts(modelImports)

	body := make([]ast.Stmt, 0)
	body = append(body, stmts...)

	initDecl := &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{
			TypeParams: nil,
			Params:     nil,
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
	// Use aliases to resolve import conflicts
	for _, importPath := range modelImports {
		alias := importAliases[importPath]
		importSpec := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", importPath),
			},
		}
		// Add alias if needed to resolve conflicts
		if alias != "" {
			importSpec.Name = ast.NewIdent(alias)
		}
		imports.Specs = append(imports.Specs, importSpec)
	}

	f := &ast.File{
		Name:  ast.NewIdent(pkgName),
		Decls: []ast.Decl{
			// NOTE: imports must appear before other declarations
		},
	}

	// Create generated code comment at the top of the file
	generatedComment := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text:  consts.CodeGeneratedComment(),
				Slash: token.Pos(1),
			},
		},
	}
	f.Comments = []*ast.CommentGroup{generatedComment}
	// Set package name position to ensure comment appears before it
	f.Name.NamePos = token.Pos(2)

	// If the caller does not pass stmts or stmts is empty, then the Init function body is empty,
	// So we should not imports any external package.
	if len(stmts) != 0 {
		// imports
		f.Decls = append(f.Decls, imports)
	}
	// Init() declarations.
	f.Decls = append(f.Decls, initDecl)

	return FormatNodeExtra(f, false)
}

// BuildRouterFile generates a router.go file, the content like below:
/*
package router

import (
	"helloworld/model"

	"github.com/forbearing/golib/router"
)

func Init() error {
	router.Register[*model.Group, *model.Group, *model.Group](router.Auth(), "group")
	router.Register[*model.User, *model.User, *model.User](router.Pub(), "user")
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

	// Create generated code comment at the top of the file
	generatedComment := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text:  consts.CodeGeneratedComment(),
				Slash: token.Pos(1),
			},
		},
	}
	f.Comments = []*ast.CommentGroup{generatedComment}
	// Set package name position to ensure comment appears before it
	f.Name.NamePos = token.Pos(2)

	// If the caller does not pass stmts or stmts is empty, then the Init function body is empty,
	// So we should not imports any external package.
	if len(stmts) != 0 {
		// imports
		f.Decls = append(f.Decls, importDecl)
	}
	// Init() declarations.
	f.Decls = append(f.Decls, initDecl)

	return FormatNodeExtra(f, false)
}

// BuildMainFile generates a main.go file, the content like below:
/*
package main

import (
	"helloworld/configx"
	"helloworld/cronjob"
	_ "helloworld/model"
	"helloworld/router"
	"helloworld/service"

	"github.com/forbearing/golib/bootstrap"
	. "github.com/forbearing/golib/util"
)

func main() {
	RunOrDie(bootstrap.Bootstrap)
	RunOrDie(configx.Init)
	RunOrDie(cronjob.Init)
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
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", projectName+"/configx")}},
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", projectName+"/cronjob")}},
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", projectName+"/middleware")}, Name: ast.NewIdent("_")},
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", projectName+"/model")}, Name: ast.NewIdent("_")},
					&ast.ImportSpec{Path: &ast.BasicLit{Value: fmt.Sprintf("%q", projectName+"/service")}, Name: ast.NewIdent("_")},
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
										X:   ast.NewIdent("configx"),
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
										X:   ast.NewIdent("cronjob"),
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

	// Create generated code comment at the top of the file
	generatedComment := &ast.CommentGroup{
		List: []*ast.Comment{
			{
				Text:  consts.CodeGeneratedComment(),
				Slash: token.Pos(1),
			},
		},
	}
	f.Comments = []*ast.CommentGroup{generatedComment}
	// Set package name position to ensure comment appears before it
	f.Name.NamePos = token.Pos(2)

	return FormatNodeExtra(f, false)
}
