package apply

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/internal/codegen/gen"
	"github.com/sirupsen/logrus"
)

var parserLogger *logrus.Entry

func init() {
	parserLogger = logrus.WithField("component", "codegen-apply-parser")
}

// parseServiceFile parses an existing service file and extracts information
func parseServiceFile(filePath string) (*ServiceFileInfo, error) {
	parserLogger.WithField("filePath", filePath).Info("Starting service file parsing")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		parserLogger.WithField("filePath", filePath).Debug("Service file does not exist")
		return nil, nil // File doesn't exist, not an error
	}

	// Parse the file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		parserLogger.WithFields(logrus.Fields{
			"filePath": filePath,
			"error":    err,
		}).Error("Failed to parse service file")
		return nil, err
	}

	parserLogger.WithField("packageName", file.Name.Name).Debug("Successfully parsed service file AST")

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
				
				parserLogger.WithFields(logrus.Fields{
					"methodName": methodName,
					"filePath":   filePath,
				}).Debug("Found existing method")
				
				// Extract business logic
				businessLogic := extractBusinessLogic(funcDecl.Body, fset)
				if len(businessLogic) > 0 {
					info.BusinessCode[methodName] = businessLogic
					parserLogger.WithFields(logrus.Fields{
						"methodName": methodName,
						"linesCount": len(businessLogic),
					}).Debug("Extracted business logic")
				}
			}
		}
	}

	parserLogger.WithFields(logrus.Fields{
		"filePath":    filePath,
		"methodCount": methodCount,
		"methods":     getMethodNames(info.Methods),
	}).Info("Service file parsing completed")

	return info, nil
}

// getMethodNames extracts method names from a map of function declarations
func getMethodNames(methods map[string]*ast.FuncDecl) []string {
	names := make([]string, 0, len(methods))
	for name := range methods {
		names = append(names, name)
	}
	return names
}

// scanServiceFiles scans the service directory for existing service files
func scanServiceFiles(serviceDir string, models []*gen.ModelInfo) (map[string]*ServiceFileInfo, error) {
	parserLogger.WithFields(logrus.Fields{
		"serviceDir":  serviceDir,
		"modelCount":  len(models),
		"models":      extractModelNames(models),
	}).Info("Starting service files scan")

	serviceFiles := make(map[string]*ServiceFileInfo)
	scannedCount := 0
	foundCount := 0

	for _, model := range models {
		scannedCount++
		modelLogger := parserLogger.WithFields(logrus.Fields{
			"modelName":   model.ModelName,
			"packageName": model.PackageName,
		})

		// Generate service file path based on model info
		// Extract the subdirectory from model file path to maintain directory structure
		modelRelPath, err := filepath.Rel(filepath.Join(serviceDir, "..", "model"), model.ModelFileDir)
		if err != nil {
			modelLogger.WithFields(logrus.Fields{
				"modelFileDir": model.ModelFileDir,
				"error":        err,
			}).Error("Failed to get relative path for model directory")
			return nil, err
		}
		
		servicePkgName := gen.ModelPkg2ServicePkg(model.PackageName)
		serviceFilePath := filepath.Join(serviceDir, modelRelPath, strings.ToLower(model.ModelName)+".go")
		
		modelLogger.WithFields(logrus.Fields{
			"serviceFilePath": serviceFilePath,
			"servicePkgName":  servicePkgName,
		}).Debug("Scanning for service file")

		// Update model with service file path
		model.ServiceFilePath = serviceFilePath
		
		// Parse the service file if it exists
		serviceInfo, err := parseServiceFile(serviceFilePath)
		if err != nil {
			modelLogger.WithFields(logrus.Fields{
				"serviceFilePath": serviceFilePath,
				"error":           err,
			}).Error("Failed to parse service file")
			return nil, err
		}
		
		if serviceInfo != nil {
			// Use model name as key
			serviceFiles[model.ModelName] = serviceInfo
			foundCount++
			modelLogger.WithFields(logrus.Fields{
				"serviceFilePath": serviceFilePath,
				"methodCount":     len(serviceInfo.Methods),
			}).Info("Found existing service file")
		} else {
			modelLogger.WithField("serviceFilePath", serviceFilePath).Debug("No existing service file found")
		}
	}

	parserLogger.WithFields(logrus.Fields{
		"serviceDir":         serviceDir,
		"scannedModels":      scannedCount,
		"foundServiceFiles":  foundCount,
		"serviceFiles":       getServiceFileKeys(serviceFiles),
	}).Info("Service files scan completed")

	return serviceFiles, nil
}

// getServiceFileKeys extracts keys from service files map for logging
func getServiceFileKeys(serviceFiles map[string]*ServiceFileInfo) []string {
	keys := make([]string, 0, len(serviceFiles))
	for key := range serviceFiles {
		keys = append(keys, key)
	}
	return keys
}

// needsRegeneration checks if a service file needs to be regenerated
func needsRegeneration(model *gen.ModelInfo, serviceInfo *ServiceFileInfo) bool {
	regenLogger := parserLogger.WithFields(logrus.Fields{
		"modelName":       model.ModelName,
		"serviceFilePath": model.ServiceFilePath,
	})

	regenLogger.Debug("Checking if service file needs regeneration")

	if serviceInfo == nil {
		regenLogger.Info("Service file needs generation - file doesn't exist")
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

	regenLogger.WithFields(logrus.Fields{
		"requiredMethodCount": len(requiredMethods),
		"existingMethodCount": len(serviceInfo.Methods),
	}).Debug("Checking required methods")

	missingMethods := []string{}
	for _, methodName := range requiredMethods {
		if _, exists := serviceInfo.Methods[methodName]; !exists {
			missingMethods = append(missingMethods, methodName)
		}
	}

	if len(missingMethods) > 0 {
		regenLogger.WithFields(logrus.Fields{
			"missingMethodCount": len(missingMethods),
			"missingMethods":     missingMethods,
		}).Info("Service file needs regeneration - missing methods")
		return true
	}

	regenLogger.WithField("existingMethods", getMethodNames(serviceInfo.Methods)).Info("Service file is up to date - all required methods exist")
	return false // All methods exist, no regeneration needed
}