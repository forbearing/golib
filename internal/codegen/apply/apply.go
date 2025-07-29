package apply

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"
	"slices"

	"github.com/forbearing/golib/internal/codegen"
	"github.com/forbearing/golib/internal/codegen/gen"
)

// ApplyServiceGeneration applies service generation with business logic preservation
func ApplyServiceGeneration(config *ApplyConfig) error {
	models, err := codegen.FindModelsInDirectory(config.Module, config.ModelDir, config.ServiceDir, config.Excludes)
	if err != nil {
		return fmt.Errorf("failed to find models: %w", err)
	}

	// Scan existing service files and set ServiceFilePath for each model
	serviceFiles, err := scanServiceFiles(config.ServiceDir, models)
	if err != nil {
		return fmt.Errorf("failed to scan service files: %w", err)
	}

	// Process each model
	processedCount := 0
	skippedCount := 0
	for _, model := range models {

		if slices.Contains(config.Excludes, model.ModelName) {
			skippedCount++
			continue
		}

		serviceInfo := serviceFiles[model.ModelName]

		// Check if regeneration is needed
		if !needsRegeneration(model, serviceInfo) {
			skippedCount++
			continue
		}

		// Generate service file
		err := generateServiceFile(model, serviceInfo, config)
		if err != nil {
			return fmt.Errorf("failed to generate service for %s: %w", model.ModelName, err)
		}

		processedCount++
	}

	return nil
}

// generateServiceFile generates a complete service file
func generateServiceFile(model *gen.ModelInfo, serviceInfo *ServiceFileInfo, config *ApplyConfig) error {
	// Create file set and AST file
	fset := token.NewFileSet()
	file := &ast.File{
		Name: ast.NewIdent(gen.ModelPkg2ServicePkg(model.PackageName)),
	}

	// Add imports
	file.Decls = append(file.Decls, gen.Imports(model.ModulePath, model.ModelFileDir, model.PackageName))

	// Add type declarations
	file.Decls = append(file.Decls, gen.Types(model.ModelName, model.PackageName))

	// Add init function
	file.Decls = append(file.Decls, gen.Inits(model.ModelName))

	// Generate all service methods
	methods := []string{
		"CreateBefore", "CreateAfter",
		"UpdateBefore", "UpdateAfter",
		"DeleteBefore", "DeleteAfter",
		"GetBefore", "GetAfter",
		"ListBefore", "ListAfter",
		"BatchCreateBefore", "BatchCreateAfter",
		"BatchUpdateBefore", "BatchUpdateAfter",
		"BatchDeleteBefore", "BatchDeleteAfter",
	}

	generatedMethods := 0
	for _, methodName := range methods {
		method := generateServiceMethod(model, methodName, serviceInfo)
		if method != nil {
			file.Decls = append(file.Decls, method)
			generatedMethods++
		} else {
		}
	}

	// Ensure service directory exists
	serviceDir := filepath.Dir(model.ServiceFilePath)

	if err := os.MkdirAll(serviceDir, 0o755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Write file
	outputFile, err := os.Create(model.ServiceFilePath)
	if err != nil {
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer outputFile.Close()

	// Format and write
	if err := format.Node(outputFile, fset, file); err != nil {
		return fmt.Errorf("failed to format service file: %w", err)
	}

	return nil
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
