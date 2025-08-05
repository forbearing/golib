package apply

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/forbearing/golib/internal/codegen/gen"
)

// generateServiceMethod creates a new service method with preserved business logic
func generateServiceMethod(model *gen.ModelInfo, methodName string, serviceInfo *ServiceFileInfo) *ast.FuncDecl {
	// Get existing business logic if available
	var businessLogic []ast.Stmt
	if serviceInfo != nil {
		if existingCode, exists := serviceInfo.BusinessCode[methodName]; exists {
			// Convert business logic strings back to AST statements
			businessLogic = convertBusinessLogicToAST(existingCode)
		}
	}

	// Create method based on type
	var method *ast.FuncDecl
	switch {
	case strings.HasSuffix(methodName, "Before") || strings.HasSuffix(methodName, "After"):
		if strings.Contains(methodName, "List") {
			// List methods use slice parameters
			method = generateListMethod(model, methodName, businessLogic)
		} else if strings.Contains(methodName, "Batch") {
			// Batch methods use variadic parameters
			method = generateBatchMethod(model, methodName, businessLogic)
		} else {
			// Regular CRUD methods
			method = generateCRUDMethod(model, methodName, businessLogic)
		}
	}

	return method
}

// generateCRUDMethod generates regular CRUD methods (Create, Update, Delete, Get)
func generateCRUDMethod(model *gen.ModelInfo, methodName string, businessLogic []ast.Stmt) *ast.FuncDecl {
	modelVarName := strings.ToLower(model.ModelName)

	// Create method body
	body := []ast.Stmt{
		gen.StmtLogWithServiceContext(modelVarName),
		gen.StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(model.ModelName),
			strings.ToLower(methodName))),
		gen.EmptyLine(),
	}

	// Add business logic markers and preserved code
	body = append(body, createBusinessLogicSection(businessLogic)...)
	body = append(body, gen.Returns("nil"))

	return gen.ServiceMethod1(modelVarName, model.ModelName, methodName, model.ModelPkgName, body...)
}

// generateListMethod generates List methods that work with slices
func generateListMethod(model *gen.ModelInfo, methodName string, businessLogic []ast.Stmt) *ast.FuncDecl {
	modelVarName := strings.ToLower(model.ModelName)

	// Create method body
	body := []ast.Stmt{
		gen.StmtLogWithServiceContext(modelVarName),
		gen.StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(model.ModelName),
			strings.ToLower(methodName))),
		gen.EmptyLine(),
	}

	// Add business logic markers and preserved code
	body = append(body, createBusinessLogicSection(businessLogic)...)
	body = append(body, gen.Returns("nil"))

	return gen.ServiceMethod2(modelVarName, model.ModelName, methodName, model.ModelPkgName, body...)
}

// generateBatchMethod generates Batch methods that use variadic parameters
func generateBatchMethod(model *gen.ModelInfo, methodName string, businessLogic []ast.Stmt) *ast.FuncDecl {
	modelVarName := strings.ToLower(model.ModelName)

	// Create method body
	body := []ast.Stmt{
		gen.StmtLogWithServiceContext(modelVarName),
		gen.StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(model.ModelName),
			strings.ToLower(methodName))),
		gen.EmptyLine(),
	}

	// Add business logic markers and preserved code
	body = append(body, createBusinessLogicSection(businessLogic)...)
	body = append(body, gen.Returns("nil"))

	return gen.ServiceMethod3(modelVarName, model.ModelName, methodName, model.ModelPkgName, body...)
}

// createBusinessLogicSection creates the business logic section with markers
func createBusinessLogicSection(businessLogic []ast.Stmt) []ast.Stmt {
	section := []ast.Stmt{
		createCommentStmt("// ===== business logic start ====="),
		gen.EmptyLine(),
	}

	// Add preserved business logic
	section = append(section, businessLogic...)

	section = append(section,
		gen.EmptyLine(),
		createCommentStmt("// ===== business logic end ====="),
	)

	return section
}

// createCommentStmt creates a comment statement
func createCommentStmt(comment string) ast.Stmt {
	return &ast.ExprStmt{
		X: &ast.BasicLit{
			Kind:  0, // This will be handled by the formatter
			Value: comment,
		},
	}
}

// convertBusinessLogicToAST converts business logic strings back to AST statements
func convertBusinessLogicToAST(businessCode []string) []ast.Stmt {
	if len(businessCode) == 0 {
		return nil
	}

	var statements []ast.Stmt
	for _, code := range businessCode {
		trimmed := strings.TrimSpace(code)
		if trimmed == "" {
			statements = append(statements, gen.EmptyLine())
		} else {
			// Parse the code back to proper AST
			// For now, create a raw statement that will be formatted properly
			statements = append(statements, &ast.ExprStmt{
				X: &ast.Ident{Name: trimmed},
			})
		}
	}

	return statements
}
