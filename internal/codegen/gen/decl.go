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
func imports(modulePath, modelFileDir, modelPkgName string, otherPkg ...string) *ast.GenDecl {
	importModel := filepath.Join(modulePath, modelFileDir)
	fields := strings.Split(importModel, "/")
	if len(fields) > 0 && fields[len(fields)-1] != modelPkgName {
		// model_setting "mymodule/model/setting"
		importModel = fmt.Sprintf("%s %q", modelPkgName, importModel)
	} else {
		// "mymodule/model"
		importModel = fmt.Sprintf("%q", importModel)
	}

	genDecl := &ast.GenDecl{
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

	for _, pkg := range otherPkg {
		if len(pkg) == 0 {
			continue
		}
		genDecl.Specs = append(genDecl.Specs, &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: fmt.Sprintf("%q", pkg),
			},
		})
	}

	return genDecl
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

	// if reqName is equal to modelName or reqName starts with *, then the reqExpr use StarExpr,
	// otherwise use SelectorExpr
	var reqExpr ast.Expr
	if strings.HasPrefix(reqName, "*") || modelName == reqName {
		reqExpr = &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(modelPkgName),
				Sel: ast.NewIdent(strings.TrimPrefix(reqName, "*")),
			},
		}
	} else {
		reqExpr = &ast.SelectorExpr{
			X:   ast.NewIdent(modelPkgName),
			Sel: ast.NewIdent(reqName),
		}
	}

	// if rspName is equal to modelName or rspName starts with *, then the rspExpr use StarExpr,
	// otherwise use SelectorExpr
	var rspExpr ast.Expr
	if strings.HasPrefix(rspName, "*") || modelName == rspName {
		rspExpr = &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(modelPkgName),
				Sel: ast.NewIdent(strings.TrimPrefix(rspName, "*")),
			},
		}
	} else {
		rspExpr = &ast.SelectorExpr{
			X:   ast.NewIdent(modelPkgName),
			Sel: ast.NewIdent(rspName),
		}
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
										reqExpr,
										rspExpr,
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
//	"func (u *Creator) CreateBefore(ctx *types.ServiceContext, user *model.User) error {\n}"
//	"func (g *Updater) UpdateAfter(ctx *types.ServiceContext, group *model.Group) error {\n}",
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
//	"func (u *Lister) ListBefore(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
//	"func (u *Lister) ListAfter(ctx *types.ServiceContext, users *[]*model.User) error {\n}"
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
//	"func (u *ManyCreator) CreateManyBefore(ctx *types.ServiceContext, users ...*model.User) error {\n}"
//	"func (u *ManyCreator) CreateManyAfter(ctx *types.ServiceContext, users ...*model.User) error {\n}"
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
//	func (u *Creator) Create(ctx *types.ServiceContext, user *model.User) (rsp *model.User, err error) {\n}
func serviceMethod4(recvName, modelName, modelPkgName, reqName, rspName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
	// if reqName is equal to modelName or reqName starts with *, then the reqExpr use StarExpr,
	// otherwise use SelectorExpr
	var reqExpr ast.Expr
	if strings.HasPrefix(reqName, "*") || modelName == reqName {
		reqExpr = &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(modelPkgName),
				Sel: ast.NewIdent(strings.TrimPrefix(reqName, "*")),
			},
		}
	} else {
		reqExpr = &ast.SelectorExpr{
			X:   ast.NewIdent(modelPkgName),
			Sel: ast.NewIdent(reqName),
		}
	}

	// if rspName is equal to modelName or rspName starts with *, then the rspExpr use StarExpr,
	// otherwise use SelectorExpr
	var rspExpr ast.Expr
	if strings.HasPrefix(rspName, "*") || modelName == rspName {
		rspExpr = &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(modelPkgName),
				Sel: ast.NewIdent(strings.TrimPrefix(rspName, "*")),
			},
		}
	} else {
		rspExpr = &ast.SelectorExpr{
			X:   ast.NewIdent(modelPkgName),
			Sel: ast.NewIdent(rspName),
		}
	}

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
						Type:  reqExpr,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("rsp")},
						Type:  rspExpr,
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

// serviceMethod5 generates an ast node that represents the declaration of below:
// For example:
//
//	func (a *Importer) Import(ctx *types.ServiceContext, reader io.Reader) ([]*model.Asset, error) {\n}
func serviceMethod5(recvName, modelName, modelPkgName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
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
		Name: ast.NewIdent("Import"),
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
						Names: []*ast.Ident{ast.NewIdent("reader")},
						Type: &ast.SelectorExpr{
							X:   ast.NewIdent("io"),
							Sel: ast.NewIdent("Reader"),
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent(pluralizeCli.Plural(strings.ToLower(modelName)))},
						Type: &ast.ArrayType{
							Elt: &ast.StarExpr{
								X: &ast.SelectorExpr{
									X:   ast.NewIdent(modelPkgName),
									Sel: ast.NewIdent(modelName),
								},
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

// serviceMethod6 generates an ast node that represents the declaration of below:
// For example:
//
//	func (a *Exporter) Export(ctx *types.ServiceContext, data ...*model.Asset) ([]byte, error) {\n}
func serviceMethod6(recvName, modelName, modelPkgName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
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
		Name: ast.NewIdent("Export"),
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
						Names: []*ast.Ident{ast.NewIdent("data")},
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
						Names: []*ast.Ident{ast.NewIdent("data")},
						Type: &ast.ArrayType{
							Elt: ast.NewIdent("byte"),
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

// serviceMethod7 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *Lister) Filter(ctx *types.ServiceContext, user *model.User) *model.User {\n}"
//	"func (g *Lister) Filter(ctx *types.ServiceContext, group *model_auth.Group) *model_auth.Group {\n}"
func serviceMethod7(recvName, modelName, modelPkgName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
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
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent(modelPkgName),
								Sel: ast.NewIdent(modelName),
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}

// serviceMethod8 generates an ast node that represents the declaration of below:
// For example:
//
//	"func (u *Lister) FilterRaw(ctx *types.ServiceContext) string {\n}"
//	"func (g *Lister) FilterRaw(ctx *types.ServiceContext) string {\n}"
func serviceMethod8(recvName, modelName, modelPkgName string, phase consts.Phase, body ...ast.Stmt) *ast.FuncDecl {
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
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("string"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: body,
		},
	}
}
