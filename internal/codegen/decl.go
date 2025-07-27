package codegen

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/gertd/go-pluralize"
)

var pluralizeCli = pluralize.NewClient()

// imports
/*
import (
	"codegen/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)
*/
func imports(modulePath, modelFileDir, modelPkgName string) *ast.GenDecl {
	importModel := filepath.Join(modulePath, modelFileDir)
	fields := strings.Split(importModel, "/")
	if len(fields) > 0 && fields[len(fields)-1] != modelPkgName {
		// model_setting "mymodule/model/setting"
		importModel = fmt.Sprintf("%s %q", modelPkgName, importModel)
	} else {
		// "mymodule/model"
		importModel = fmt.Sprintf("%q", importModel)
	}

	return &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: importModel,
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
func types(modelName, modelPkgName string) *ast.GenDecl {
	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{
					Text: fmt.Sprintf("// %s implements the types.Service[*%s.%s] interface.",
						strings.ToLower(modelName), modelPkgName, modelName),
				},
			},
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(strings.ToLower(modelName)),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.IndexExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("service"),
										Sel: ast.NewIdent("Base"),
									},
									Index: &ast.StarExpr{
										X: &ast.SelectorExpr{
											X:   ast.NewIdent(modelPkgName),
											Sel: ast.NewIdent(modelName),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// service_method_1 creates a function declaration for a method in a service interface.
// For example:
//
//	"func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}"
//	"func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {\n}",
func service_method_1(recvName, modelName, methodName, modelPkgName string, body ...ast.Stmt) *ast.FuncDecl {
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
								X:   ast.NewIdent(modelPkgName),
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

// service_method_2
// For example:
//
//	"func (u *user) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
//	"func (u *user) ListAfter(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
func service_method_2(recvName, modelName, methodName, modelPkgName string, body ...ast.Stmt) *ast.FuncDecl {
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
						Names: []*ast.Ident{ast.NewIdent(pluralizeCli.Plural(strings.ToLower(modelName)))},
						Type: &ast.StarExpr{
							X: &ast.ArrayType{
								Elt: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent(modelPkgName),
										Sel: ast.NewIdent(modelName),
									},
								},
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

// service_method_3
// For example:
//
//	"func (u *user) BatchCreateBefore(ctx *types.ServiceContext, users ...*model.User) error {\n}"
//	"func (u *user) BatchCreateAfter(ctx *types.ServiceContext, users ...*model.User) error {\n}"
func service_method_3(recvName, modelName, methodName, modelPkgName string, body ...ast.Stmt) *ast.FuncDecl {
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
						Names: []*ast.Ident{ast.NewIdent(pluralizeCli.Plural(strings.ToLower(modelName)))},
						Type: &ast.Ellipsis{
							Elt: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent(modelPkgName),
									Sel: ast.NewIdent(modelName),
								},
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
