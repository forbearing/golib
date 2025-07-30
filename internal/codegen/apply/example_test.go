package apply_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/forbearing/golib/internal/codegen/apply"
)

// ExampleApplyServiceGeneration demonstrates how to use the apply package
// to generate service files from model definitions.
func ExampleApplyServiceGeneration() {
	// Create temporary directories for the example
	tempDir, err := os.MkdirTemp("", "apply_example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	modelDir := filepath.Join(tempDir, "model")
	serviceDir := filepath.Join(tempDir, "service")

	// Create model directory and a sample model file
	if err := os.MkdirAll(modelDir, 0o755); err != nil {
		log.Fatal(err)
	}

	// Create go.mod file in temp directory for module detection
	goModContent := `module example.com/myapp

go 1.21
`
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0o644); err != nil {
		log.Fatal(err)
	}

	modelContent := `package model

import (
	"time"
	"github.com/forbearing/golib/model"
)

type User struct {
	model.Base
	Name      string    ` + "`" + `json:"name"` + "`" + `
	Email     string    ` + "`" + `json:"email"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at"` + "`" + `
}
`

	if err := os.WriteFile(filepath.Join(modelDir, "user.go"), []byte(modelContent), 0o644); err != nil {
		log.Fatal(err)
	}

	// Create configuration for service generation
	config := apply.NewApplyConfig(
		"example.com/myapp", // module path
		modelDir,            // model directory
		serviceDir,          // service directory
	)

	// Apply service generation
	if err := apply.ApplyServiceGeneration(config); err != nil {
		fmt.Printf("Service generation failed: %v", err)
		return
	}

	// Check if service file was created
	serviceFile := filepath.Join(serviceDir, "user.go")
	if _, err := os.Stat(serviceFile); err == nil {
		fmt.Println("Service file generated successfully")
	} else {
		fmt.Printf("Service file not found: %v", err)
	}

	// Output: Service file generated successfully
}

