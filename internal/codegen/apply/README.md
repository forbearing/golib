# Apply Package

The `apply` package provides functionality for generating service files from model definitions while preserving existing business logic.

## Features

- **Model Discovery**: Automatically discovers model definitions from specified directories
- **Service Generation**: Generates complete service files with CRUD operations
- **Business Logic Preservation**: Preserves existing business logic when regenerating service files
- **Configurable Exclusions**: Allows excluding specific models from generation
- **Path Management**: Handles proper file paths and package naming conventions

## Core Components

### Types

- **`ApplyConfig`**: Configuration for service generation operations
- **`ServiceFileInfo`**: Metadata about service files including business logic sections
- **`BusinessLogicSection`**: Represents preserved business logic code segments

### Main Functions

- **`ApplyServiceGeneration(config *ApplyConfig)`**: Main entry point for service generation
- **`NewApplyConfig(module, modelDir, serviceDir string)`**: Creates new configuration
- **`WithExclusions(excludes []string)`**: Adds model exclusions to configuration

## Usage Example

```go
package main

import (
    "log"
    "github.com/forbearing/golib/internal/codegen/apply"
)

func main() {
    // Create configuration
    config := apply.NewApplyConfig(
        "example.com/myapp",  // module path
        "./model",            // model directory
        "./service",          // service directory
    ).WithExclusions([]string{"internal_model"}) // exclude specific models

    // Generate services
    if err := apply.ApplyServiceGeneration(config); err != nil {
        log.Fatal(err)
    }
}
```

## How It Works

1. **Model Discovery**: Scans the model directory for Go files containing struct definitions
2. **Service File Scanning**: Checks existing service files and extracts business logic
3. **Code Generation**: Generates new service files using the `gen` package
4. **Business Logic Preservation**: Merges preserved business logic into generated files
5. **File Writing**: Writes the final service files to the specified directory

## Generated Service Structure

Each generated service file includes:

- Package declaration and imports
- Service type definition implementing the Service interface
- CRUD methods (Create, Update, Delete, Get, List)
- Initialization function for service registration
- Preserved business logic from existing files

## Configuration Options

- **Module Path**: Go module path for import generation
- **Model Directory**: Directory containing model definitions
- **Service Directory**: Target directory for generated service files
- **Exclusions**: List of model names to exclude from generation

## Business Logic Preservation

The package automatically preserves business logic by:

1. Parsing existing service files using Go AST
2. Identifying business logic sections between generated method signatures
3. Extracting and storing these sections with line number information
4. Reinserting preserved logic into newly generated files

This ensures that custom business logic is not lost during regeneration cycles.