package apply

import (
	"fmt"
	"go/ast"
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
	models, err := codegen.FindModels(config.Module, config.ModelDir, config.ServiceDir, config.Excludes)
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
			if err := generateServiceFile(model); err != nil {
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
			updated, err := applyServiceChanges(model, serviceInfo)
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

// generateServiceFile creates a new service file for the given model
func generateServiceFile(info *gen.ModelInfo) error {
	// Generate the service file using gen package
	file := gen.GenerateService(info)
	if file == nil {
		return fmt.Errorf("failed to generate service file for model %s", info.ModelName)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(info.ServiceFilePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Format the generated code
	code, err := gen.FormatNodeExtra(file)
	if err != nil {
		return fmt.Errorf("failed to format service file: %w", err)
	}

	// Write the file
	if err := os.WriteFile(info.ServiceFilePath, []byte(code), 0o644); err != nil {
		return fmt.Errorf("failed to write service file %s: %w", info.ServiceFilePath, err)
	}

	return nil
}

// applyServiceChanges applies changes to existing service file while preserving business logic
func applyServiceChanges(info *gen.ModelInfo, serviceInfo *ServiceFileInfo) (bool, error) {
	// Find the service struct that inherits from service.Base[*Model]
	serviceStruct := codegen.FindServiceStruct(serviceInfo.File, info.ModelName)
	if serviceStruct == nil {
		return false, fmt.Errorf("service struct inheriting from service.Base[*%s] not found", info.ModelName)
	}

	// Always regenerate all hook methods to ensure consistency
	// This preserves business logic while updating method signatures and structure
	// Only generate methods that are actually missing to avoid duplicates
	var methodsToGenerate []string
	for _, method := range gen.Methods {
		if !codegen.HasMethod(serviceInfo.File, serviceStruct.Name.Name, method) {
			methodsToGenerate = append(methodsToGenerate, method)
		}
	}

	// If no methods are missing, no need to regenerate
	if len(methodsToGenerate) == 0 {
		return false, nil
	}

	// Generate methods and add them to the file
	for _, methodName := range methodsToGenerate {
		// Generate new method with preserved business logic
		method := generateServiceMethod(info, methodName, serviceInfo)
		if method != nil {
			serviceInfo.File.Decls = append(serviceInfo.File.Decls, method)
		}
	}

	// Format and write
	code, err := gen.FormatNodeExtra(serviceInfo.File)
	if err != nil {
		return false, fmt.Errorf("failed to format service file: %w", err)
	}
	if err := os.WriteFile(serviceInfo.FilePath, []byte(code), 0o644); err != nil {
		return false, fmt.Errorf("failed to write service file: %w", err)
	}

	return true, nil
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

	if model.ModelPkgName == "model" {
		// Simple case: *ModelName
		modelTypeExpr = &ast.StarExpr{X: ast.NewIdent(model.ModelName)}
	} else {
		// Qualified case: *model_package.ModelName
		modelTypeExpr = &ast.StarExpr{
			X: &ast.SelectorExpr{
				X:   ast.NewIdent(model.ModelPkgName),
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
