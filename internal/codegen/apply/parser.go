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
		return nil, nil // File doesn't exist, not an error
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

// getMethodNames extracts method names from a map of function declarations
// TODO: 未来可能需要用于提取方法名列表的功能，暂时保留
func getMethodNames(methods map[string]*ast.FuncDecl) []string {
	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	return names
}

// scanServiceFiles scans the service directory for existing service files
// TODO: 未来可能需要用于批量扫描服务文件的功能，暂时保留
func scanServiceFiles(serviceDir string, models []*gen.ModelInfo) (map[string]*ServiceFileInfo, error) {
	serviceFiles := make(map[string]*ServiceFileInfo)
	scannedCount := 0
	foundCount := 0

	for _, model := range models {
		scannedCount++

		// Generate service file path based on model info
		// Extract the subdirectory from model file path to maintain directory structure
		modelRelPath, err := filepath.Rel(filepath.Join(serviceDir, "..", "model"), model.ModelFileDir)
		if err != nil {
			return nil, err
		}

		serviceFilePath := filepath.Join(serviceDir, modelRelPath, strings.ToLower(model.ModelName)+".go")

		// Update model with service file path
		model.ServiceFilePath = serviceFilePath

		// Parse the service file if it exists
		serviceInfo, err := parseServiceFile(serviceFilePath)
		if err != nil {
			return nil, err
		}

		if serviceInfo != nil {
			// Use model name as key
			serviceFiles[model.ModelName] = serviceInfo
			foundCount++
		}
	}

	return serviceFiles, nil
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
