package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/forbearing/golib/dsl"
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
4. Model directories and files must be singular
5. Model struct json tags should use snake_case naming`,
	Run: func(cmd *cobra.Command, args []string) {
		checkRun()
	},
}

// performArchitectureCheckForCheck performs architecture dependency checks for check command
func performArchitectureCheckForCheck(cwd string) []string {
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

	return violations
}

func checkRun() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		logError(fmt.Sprintf("getting current directory: %v", err))
		os.Exit(1)
	}

	var totalViolations int

	// Architecture Dependency Check
	logSection("Architecture Dependency Check")
	archViolations := performArchitectureCheckForCheck(cwd)
	if len(archViolations) > 0 {
		for _, violation := range archViolations {
			fmt.Printf("  %s %s\n", red("→"), violation)
		}
		totalViolations += len(archViolations)
	} else {
		fmt.Printf("  %s No architecture violations found\n", green("✔"))
	}

	// Model Singular Naming Check
	logSection("Model Singular Naming Check")
	singularViolations := checkModelSingularNaming(filepath.Join(cwd, modelDir))
	if len(singularViolations) > 0 {
		for _, violation := range singularViolations {
			fmt.Printf("  %s %s\n", red("→"), violation)
		}
		totalViolations += len(singularViolations)
	} else {
		fmt.Printf("  %s No singular naming violations found\n", green("✔"))
	}

	// Model JSON Tag Naming Check
	logSection("Model JSON Tag Naming Check")
	jsonTagViolations := checkModelJSONTagNaming(filepath.Join(cwd, modelDir))
	if len(jsonTagViolations) > 0 {
		for _, violation := range jsonTagViolations {
			fmt.Printf("  %s %s\n", red("→"), violation)
		}
		totalViolations += len(jsonTagViolations)
	} else {
		fmt.Printf("  %s No JSON tag naming violations found\n", green("✔"))
	}

	// Summary
	logSection("Summary")
	if totalViolations > 0 {
		fmt.Printf("  %s %d violations found\n", red("✘"), totalViolations)
		os.Exit(1)
	} else {
		fmt.Printf("  %s All checks passed\n", green("✔"))
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

// checkModelJSONTagNaming checks if model struct json tags use camelCase naming
func checkModelJSONTagNaming(modelDir string) []string {
	var violations []string

	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		return violations
	}

	err := filepath.Walk(modelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip generated files
		if strings.HasSuffix(path, "model.go") {
			return nil
		}

		fileViolations := checkFileJSONTagNaming(path)
		violations = append(violations, fileViolations...)

		return nil
	})
	if err != nil {
		logError(fmt.Sprintf("walking model directory: %v", err))
	}

	return violations
}

// checkFileJSONTagNaming checks json tag naming in a single file
func checkFileJSONTagNaming(filePath string) []string {
	var violations []string

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return violations
	}

	// Find all model structs in this file
	modelBaseNames := dsl.FindAllModelBase(node)
	modelEmptyNames := dsl.FindAllModelEmpty(node)
	allModelNames := append(modelBaseNames, modelEmptyNames...)

	// If no model structs found, skip this file
	if len(allModelNames) == 0 {
		return violations
	}

	// Get relative path for cleaner output
	cwd, _ := os.Getwd()
	relPath, _ := filepath.Rel(cwd, filePath)

	// Check only model structs
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// Check if this struct is a model
			isModel := false
			for _, modelName := range allModelNames {
				if typeSpec.Name.Name == modelName {
					isModel = true
					break
				}
			}
			if !isModel {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok || structType.Fields == nil {
				continue
			}

			// Check JSON tags in this model struct
			for _, field := range structType.Fields.List {
				if field.Tag != nil {
					tagValue := strings.Trim(field.Tag.Value, "`")
					if jsonTag := extractJSONTag(tagValue); jsonTag != "" {
						if !isSnakeCase(jsonTag) {
							fieldName := ""
							if len(field.Names) > 0 {
								fieldName = field.Names[0].Name
							}
							violations = append(violations, fmt.Sprintf(
								"%s: field '%s' json tag '%s' should be '%s'",
								relPath, fieldName, jsonTag, toSnakeCase(jsonTag)))
						}
					}
				}
			}
		}
	}

	return violations
}

// extractJSONTag extracts the json tag value from struct tag
func extractJSONTag(tag string) string {
	re := regexp.MustCompile(`json:"([^"]+)"`)
	matches := re.FindStringSubmatch(tag)
	if len(matches) > 1 {
		// Remove options like omitempty
		parts := strings.Split(matches[1], ",")
		return parts[0]
	}
	return ""
}

// isSnakeCase checks if a string is in snake_case format
func isSnakeCase(s string) bool {
	if s == "" {
		return true
	}

	// Skip special cases like "-" or single characters
	if s == "-" || len(s) == 1 {
		return true
	}

	// Check if it contains hyphens (kebab-case) or uppercase letters
	if strings.Contains(s, "-") {
		return false
	}

	// Check for uppercase letters (not snake_case)
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return false
		}
	}

	return true
}

// toSnakeCase converts camelCase or kebab-case to snake_case
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	// Replace hyphens with underscores
	s = strings.ReplaceAll(s, "-", "_")

	// Convert camelCase to snake_case
	var result strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(r - 'A' + 'a')
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}
