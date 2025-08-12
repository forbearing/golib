package gen

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/forbearing/golib/types/consts"
)

// StmtLogInfo create *ast.ExprStmt represents `log.Info(str)`
func StmtLogInfo(str string) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			// log.Info
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("log"),
				Sel: ast.NewIdent("Info"),
			},
			// str
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: str,
				},
			},
		},
	}
}

func EmptyLine() *ast.EmptyStmt {
	return &ast.EmptyStmt{}
}

func Returns(exprs ...ast.Expr) *ast.ReturnStmt {
	return &ast.ReturnStmt{
		Results: exprs,
	}
}

// StmtLogWithServiceContext create *ast.AssignStmt represents `log := u.WithServiceContext(ctx, ctx.GetPhase())`
// modelVarName is model variable name.
func StmtLogWithServiceContext(modelVarName string) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{
			ast.NewIdent("log"),
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   ast.NewIdent(modelVarName),
					Sel: ast.NewIdent("WithServiceContext"),
				},
				Args: []ast.Expr{
					ast.NewIdent("ctx"),
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("ctx"),
							Sel: ast.NewIdent("GetPhase"),
						},
					},
				},
			},
		},
	}
}

// StmtModelRegister creates a *ast.ExprStmt represents golang code like below:
//
//	model.Register[*User]()
func StmtModelRegister(modelName string) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.IndexExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("model"),
					Sel: ast.NewIdent("Register"),
				},
				Index: &ast.StarExpr{
					X: ast.NewIdent(modelName),
				},
			},
		},
	}
}

// StmtServiceRegister creates a *ast.ExprStmt represents golang code like below:
//
//	service.Register[*user.Creator](consts.PHASE_CREATE)
func StmtServiceRegister(structName string, phase consts.Phase) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.IndexExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("service"),
					Sel: ast.NewIdent("Register"),
				},
				Index: &ast.StarExpr{
					X: ast.NewIdent(structName),
				},
			},
			Args: []ast.Expr{
				&ast.SelectorExpr{
					X:   ast.NewIdent("consts"),
					Sel: ast.NewIdent(phase.Name()),
				},
			},
		},
	}
}

// StmtRouterRegister creates a *ast.ExprStmt represents golang code like below:
//
//	router.Register[*model.Group, *model.Group, *model.Group](router.API(), "group")
func StmtRouterRegister(modelPkgName, modelName, reqName, respName string, endpoint string, verb string) *ast.ExprStmt {
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.IndexListExpr{
				X: &ast.SelectorExpr{
					X:   ast.NewIdent("router"),
					Sel: ast.NewIdent("Register"),
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
							Sel: ast.NewIdent(respName),
						},
					},
				},
			},
			Args: []ast.Expr{
				&ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   ast.NewIdent("router"),
						Sel: ast.NewIdent("API"),
					},
				},
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", endpoint),
				},
				&ast.SelectorExpr{
					X:   ast.NewIdent("consts"),
					Sel: ast.NewIdent(verb),
				},
			},
		},
	}
}
