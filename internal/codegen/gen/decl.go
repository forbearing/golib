package gen

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/types/consts"
	"github.com/gertd/go-pluralize"
)

var pluralizeCli = pluralize.NewClient()

// imports generates an ast node that represents the declaration of below:
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

// inits returns an ast node that represents the declaration of below:
/*
func init() {
	service.Register[*user]()
}
*/
func inits(modelNames ...string) *ast.FuncDecl {
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

// types returns an ast node that represents the declaration of below:
/*
type userCreator struct {
	service.Base[*model.User, *model.User, *model.User]
}
*/
func types(modelPkgName, modelName, reqName, rspName string, phase consts.Phase, withComment bool) *ast.GenDecl {
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
				// eg: Creator, Updater, Deleter.
				Name: ast.NewIdent(phase.RoleName()),
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

// serviceMethod1 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *userCreator) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}"
//	"func (g *groupUpdater) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {\n}",
func serviceMethod1(recvName, modelName, modelPkgName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(recvName)},
					Type: &ast.StarExpr{
						X: ast.NewIdent(phase.RoleName()),
					},
				},
			},
		},
		Name: ast.NewIdent(phase.MethodName()),
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

// serviceMethod2 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *userLister) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
//	"func (u *userLister) ListAfter(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
func serviceMethod2(recvName, modelName, modelPkgName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(recvName)},
					Type: &ast.StarExpr{
						X: ast.NewIdent(phase.RoleName()),
					},
				},
			},
		},
		Name: ast.NewIdent(phase.MethodName()),
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

// serviceMethod3 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *userManyCreator) CreateManyBefore(ctx *types.ServiceContext, users ...*model.User) error {\n}"
//	"func (u *userManyCreator) CreateManyAfter(ctx *types.ServiceContext, users ...*model.User) error {\n}"
func serviceMethod3(recvName, modelName, modelPkgName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(recvName)},
					Type: &ast.StarExpr{
						X: ast.NewIdent(phase.RoleName()),
					},
				},
			},
		},
		Name: ast.NewIdent(phase.MethodName()),
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

// serviceMethod4 generates an ast node that represents the declaration of below:
// For example:
//
//	func (u *userCreator) Create(ctx *types.ServiceContext, user *model.User) (rsp *model.User, err error) {\n}
func serviceMethod4(recvName, modelName, modelPkgName, reqName, rspName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
	return &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(recvName)},
					Type: &ast.StarExpr{
						X: ast.NewIdent(phase.RoleName()),
					},
				},
			},
		},
		Name: ast.NewIdent(phase.MethodName()),
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
						Names: []*ast.Ident{ast.NewIdent("rsp")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(modelPkgName),
								Sel: ast.NewIdent(rspName),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("err")},
						Type:  ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}
