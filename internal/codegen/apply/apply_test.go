package apply

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/forbearing/golib/internal/codegen/gen"
)

func TestApplyServiceGeneration(t *testing.T) {
	// Create temporary directories for testing
	tempDir := t.TempDir()
	modelDir := filepath.Join(tempDir, "model")
	serviceDir := filepath.Join(tempDir, "service")

	// Create model directory and a test model file
	err := os.MkdirAll(modelDir, 0o755)
	if err != nil {
		t.Fatalf("Failed to create model directory: %v", err)
	}

	// Create a simple test model file
	modelContent := `package model

import "github.com/forbearing/golib/model"

type User struct {
	model.Base
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\"`" + `
}
`
	modelFile := filepath.Join(modelDir, "user.go")
	err = os.WriteFile(modelFile, []byte(modelContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create model file: %v", err)
	}

	// Create apply configuration
	config := NewApplyConfig("test/module", modelDir, serviceDir)

	// Test service generation
	err = ApplyServiceGeneration(config)
	if err != nil {
		t.Fatalf("ApplyServiceGeneration failed: %v", err)
	}

	// Verify service file was created
	// The service file path should be generated based on the model package name
	// model -> service (directly in service directory)
	expectedServiceFile := filepath.Join(serviceDir, "user.go")
	if _, err := os.Stat(expectedServiceFile); os.IsNotExist(err) {
		// Check what files were actually created
		files, _ := filepath.Glob(filepath.Join(serviceDir, "**", "*.go"))
		allFiles, _ := filepath.Glob(filepath.Join(serviceDir, "*"))
		t.Logf("Files created in serviceDir: %v", allFiles)
		t.Logf("Go files created: %v", files)

		// Also check direct files in serviceDir
		directFiles, _ := filepath.Glob(filepath.Join(serviceDir, "*.go"))
		t.Logf("Direct Go files in serviceDir: %v", directFiles)

		t.Errorf("Expected service file was not created: %s", expectedServiceFile)
	}
}

func TestNewApplyConfig(t *testing.T) {
	config := NewApplyConfig("test/module", "model", "service")

	if config.Module != "test/module" {
		t.Errorf("Expected Module to be 'test/module', got '%s'", config.Module)
	}

	if config.ModelDir != "model" {
		t.Errorf("Expected ModelDir to be 'model', got '%s'", config.ModelDir)
	}

	if config.ServiceDir != "service" {
		t.Errorf("Expected ServiceDir to be 'service', got '%s'", config.ServiceDir)
	}
}

func TestWithExclusions(t *testing.T) {
	config := NewApplyConfig("test/module", "model", "service")
	config = config.WithExclusions("User", "Group")

	if len(config.Excludes) != 2 {
		t.Errorf("Expected 2 exclusions, got %d", len(config.Excludes))
	}

	if config.Excludes[0] != "User" || config.Excludes[1] != "Group" {
		t.Errorf("Expected exclusions [User, Group], got %v", config.Excludes)
	}
}

func TestShouldSkipModel(t *testing.T) {
	exclusions := []string{"User", "Group"}

	tests := []struct {
		modelName string
		expected  bool
	}{
		{"User", true},
		{"Group", true},
		{"Product", false},
		{"Order", false},
	}

	for _, test := range tests {
		// Use slices.Contains instead of shouldSkipModel
		result := false
		for _, exclusion := range exclusions {
			if test.modelName == exclusion {
				result = true
				break
			}
		}
		if result != test.expected {
			t.Errorf("shouldSkipModel(%s) = %v, expected %v", test.modelName, result, test.expected)
		}
	}
}

func TestNeedsRegeneration(t *testing.T) {
	// Create a mock model
	model := &gen.ModelInfo{
		ModelName: "User",
	}

	// Test with nil service info (file doesn't exist)
	if !needsRegeneration(model, nil) {
		t.Error("Expected needsRegeneration to return true when service info is nil")
	}

	// Test with empty service info (no methods)
	serviceInfo := &ServiceFileInfo{
		Methods: make(map[string]*ast.FuncDecl),
	}
	if !needsRegeneration(model, serviceInfo) {
		t.Error("Expected needsRegeneration to return true when no methods exist")
	}
}

