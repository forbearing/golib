package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check architecture dependencies in generated code",
	Long: `Check architecture dependencies in generated code:
1. Service code should not call other service code
2. DAO code should not call service code
3. Model code should not call service code
4. Model directories and files must be singular`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRun()
	},
}

func checkRun() {
	logSection("Architecture Dependency Check")

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logError(fmt.Sprintf("getting current directory: %v", err))
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

	// Check model singular naming
	singularViolations := checkModelSingularNaming(filepath.Join(cwd, modelDir))
	violations = append(violations, singularViolations...)

	if len(violations) > 0 {
		fmt.Printf("\n%s Architecture violations found:\n", red("✘"))
		for _, violation := range violations {
			fmt.Printf("  %s %s\n", red("→"), violation)
		}
		os.Exit(1)
	} else {
		fmt.Printf("\n%s No architecture violations found\n", green("✔"))
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
		logError(fmt.Sprintf("walking service directory: %v", err))
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
		logError(fmt.Sprintf("walking dao directory: %v", err))
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
		logError(fmt.Sprintf("walking model directory: %v", err))
	}

	return violations
}

// checkFileForServiceImports checks a single file for service imports
func checkFileForServiceImports(filePath, layerType string) []string {
	var violations []string

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		// Treat parse errors as violations to prevent code generation
		violation := fmt.Sprintf("%s file '%s' has parse error: %v",
			strings.Title(layerType), filePath, err)
		violations = append(violations, violation)
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

// checkModelSingularNaming checks if model directories and files use singular names
func checkModelSingularNaming(modelDir string) []string {
	var violations []string

	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		return violations
	}

	client := pluralize.NewClient()

	err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from model directory
		relPath, err := filepath.Rel(modelDir, path)
		if err != nil {
			return err
		}

		// Skip the root model directory itself
		if relPath == "." {
			return nil
		}

		if info.IsDir() {
			// Check directory name
			dirName := info.Name()
			if client.IsPlural(dirName) {
				violation := fmt.Sprintf("Model directory '%s' should be singular (suggested: %s)",
					path, client.Singular(dirName))
				violations = append(violations, violation)
			}
		} else if strings.HasSuffix(path, ".go") && !strings.Contains(path, "_test.go") {
			// Skip model.go registration file
			if strings.HasSuffix(path, "model.go") {
				return nil
			}

			// Check Go file name (without .go extension)
			fileName := strings.TrimSuffix(info.Name(), ".go")
			if client.IsPlural(fileName) {
				violation := fmt.Sprintf("Model file '%s' should be singular (suggested: %s.go)",
					path, client.Singular(fileName))
				violations = append(violations, violation)
			}
		}

		return nil
	})
	if err != nil {
		logError(fmt.Sprintf("walking model directory: %v", err))
	}

	return violations
}

// containsServiceImport checks if an import path contains service code
func containsServiceImport(importPath, layerType string) bool {
	// Split import path by '/'
	parts := strings.Split(importPath, "/")

	for _, part := range parts {
		if part == "service" {
			// For service layer, check if it's importing other service packages
			if layerType == "service" {
				// Allow importing the base golib service package only
				if strings.Contains(importPath, "github.com/forbearing/golib/service") {
					return false
				}
				// Forbid importing any other service implementations
				return true
			}
			// For dao and model layers, any service import is forbidden except golib service
			if layerType == "dao" || layerType == "model" {
				// Allow importing the base golib service package for interfaces
				if strings.Contains(importPath, "github.com/forbearing/golib/service") {
					return false
				}
				// Forbid importing service implementations
				return true
			}
		}
	}

	return false
}
