package apply

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"os"
	"path/filepath"

	"github.com/forbearing/golib/internal/codegen/gen"
	"github.com/sirupsen/logrus"
)

var (
	// logger is the dedicated logger for apply package with caller information
	logger *logrus.Entry
)

func init() {
	logger = logrus.WithField("component", "codegen-apply")
}

// extractModelNames extracts model names from ModelInfo slice for logging
func extractModelNames(models []*gen.ModelInfo) []string {
	names := make([]string, len(models))
	for i, model := range models {
		names[i] = model.ModelName
	}
	return names
}

// ApplyServiceGeneration applies service generation with business logic preservation
func ApplyServiceGeneration(config *ApplyConfig) error {
	logger.WithFields(logrus.Fields{
		"module":     config.Module,
		"modelDir":   config.ModelDir,
		"serviceDir": config.ServiceDir,
		"excludes":   config.Excludes,
	}).Info("Starting service generation")

	// Find all models by scanning directory
	models, err := FindModelsInDirectory(config.Module, config.ModelDir)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":    err,
			"modelDir": config.ModelDir,
		}).Error("Failed to find models")
		return fmt.Errorf("failed to find models: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"count":  len(models),
		"models": extractModelNames(models),
	}).Info("Found models for processing")

	// Scan existing service files and set ServiceFilePath for each model
	serviceFiles, err := scanServiceFiles(config.ServiceDir, models)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error":      err,
			"serviceDir": config.ServiceDir,
		}).Error("Failed to scan service files")
		return fmt.Errorf("failed to scan service files: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"serviceFilesCount": len(serviceFiles),
	}).Info("Scanned existing service files")

	// Process each model
	processedCount := 0
	skippedCount := 0
	for _, model := range models {
		modelLogger := logger.WithFields(logrus.Fields{
			"modelName":       model.ModelName,
			"packageName":     model.PackageName,
			"serviceFilePath": model.ServiceFilePath,
		})

		if shouldSkipModel(model.ModelName, config.Excludes) {
			modelLogger.Info("Skipping model due to exclusion")
			skippedCount++
			continue
		}

		serviceInfo := serviceFiles[model.ModelName]

		// Check if regeneration is needed
		if !needsRegeneration(model, serviceInfo) {
			modelLogger.Info("Service file is up to date, skipping regeneration")
			skippedCount++
			continue
		}

		modelLogger.Info("Generating service file")

		// Generate service file
		err := generateServiceFile(model, serviceInfo, config)
		if err != nil {
			modelLogger.WithField("error", err).Error("Failed to generate service file")
			return fmt.Errorf("failed to generate service for %s: %w", model.ModelName, err)
		}

		modelLogger.Info("Successfully generated service file")
		processedCount++
	}

	logger.WithFields(logrus.Fields{
		"totalModels": len(models),
		"processed":   processedCount,
		"skipped":     skippedCount,
	}).Info("Service generation completed")

	return nil
}

// FindModelsInDirectory finds all models in a directory
func FindModelsInDirectory(modulePath, modelDir string) ([]*gen.ModelInfo, error) {
	logger.WithFields(logrus.Fields{
		"modulePath": modulePath,
		"modelDir":   modelDir,
	}).Info("Starting model discovery")

	var allModels []*gen.ModelInfo

	err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.WithFields(logrus.Fields{
				"path":  path,
				"error": err,
			}).Warn("Error accessing path during walk")
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		// Skip test files
		if filepath.Base(path) != "go" {
			if len(filepath.Base(path)) > 8 && filepath.Base(path)[len(filepath.Base(path))-8:] == "_test.go" {
				logger.WithField("path", path).Debug("Skipping test file")
				return nil
			}
		}

		logger.WithField("path", path).Debug("Processing Go file")

		// Find models in this file
		models, err := gen.FindModels(modulePath, path)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"path":  path,
				"error": err,
			}).Debug("No valid models found in file")
			// Skip files that don't contain valid models
			return nil
		}

		if len(models) > 0 {
			modelNames := extractModelNames(models)
			logger.WithFields(logrus.Fields{
				"path":   path,
				"count":  len(models),
				"models": modelNames,
			}).Info("Found models in file")
			allModels = append(allModels, models...)
		}

		return nil
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"modelDir": modelDir,
			"error":    err,
		}).Error("Error during model directory walk")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"totalModels": len(allModels),
		"modelNames":  extractModelNames(allModels),
	}).Info("Model discovery completed")

	return allModels, err
}

