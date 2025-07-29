package apply

import (
	"go/ast"
	"go/token"
)

// ServiceFileInfo contains information about an existing service file
type ServiceFileInfo struct {
	FilePath     string                    // Path to the service file
	ModelName    string                    // Name of the model (e.g., "User")
	PackageName  string                    // Package name (e.g., "service_system")
	Methods      map[string]*ast.FuncDecl  // Existing hook methods
	FileSet      *token.FileSet            // File set for parsing
	File         *ast.File                 // Parsed AST file
	BusinessCode map[string][]string       // Business logic code for each method
}

// ApplyConfig contains configuration for the apply operation
type ApplyConfig struct {
	ModelDir   string   // Model directory path
	ServiceDir string   // Service directory path
	Module     string   // Module path
	Excludes   []string // Files to exclude
}

// BusinessLogicSection represents a section of business logic code
type BusinessLogicSection struct {
	MethodName string   // Name of the method containing this logic
	StartLine  int      // Start line of business logic
	EndLine    int      // End line of business logic
	Code       []string // Lines of business logic code
}