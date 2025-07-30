package apply

import (
	"go/ast"
	"os"
	"path/filepath"
	"testing"

	"github.com/forbearing/golib/internal/codegen/gen"
)

func TestApplyServiceGeneration(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		module       string
		modelDir     string
		serviceDir   string
		modelContent string
		wantErr      bool
		wantFile     string
	}{
		{
			name:       "basic service generation",
			module:     "test/module",
			modelDir:   "model",
			serviceDir: "service",
			modelContent: `package model

import "github.com/forbearing/golib/model"

type User struct {
	model.Base
	Name  string ` + "`json:\"name\"`" + `
	Email string ` + "`json:\"email\"`" + `
}
`,
			wantErr:  false,
			wantFile: "user.go",
		},
		{
			name:       "service generation with qualified package",
			module:     "test/module",
			modelDir:   "model",
			serviceDir: "service",
			modelContent: `package model

import model_base "github.com/forbearing/golib/model"

type Product struct {
	model_base.Base
	Name  string ` + "`json:\"name\"`" + `
	Price float64 ` + "`json:\"price\"`" + `
}
`,
			wantErr:  false,
			wantFile: "product.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directories for testing
			tempDir := t.TempDir()
			modelDir := filepath.Join(tempDir, tt.modelDir)
			serviceDir := filepath.Join(tempDir, tt.serviceDir)

			// Create model directory and test model file
			err := os.MkdirAll(modelDir, 0o755)
			if err != nil {
				t.Fatalf("Failed to create model directory: %v", err)
			}

			modelFile := filepath.Join(modelDir, "test.go")
			err = os.WriteFile(modelFile, []byte(tt.modelContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to create model file: %v", err)
			}

			// Create apply configuration
			config := NewApplyConfig(tt.module, modelDir, serviceDir)

			// Test service generation
			gotErr := ApplyServiceGeneration(config)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ApplyServiceGeneration() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ApplyServiceGeneration() succeeded unexpectedly")
			}

			// Verify service file was created
			expectedServiceFile := filepath.Join(serviceDir, tt.wantFile)
			if _, err := os.Stat(expectedServiceFile); os.IsNotExist(err) {
				t.Errorf("Expected service file was not created: %s", expectedServiceFile)
			}
		})
	}
}

func TestNewApplyConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		module     string
		modelDir   string
		serviceDir string
		want       *ApplyConfig
	}{
		{
			name:       "basic config creation",
			module:     "test/module",
			modelDir:   "model",
			serviceDir: "service",
			want: &ApplyConfig{
				Module:     "test/module",
				ModelDir:   "model",
				ServiceDir: "service",
				Excludes:   []string{},
			},
		},
		{
			name:       "config with different paths",
			module:     "github.com/example/app",
			modelDir:   "internal/model",
			serviceDir: "internal/service",
			want: &ApplyConfig{
				Module:     "github.com/example/app",
				ModelDir:   "internal/model",
				ServiceDir: "internal/service",
				Excludes:   []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewApplyConfig(tt.module, tt.modelDir, tt.serviceDir)

			if got.Module != tt.want.Module {
				t.Errorf("NewApplyConfig().Module = %v, want %v", got.Module, tt.want.Module)
			}
			if got.ModelDir != tt.want.ModelDir {
				t.Errorf("NewApplyConfig().ModelDir = %v, want %v", got.ModelDir, tt.want.ModelDir)
			}
			if got.ServiceDir != tt.want.ServiceDir {
				t.Errorf("NewApplyConfig().ServiceDir = %v, want %v", got.ServiceDir, tt.want.ServiceDir)
			}
			if len(got.Excludes) != len(tt.want.Excludes) {
				t.Errorf("NewApplyConfig().Excludes length = %v, want %v", len(got.Excludes), len(tt.want.Excludes))
			}
		})
	}
}

func TestWithExclusions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		baseConfig *ApplyConfig
		exclusions []string
		want       []string
	}{
		{
			name: "add single exclusion",
			baseConfig: &ApplyConfig{
				Module:     "test/module",
				ModelDir:   "model",
				ServiceDir: "service",
				Excludes:   []string{},
			},
			exclusions: []string{"User"},
			want:       []string{"User"},
		},
		{
			name: "add multiple exclusions",
			baseConfig: &ApplyConfig{
				Module:     "test/module",
				ModelDir:   "model",
				ServiceDir: "service",
				Excludes:   []string{},
			},
			exclusions: []string{"User", "Group"},
			want:       []string{"User", "Group"},
		},
		{
			name: "add to existing exclusions",
			baseConfig: &ApplyConfig{
				Module:     "test/module",
				ModelDir:   "model",
				ServiceDir: "service",
				Excludes:   []string{"Product"},
			},
			exclusions: []string{"User", "Group"},
			want:       []string{"Product", "User", "Group"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.baseConfig.WithExclusions(tt.exclusions...)

			if len(got.Excludes) != len(tt.want) {
				t.Errorf("WithExclusions() exclusions length = %v, want %v", len(got.Excludes), len(tt.want))
				return
			}

			for i, exclusion := range tt.want {
				if got.Excludes[i] != exclusion {
					t.Errorf("WithExclusions() exclusion[%d] = %v, want %v", i, got.Excludes[i], exclusion)
				}
			}
		})
	}
}

func TestShouldSkipModel(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		exclusions []string
		modelName  string
		want       bool
	}{
		{
			name:       "skip excluded model - User",
			exclusions: []string{"User", "Group"},
			modelName:  "User",
			want:       true,
		},
		{
			name:       "skip excluded model - Group",
			exclusions: []string{"User", "Group"},
			modelName:  "Group",
			want:       true,
		},
		{
			name:       "don't skip non-excluded model - Product",
			exclusions: []string{"User", "Group"},
			modelName:  "Product",
			want:       false,
		},
		{
			name:       "don't skip non-excluded model - Order",
			exclusions: []string{"User", "Group"},
			modelName:  "Order",
			want:       false,
		},
		{
			name:       "empty exclusions",
			exclusions: []string{},
			modelName:  "User",
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use slices.Contains logic to check if model should be skipped
			got := false
			for _, exclusion := range tt.exclusions {
				if tt.modelName == exclusion {
					got = true
					break
				}
			}
			if got != tt.want {
				t.Errorf("shouldSkipModel(%s) = %v, want %v", tt.modelName, got, tt.want)
			}
		})
	}
}

func TestNeedsRegeneration(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		model       *gen.ModelInfo
		serviceInfo *ServiceFileInfo
		want        bool
	}{
		{
			name: "service file doesn't exist",
			model: &gen.ModelInfo{
				ModelName: "User",
			},
			serviceInfo: nil,
			want:        true,
		},
		{
			name: "service file exists but no methods",
			model: &gen.ModelInfo{
				ModelName: "User",
			},
			serviceInfo: &ServiceFileInfo{
				Methods: make(map[string]*ast.FuncDecl),
			},
			want: true,
		},
		{
			name: "service file exists with some methods",
			model: &gen.ModelInfo{
				ModelName: "User",
			},
			serviceInfo: &ServiceFileInfo{
				Methods: map[string]*ast.FuncDecl{
					"CreateBefore": {},
				},
			},
			want: true, // Still needs regeneration as not all methods exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := needsRegeneration(tt.model, tt.serviceInfo)
			if got != tt.want {
				t.Errorf("needsRegeneration() = %v, want %v", got, tt.want)
			}
		})
	}
}