// generateServiceFile generates a complete service file
func generateServiceFile(model *gen.ModelInfo, serviceInfo *ServiceFileInfo, config *ApplyConfig) error {
	logger.WithFields(logrus.Fields{
		"modelName":          model.ModelName,
		"serviceFilePath":    model.ServiceFilePath,
		"hasExistingService": serviceInfo != nil,
	}).Info("Starting service file generation")

	// Create file set and AST file
	fset := token.NewFileSet()
	file := &ast.File{
		Name: ast.NewIdent(gen.ModelPkg2ServicePkg(model.PackageName)),
	}

	logger.WithField("packageName", gen.ModelPkg2ServicePkg(model.PackageName)).Debug("Created AST file structure")

	// Add imports
	file.Decls = append(file.Decls, gen.Imports(model.ModulePath, model.ModelFileDir, model.PackageName))
	logger.Debug("Added imports to service file")

	// Add type declarations
	file.Decls = append(file.Decls, gen.Types(model.ModelName, model.PackageName))
	logger.Debug("Added type declarations to service file")

	// Add init function
	file.Decls = append(file.Decls, gen.Inits(model.ModelName))
	logger.Debug("Added init function to service file")

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

	logger.WithFields(logrus.Fields{
		"methodCount": len(methods),
		"methods":     methods,
	}).Info("Generating service methods")

	generatedMethods := 0
	for _, methodName := range methods {
		method := generateServiceMethod(model, methodName, serviceInfo)
		if method != nil {
			file.Decls = append(file.Decls, method)
			generatedMethods++
			logger.WithField("methodName", methodName).Debug("Generated service method")
		} else {
			logger.WithField("methodName", methodName).Warn("Failed to generate service method")
		}
	}

	logger.WithFields(logrus.Fields{
		"totalMethods":     len(methods),
		"generatedMethods": generatedMethods,
	}).Info("Service methods generation completed")

	// Ensure service directory exists
	serviceDir := filepath.Dir(model.ServiceFilePath)
	logger.WithField("serviceDir", serviceDir).Debug("Ensuring service directory exists")

	if err := os.MkdirAll(serviceDir, 0o755); err != nil {
		logger.WithFields(logrus.Fields{
			"serviceDir": serviceDir,
			"error":      err,
		}).Error("Failed to create service directory")
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Write file
	logger.WithField("filePath", model.ServiceFilePath).Debug("Creating service file")
	outputFile, err := os.Create(model.ServiceFilePath)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"filePath": model.ServiceFilePath,
			"error":    err,
		}).Error("Failed to create service file")
		return fmt.Errorf("failed to create service file: %w", err)
	}
	defer outputFile.Close()

	// Format and write
	logger.Debug("Formatting and writing service file")
	if err := format.Node(outputFile, fset, file); err != nil {
		logger.WithFields(logrus.Fields{
			"filePath": model.ServiceFilePath,
			"error":    err,
		}).Error("Failed to format service file")
		return fmt.Errorf("failed to format service file: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"modelName":        model.ModelName,
		"serviceFilePath":  model.ServiceFilePath,
		"generatedMethods": generatedMethods,
	}).Info("Service file generation completed successfully")

	return nil
}

// shouldSkipModel checks if a model should be skipped based on exclusions
func shouldSkipModel(modelName string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if modelName == exclusion {
			return true
		}
	}
	return false
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

