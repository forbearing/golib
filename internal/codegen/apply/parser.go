package apply

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/internal/codegen/gen"
)

// parseServiceFile parses an existing service file and extracts information
func parseServiceFile(filePath string) (*ServiceFileInfo, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &ServiceFileInfo{}, nil // File doesn't exist, not an error
	}

	// Parse the file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	info := &ServiceFileInfo{
		FilePath:     filePath,
		FileSet:      fset,
		File:         file,
		Methods:      make(map[string]*ast.FuncDecl),
		BusinessCode: make(map[string][]string),
	}

	// Extract model name from file name
	baseName := filepath.Base(filePath)
	// TODO: model name is same as file name?
	modelName := strings.TrimSuffix(baseName, ".go")
	info.ModelName = strings.Title(modelName)
	info.PackageName = file.Name.Name

	methodCount := 0
	// Extract existing methods and their business logic
	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Check if this is a service method (has receiver)
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				methodName := funcDecl.Name.Name
				info.Methods[methodName] = funcDecl
				methodCount++

				// Extract business logic
				businessLogic := extractBusinessLogic(funcDecl.Body, fset)
				if len(businessLogic) > 0 {
					info.BusinessCode[methodName] = businessLogic
				}
			}
		}
	}

	return info, nil
}

// extractBusinessLogic extracts business logic from method body between markers
func extractBusinessLogic(methodBody *ast.BlockStmt, fset *token.FileSet) []string {
	if methodBody == nil || len(methodBody.List) == 0 {
		return nil
	}

	var businessLogic []string
	// inBusinessSection := false

	for _, stmt := range methodBody.List {
		code, err := gen.FormatNode(stmt)
		if err != nil {
			continue
		}
		businessLogic = append(businessLogic, code)

		// // Check if this is a comment statement
		// if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
		// 	if callExpr, ok := exprStmt.X.(*ast.CallExpr); ok {
		// 		// Skip log statements and other generated code
		// 		if isGeneratedStatement(callExpr) {
		// 			continue
		// 		}
		// 	}
		// }
		//
		// // Convert statement to string
		// stmtStr, err := gen.FormatNode(stmt)
		// if err != nil {
		// 	continue
		// }
		//
		// // Check for business logic markers
		// if strings.Contains(stmtStr, "// ===== business logic start =====") {
		// 	inBusinessSection = true
		// 	continue
		// }
		// if strings.Contains(stmtStr, "// ===== business logic end =====") {
		// 	inBusinessSection = false
		// 	continue
		// }
		//
		// // Collect business logic
		// if inBusinessSection {
		// 	businessLogic = append(businessLogic, stmtStr)
		// }
	}

	return businessLogic
}

// // isGeneratedStatement checks if a statement is generated code that should be skipped
// func isGeneratedStatement(callExpr *ast.CallExpr) bool {
// 	if sel, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
// 		// Skip log.Info statements
// 		if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == "log" && sel.Sel.Name == "Info" {
// 			return true
// 		}
// 	}
// 	return false
// }
