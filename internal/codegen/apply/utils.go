package apply

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

// formatNode formats an AST node to Go source code
func formatNode(node ast.Node) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// extractBusinessLogic extracts business logic from method body between markers
func extractBusinessLogic(methodBody *ast.BlockStmt, fset *token.FileSet) []string {
	if methodBody == nil || len(methodBody.List) == 0 {
		return nil
	}

	var businessLogic []string
	inBusinessSection := false
	
	for _, stmt := range methodBody.List {
		// Check if this is a comment statement
		if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
			if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
				// Skip log statements and other generated code
				if isGeneratedStatement(callExpr) {
					continue
				}
			}
		}
		
		// Convert statement to string
		stmtStr, err := formatNode(stmt)
		if err != nil {
			continue
		}
		
		// Check for business logic markers
		if strings.Contains(stmtStr, "// ===== business logic start =====") {
			inBusinessSection = true
			continue
		}
		if strings.Contains(stmtStr, "// ===== business logic end =====") {
			inBusinessSection = false
			continue
		}
		
		// Collect business logic
		if inBusinessSection {
			businessLogic = append(businessLogic, stmtStr)
		}
	}
	
	return businessLogic
}

// isGeneratedStatement checks if a statement is generated code that should be skipped
func isGeneratedStatement(callExpr *ast.CallExpr) bool {
	if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
		// Skip log.Info statements
		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "log" && sel.Sel.Name == "Info" {
			return true
		}
	}
	return false
}

// findMethodByName finds a method declaration by name in the AST file
func findMethodByName(file *ast.File, methodName string) *ast.FuncDecl {
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == methodName {
				return funcDecl
			}
		}
	}
	return nil
}

// hasBusinessLogicMarkers checks if a method has business logic markers
func hasBusinessLogicMarkers(methodBody *ast.BlockStmt, fset *token.FileSet) bool {
	if methodBody == nil {
		return false
	}
	
	bodyStr, err := formatNode(methodBody)
	if err != nil {
		return false
	}
	
	return strings.Contains(bodyStr, "// ===== business logic start =====") &&
		   strings.Contains(bodyStr, "// ===== business logic end =====")
}