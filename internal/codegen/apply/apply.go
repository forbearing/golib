package apply

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/forbearing/golib/internal/codegen"
	"github.com/forbearing/golib/internal/codegen/gen"
)

// ApplyServiceGeneration applies service code generation with business logic protection
func ApplyServiceGeneration(config *ApplyConfig) error {
	// Find all models in the model directory
	models, err := codegen.FindModelsInDirectory(config.Module, config.ModelDir, config.ServiceDir, config.Excludes)
	if err != nil {
		return fmt.Errorf("failed to find models: %w", err)
	}

	for _, model := range models {
		// Skip excluded models
		if slices.Contains(config.Excludes, model.ModelName) {
			continue
		}

		// Check if service file exists
		if _, err := os.Stat(model.ServiceFilePath); os.IsNotExist(err) {
			// Create new service file
			if err := generateNewServiceFile(model, config); err != nil {
				return fmt.Errorf("failed to generate new service file for %s: %w", model.ModelName, err)
			}
			fmt.Printf("generate %s\n", model.ServiceFilePath)
		} else {
			// Parse existing service file
			serviceInfo, err := parseServiceFile(model.ServiceFilePath)
			if err != nil {
				return fmt.Errorf("failed to parse service file %s: %w", model.ServiceFilePath, err)
			}

			// Apply changes to existing service file
			updated, err := applyServiceChanges(model, serviceInfo, config)
			if err != nil {
				return fmt.Errorf("failed to apply service changes for %s: %w", model.ModelName, err)
			}

			if updated {
				fmt.Printf("update %s\n", model.ServiceFilePath)
			} else {
				fmt.Printf("skip %s\n", model.ServiceFilePath)
			}
		}
	}

	return nil
}

// generateNewServiceFile creates a new service file for the given model
func generateNewServiceFile(model *gen.ModelInfo, config *ApplyConfig) error {
	// Generate the service file using gen package
	serviceFile := gen.GenerateServiceFile(model)
	if serviceFile == nil {
		return fmt.Errorf("failed to generate service file for model %s", model.ModelName)
	}

	// Format the generated code
	var buf strings.Builder
	fset := token.NewFileSet()
	if err := format.Node(&buf, fset, serviceFile); err != nil {
		return fmt.Errorf("failed to format service file: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(model.ServiceFilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write the file
	if err := os.WriteFile(model.ServiceFilePath, []byte(buf.String()), 0o644); err != nil {
		return fmt.Errorf("failed to write service file %s: %w", model.ServiceFilePath, err)
	}

	return nil
}

// applyServiceChanges applies changes to existing service file while preserving business logic
func applyServiceChanges(model *gen.ModelInfo, serviceInfo *ServiceFileInfo, config *ApplyConfig) (bool, error) {
	// Find the service struct that inherits from service.Base[*Model]
	serviceStruct := codegen.FindServiceStruct(serviceInfo.File, model.ModelName)
	if serviceStruct == nil {
		return false, fmt.Errorf("service struct inheriting from service.Base[*%s] not found", model.ModelName)
	}

	// Get all hook methods that should exist
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

	// Check which methods are missing and need to be added
	missingMethods := []string{}
	for _, methodName := range requiredMethods {
		if !hasMethod(serviceInfo.File, serviceStruct.Name.Name, methodName) {
			missingMethods = append(missingMethods, methodName)
		}
	}

	// If no methods are missing, no need to update
	if len(missingMethods) == 0 {
		return false, nil
	}

	// Generate missing methods and add them to the file
	for _, methodName := range missingMethods {
		method := generateHookMethod(model, methodName, serviceStruct.Name.Name)
		if method != nil {
			serviceInfo.File.Decls = append(serviceInfo.File.Decls, method)
		}
	}

	// Write the updated file
	outputFile, err := os.Create(serviceInfo.FilePath)
	if err != nil {
		return false, fmt.Errorf("failed to create service file: %w", err)
	}
	defer outputFile.Close()

	// Format and write
	if err := format.Node(outputFile, serviceInfo.FileSet, serviceInfo.File); err != nil {
		return false, fmt.Errorf("failed to format service file: %w", err)
	}

	return true, nil
}

// hasMethod checks if a struct has a specific method
func hasMethod(file *ast.File, structName, methodName string) bool {
	return codegen.HasMethod(file, structName, methodName)
}

// generateHookMethod generates a hook method for the service struct
func generateHookMethod(model *gen.ModelInfo, methodName, structName string) *ast.FuncDecl {
	// Create receiver
	recv := &ast.FieldList{
		List: []*ast.Field{
			{
				Names: []*ast.Ident{ast.NewIdent("s")},
				Type:  &ast.StarExpr{X: ast.NewIdent(structName)},
			},
		},
	}

	// Create the model type expression based on package name
	// The ModelInfo.PackageName contains the model package name (e.g., "model_cmdb")
	var modelTypeExpr ast.Expr

	if model.PackageName == "model" {
		// Simple case: *ModelName
		modelTypeExpr = &ast.StarExpr{X: ast.NewIdent(model.ModelName)}
	} else {
		// Qualified case: *model_package.ModelName
		modelTypeExpr = &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(model.PackageName),
				Sel: ast.NewIdent(model.ModelName),
			},
		}
	}

	// Create parameters based on method type
	var params *ast.FieldList
	switch {
	case strings.HasSuffix(methodName, "Before") || strings.HasSuffix(methodName, "After"):
		if strings.HasPrefix(methodName, "List") {
			// List methods take *[]Model
			params = &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("types"),
								Sel: ast.NewIdent("ServiceContext"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("models")},
						Type: &ast.StarExpr{
							X: &ast.ArrayType{
								Elt: modelTypeExpr,
							},
						},
					},
				},
			}
		} else if strings.HasPrefix(methodName, "Batch") {
			// Batch methods take ...Model
			params = &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("types"),
								Sel: ast.NewIdent("ServiceContext"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("models")},
						Type: &ast.Ellipsis{
							Elt: modelTypeExpr,
						},
					},
				},
			}
		} else {
			// Regular methods take Model
			params = &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("types"),
								Sel: ast.NewIdent("ServiceContext"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent("model")},
						Type:  modelTypeExpr,
					},
				},
			}
		}
	}

	// Create return type (error)
	results := &ast.FieldList{
		List: []*ast.Field{
			{
				Type: ast.NewIdent("error"),
			},
		},
	}

	// Create method body (return nil)
	body := &ast.BlockStmt{
		List: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{ast.NewIdent("nil")},
			},
		},
	}

	return &ast.FuncDecl{
		Recv: recv,
		Name: ast.NewIdent(methodName),
		Type: &ast.FuncType{
			Params:  params,
			Results: results,
		},
		Body: body,
	}
}

// NewApplyConfig creates a new apply configuration
func NewApplyConfig(modulePath, modelDir, serviceDir string) *ApplyConfig {
	return &ApplyConfig{
		Module:     modulePath,
		ModelDir:   modelDir,
		ServiceDir: serviceDir,
	}
}

// WithExclusions adds model exclusions to the config
func (c *ApplyConfig) WithExclusions(exclusions ...string) *ApplyConfig {
	c.Excludes = append(c.Excludes, exclusions...)
	return c
}
