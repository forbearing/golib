package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check architecture dependencies in generated code",
	Long: `Check architecture dependencies in generated code:
1. Service code should not call other service code
2. DAO code should not call service code
3. Model code should not call service code`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRun()
	},
}

func checkRun() {
	logSection("Architecture Dependency Check")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	var violations []string

	// Check service files
	serviceViolations := checkServiceDependencies(filepath.Join(cwd, serviceDir))
	violations = append(violations, serviceViolations...)

	// Check dao files
	daoViolations := checkDAODependencies(filepath.Join(cwd, "dao"))
	violations = append(violations, daoViolations...)

	// Check model files
	modelViolations := checkModelDependencies(filepath.Join(cwd, modelDir))
	violations = append(violations, modelViolations...)

	if len(violations) > 0 {
		fmt.Println("\n❌ Architecture violations found:")
		for _, violation := range violations {
			fmt.Printf("  - %s\n", violation)
		}
		os.Exit(1)
	} else {
		fmt.Println("\n✅ No architecture violations found")
	}
}

// checkServiceDependencies checks if service code calls other service code
func checkServiceDependencies(serviceDir string) []string {
	var violations []string

	if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
		return violations
	}

	err := filepath.Walk(serviceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "_test.go") {
			return nil
		}

		// Skip service.go registration file
		if strings.HasSuffix(path, "service.go") {
			return nil
		}

		fileViolations := checkFileForServiceImports(path, "service")
		violations = append(violations, fileViolations...)

		return nil
	})
	if err != nil {
		fmt.Printf("Error walking service directory: %v\n", err)
	}

	return violations
}

// checkDAODependencies checks if DAO code calls service code
func checkDAODependencies(daoDir string) []string {
	var violations []string

	if _, err := os.Stat(daoDir); os.IsNotExist(err) {
		return violations
	}

	err := filepath.Walk(daoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "_test.go") {
			return nil
		}

		fileViolations := checkFileForServiceImports(path, "dao")
		violations = append(violations, fileViolations...)

		return nil
	})
	if err != nil {
		fmt.Printf("Error walking dao directory: %v\n", err)
	}

	return violations
}

// checkModelDependencies checks if model code calls service code
func checkModelDependencies(modelDir string) []string {
	var violations []string

	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		return violations
	}

	err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") || strings.Contains(path, "_test.go") {
			return nil
		}

		// Skip model.go registration file
		if strings.HasSuffix(path, "model.go") {
			return nil
		}

		fileViolations := checkFileForServiceImports(path, "model")
		violations = append(violations, fileViolations...)

		return nil
	})
	if err != nil {
		fmt.Printf("Error walking model directory: %v\n", err)
	}

	return violations
}

// checkFileForServiceImports checks a single file for service imports
func checkFileForServiceImports(filePath, layerType string) []string {
	var violations []string

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Error parsing file %s: %v\n", filePath, err)
		return violations
	}

	// Check imports
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)

		// Check for service imports
		if containsServiceImport(importPath, layerType) {
			violation := fmt.Sprintf("%s file '%s' imports service code: %s",
				strings.Title(layerType), filePath, importPath)
			violations = append(violations, violation)
		}
	}

	return violations
}

// containsServiceImport checks if an import path contains service code
func containsServiceImport(importPath, layerType string) bool {
	// Split import path by '/'
	parts := strings.Split(importPath, "/")

	for i, part := range parts {
		if part == "service" {
			// For service layer, check if it's importing other service packages
			if layerType == "service" {
				// Allow importing the base service package
				if i == len(parts)-1 && part == "service" {
					return false
				}
				// Check if it's importing another service implementation
				if i < len(parts)-1 {
					return true
				}
			}
			// For dao and model layers, any service import is forbidden
			if layerType == "dao" || layerType == "model" {
				// Allow importing the base service package for interfaces
				if i == len(parts)-1 && part == "service" {
					return false
				}
				// Forbid importing service implementations
				return true
			}
		}
	}

	return false
}
