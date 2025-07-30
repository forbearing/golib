package apply

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/forbearing/golib/internal/codegen/gen"
)

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
		stmtStr, err := gen.FormatNodeExtra(stmt)
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
