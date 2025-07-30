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

// needsRegeneration checks if a service file needs to be regenerated
func needsRegeneration(model *gen.ModelInfo, serviceInfo *ServiceFileInfo) bool {
	if serviceInfo == nil {
		return true // File doesn't exist, needs generation
	}

	// Check if all required methods exist
	requiredMethods := []string{
		"CreateBefore", "CreateAfter",
		"UpdateBefore", "UpdateAfter",
		"DeleteBefore", "DeleteAfter",
		"GetBefore", "GetAfter",
		"ListBefore", "ListAfter",
		"BatchCreateBefore", "BatchCreateAfter",
		"BatchUpdateBefore", "BatchUpdateAfter",
		"BatchDeleteBefore", "BatchDeleteAfter",
	}

	missingMethods := []string{}
	for _, methodName := range requiredMethods {
		if _, exists := serviceInfo.Methods[methodName]; !exists {
			missingMethods = append(missingMethods, methodName)
		}
	}

	if len(missingMethods) > 0 {
		return true
	}

	return false // All methods exist, no regeneration needed
}
