package gen

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/gertd/go-pluralize"
)

var pluralizeCli = pluralize.NewClient()

// Imports generates an ast node that represents the declaration of below:
/*
import (
	"codegen/model"
	"github.com/forbearing/golib/service"
	"github.com/forbearing/golib/types"
)
*/
func Imports(modulePath, modelFileDir, modelPkgName string) *ast.GenDecl {
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

// Inits returns an ast node that represents the declaration of below:
/*
func init() {
	service.Register[*user]()
}
*/
func Inits(modelNames ...string) *ast.FuncDecl {
	list := make([]ast.Stmt, 0)

	for _, name := range modelNames {
		list = append(list,
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.IndexExpr{
						X: &ast.SelectorExpr{
							X:   ast.NewIdent("service"),
							Sel: ast.NewIdent("Register"),
						},
						Index: &ast.StarExpr{
							X: ast.NewIdent(strings.ToLower(name)),
						},
					},
				},
			},
		)
	}

	return &ast.FuncDecl{
		Name: ast.NewIdent("init"),
		Type: &ast.FuncType{},
		Body: &ast.BlockStmt{
			List: list,
		},
	}
}

// Types returns an ast node that represents the declaration of below:
/*
// user implements the types.Service[*model.User, *model.User, *model.User] interface.
type user struct {
	service.Base[*model.User, *model.User, *model.User]
}
*/
func Types(modelPkgName, modelName, reqName, rspName string, withComment bool) *ast.GenDecl {
	comments := []*ast.Comment{}

	if withComment {
		comments = append(comments, &ast.Comment{
			Text: fmt.Sprintf("// %s implements the types.Service[*%s.%s, *%s.%s, *%s.%s] interface.",
				strings.ToLower(modelName), modelPkgName, modelName, modelPkgName, modelName, modelPkgName, modelName),
		})
	}

	return &ast.GenDecl{
		Doc: &ast.CommentGroup{
			List: comments,
		},
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: ast.NewIdent(strings.ToLower(modelName)),
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							// {
							// 	Type: &ast.IndexExpr{
							// 		X: &ast.SelectorExpr{
							// 			X:   ast.NewIdent("service"),
							// 			Sel: ast.NewIdent("Base"),
							// 		},
							// 		Index: &ast.StarExpr{
							// 			X: &ast.SelectorExpr{
							// 				X:   ast.NewIdent(modelPkgName),
							// 				Sel: ast.NewIdent(modelName),
							// 			},
							// 		},
							// 	},
							// },
							{
								Type: &ast.IndexListExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("service"),
										Sel: ast.NewIdent("Base"),
									},
									Indices: []ast.Expr{
										&ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent(modelPkgName),
												Sel: ast.NewIdent(modelName),
											},
										},
										&ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent(modelPkgName),
												Sel: ast.NewIdent(reqName),
											},
										},
										&ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent(modelPkgName),
												Sel: ast.NewIdent(rspName),
											},
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

// ServiceMethod1 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *user) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}"
//	"func (g *group) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {\n}",
func ServiceMethod1(recvName, modelName, methodName, modelPkgName string, body ...ast.Stmt) *ast.FuncDecl {
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

// ServiceMethod2 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *user) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
//	"func (u *user) ListAfter(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
func ServiceMethod2(recvName, modelName, methodName, modelPkgName string, body ...ast.Stmt) *ast.FuncDecl {
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

// ServiceMethod3 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *user) CreateManyBefore(ctx *types.ServiceContext, users ...*model.User) error {\n}"
//	"func (u *user) CreateManyAfter(ctx *types.ServiceContext, users ...*model.User) error {\n}"
func ServiceMethod3(recvName, modelName, methodName, modelPkgName string, body ...ast.Stmt) *ast.FuncDecl {
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

// ServiceMethod4 generates an ast node that represents the declaration of below:
// For example:
//
//	func (u *user) Create(ctx *types.ServiceContext, user *model.User) (*model.User, error) {\n}
func ServiceMethod4(recvName, modelName, methodName, modelPkgName, reqName, rspName string, body ...ast.Stmt) *ast.FuncDecl {
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
						Names: []*ast.Ident{ast.NewIdent("req")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(modelPkgName),
								Sel: ast.NewIdent(reqName),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(modelPkgName),
								Sel: ast.NewIdent(rspName),
							},
						},
					},
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
