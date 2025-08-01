package codegen

import (
	"go/ast"
	"go/token"
)

// expr_log_info create *ast.ExprStmt represents `log.Info(str)`
func expr_log_info(str string) *ast.ExprStmt {
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

func empty_line() *ast.EmptyStmt {
	return &ast.EmptyStmt{}
}

func empty_line2() *ast.ExprStmt {
	return &ast.ExprStmt{}
}

func returns(strs ...string) *ast.ReturnStmt {
	exprs := make([]ast.Expr, 0, len(strs))
	for _, str := range strs {
		if len(str) == 0 {
			continue
		}
		exprs = append(exprs, ast.NewIdent(str))
	}
	return &ast.ReturnStmt{
		Results: exprs,
	}
}

// assign_with_service_context create *ast.AssignStmt represents `log := u.WithServiceContext(ctx, ctx.GetPhase())`
// modelVarName is model variable name.
func assign_with_service_context(modelVarName string) *ast.AssignStmt {
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
