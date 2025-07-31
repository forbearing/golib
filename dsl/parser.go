package dsl

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/kr/pretty"
)

func Parse(stmts []ast.Stmt) *Design {
	design := &Design{}
	if len(stmts) == 0 {
		return design
	}

	for _, stmt := range stmts {
		callExpr, ok := stmt.(*ast.ExprStmt)
		if !ok || callExpr == nil {
			continue
		}
		call, ok := callExpr.X.(*ast.CallExpr)
		if !ok || call == nil || call.Fun == nil || len(call.Args) == 0 {
			continue
		}
		ident, ok := call.Fun.(*ast.Ident)
		if !ok || ident == nil {
			continue
		}
		if !is(ident.Name) {
			continue
		}

		// Parse "Enabled" design.
		if ident.Name == "Enabled" && len(call.Args) == 1 {
			arg, ok := call.Args[0].(*ast.Ident)
			if ok && arg != nil {
				design.Enabled = arg.Name == "true"
			}
		}

		// Parse "Endpoint" design
		if ident.Name == "Endpoint" && len(call.Args) == 1 {
			if arg, ok := call.Args[0].(*ast.BasicLit); ok && arg != nil && arg.Kind == token.STRING {
				design.Endpoint = trimQuote(arg.Value)
			}
		}

		// Parse "Create" design.
		if ident.Name == "Create" && len(call.Args) == 1 {
			if flit, ok := call.Args[0].(*ast.FuncLit); ok && flit != nil && flit.Body != nil {
				// Payload or Result.
				for _, stmt := range flit.Body.List {
					if expr, ok := stmt.(*ast.ExprStmt); ok && expr != nil {
						if call, ok := expr.X.(*ast.CallExpr); ok && call != nil && call.Fun != nil {
							if identExpr, ok := call.Fun.(*ast.IndexExpr); ok && identExpr != nil {
								var isPayload bool
								var isResult bool
								_ = isResult
								if ident, ok := identExpr.X.(*ast.Ident); ok && ident != nil {
									if ident.Name == "Payload" {
										isPayload = true
									} else if ident.Name == "Result" {
										isResult = true
									}
								}
								if isPayload {
									// Payload[User]
									if ident, ok := identExpr.Index.(*ast.Ident); ok && ident != nil {
										fmt.Println("----- Payload", ident.Name)
									}
									// Payload[*User]
									if starExpr, ok := identExpr.Index.(*ast.StarExpr); ok && starExpr != nil {
										if ident, ok := starExpr.X.(*ast.Ident); ok && ident != nil {
											fmt.Println("----- Payload", ident.Name)
										}
									}
								}
							}
						}
					}
				}
			}
		}

	}
	pretty.Println(design)

	return design
}
