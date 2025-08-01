package dsl

import (
	"go/ast"
	"go/token"
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

		if payload, result, exists := parseAction("Create", ident, call.Args); exists {
			design.Create = Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction("Update", ident, call.Args); exists {
			design.Update = Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction("UpdatePartial", ident, call.Args); exists {
			design.UpdatePartial = Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction("BatchCreate", ident, call.Args); exists {
			design.BatchCreate = Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction("BatchUpdate", ident, call.Args); exists {
			design.BatchUpdate = Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction("BatchUpdatePartial", ident, call.Args); exists {
			design.BatchUpdatePartial = Action{Payload: payload, Result: result}
		}

	}

	return design
}

func parseAction(name string, ident *ast.Ident, args []ast.Expr) (string, string, bool) {
	var payload string
	var result string

	if ident.Name == name && len(args) == 1 {
		if flit, ok := args[0].(*ast.FuncLit); ok && flit != nil && flit.Body != nil {
			// Payload or Result.
			for _, stmt := range flit.Body.List {
				if expr, ok := stmt.(*ast.ExprStmt); ok && expr != nil {
					if call, ok := expr.X.(*ast.CallExpr); ok && call != nil && call.Fun != nil {
						if identExpr, ok := call.Fun.(*ast.IndexExpr); ok && identExpr != nil {
							var isPayload bool
							var isResult bool
							_ = isResult
							if ident, ok := identExpr.X.(*ast.Ident); ok && ident != nil {
								switch ident.Name {
								case "Payload":
									isPayload = true
								case "Result":
									isResult = true
								}
							}
							if isPayload {
								if ident, ok := identExpr.Index.(*ast.Ident); ok && ident != nil { // Payload[User]
									payload = ident.Name
								} else if starExpr, ok := identExpr.Index.(*ast.StarExpr); ok && starExpr != nil { // Payload[*User]
									if ident, ok := starExpr.X.(*ast.Ident); ok && ident != nil {
										payload = ident.Name
									}
								}
							}
							if isResult {
								if ident, ok := identExpr.Index.(*ast.Ident); ok && ident != nil { // Result[User]
									result = ident.Name
								} else if starExpr, ok := identExpr.Index.(*ast.StarExpr); ok && starExpr != nil { // Result[*User]
									if ident, ok := starExpr.X.(*ast.Ident); ok && ident != nil {
										result = ident.Name
									}
								}
							}
						}
					}
				}
			}
		}
		return payload, result, true
	}

	return "", "", false
}
