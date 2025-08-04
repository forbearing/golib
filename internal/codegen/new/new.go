package new

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

const mainContent = `package main

import (
	"%s/router"

	"github.com/forbearing/golib/bootstrap"
	. "github.com/forbearing/golib/util"
)

func main() {
	RunOrDie(bootstrap.Bootstrap)
	RunOrDie(router.Init)
	RunOrDie(bootstrap.Run)
}
`

const routerContent = `package router

func Init() error {
	return nil
}
`

const gitignoreContent = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Log files
*.log

# Temporary files
tmp/
temp/

# Build output
dist/
build/`

// Run initializes a new Go project with the specified project name.
func Run(projectName string) error {
	// Extract project name from module path (e.g., "github.com/user/project" -> "project")
	projectDir := filepath.Base(projectName)

	// Step 1: Create project directory
	fmt.Printf("Creating project directory: %s\n", projectDir)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return err
	}

	// Step 2: Change to project directory
	fmt.Printf("Entering directory: %s\n", projectDir)
	if err := os.Chdir(projectDir); err != nil {
		return err
	}

	// Step 3: Initialize Go module
	fmt.Printf("Initializing Go module: %s\n", projectName)
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Step 4: Create directories with .gitkeep files
	dirs := []string{"model", "service", "router"}
	for _, dir := range dirs {
		fmt.Printf("Creating directory: %s\n", dir)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		// Create .gitkeep file in each directory
		gitkeepPath := filepath.Join(dir, ".gitkeep")
		fmt.Printf("Creating file: %s\n", gitkeepPath)
		var file *os.File
		var err error
		if file, err = os.Create(gitkeepPath); err != nil {
			return err
		}
		file.Close()
	}
	// Step 5: Create main.go file
	fmt.Println("Creating main.go")
	if err := os.WriteFile("main.go", fmt.Appendf(nil, mainContent, projectName), 0o644); err != nil {
		return err
	}

	// Step 6: Create router/router.go
	fmt.Println("Creating router/router.go")
	if err := os.WriteFile("router/router.go", []byte(routerContent), 0o644); err != nil {
		return err
	}

	// Step 7: Create .gitignore file
	fmt.Println("Creating .gitignore")
	if err := os.WriteFile(".gitignore", []byte(gitignoreContent), 0o644); err != nil {
		return err
	}

	// Step 8: Run go mod tidy
	fmt.Println("Running go mod tidy...")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return err
	}

	// Step 9: Initialize git repository
	fmt.Println("Initializing git repository...")
	cmd = exec.Command("git", "init")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("Project initialization completed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  cd", projectDir)
	fmt.Println("  git add .")
	fmt.Println("  git commit -m \"Initial commit\"")

	return nil
}
