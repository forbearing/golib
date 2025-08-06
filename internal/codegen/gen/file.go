package gen

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/imports"
)

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
