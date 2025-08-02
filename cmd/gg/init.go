package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const mainContent = `package main

import (
	"github.com/forbearing/golib/bootstrap"
	. "github.com/forbearing/golib/util"
)

func main() {
	RunOrDie(bootstrap.Bootstrap)
	RunOrDie(bootstrap.Run)
}`

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

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init project",
	Long:  "init project",
	Args:  cobra.ExactArgs(1),
	// Steps to execute:
	// 1. Initialize a Go project
	// 2. Create three directories: model, service, router, each containing a .gitkeep empty file
	// 3. Create main.go file with content: mainContent
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		initProject(projectName)
	},
}

// initProject initializes a new Go project with the specified name
func initProject(projectName string) {
	// Extract project name from module path (e.g., "github.com/user/project" -> "project")
	projectDir := filepath.Base(projectName)
	
	// Step 1: Create project directory
	fmt.Printf("Creating project directory: %s\n", projectDir)
	checkErr(os.MkdirAll(projectDir, 0o755))
	
	// Step 2: Change to project directory
	fmt.Printf("Entering directory: %s\n", projectDir)
	checkErr(os.Chdir(projectDir))
	
	// Step 3: Initialize Go module
	fmt.Printf("Initializing Go module: %s\n", projectName)
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	checkErr(cmd.Run())

	// Step 4: Create directories with .gitkeep files
	dirs := []string{"model", "service", "router"}
	for _, dir := range dirs {
		fmt.Printf("Creating directory: %s\n", dir)
		checkErr(os.MkdirAll(dir, 0o755))

		// Create .gitkeep file in each directory
		gitkeepPath := filepath.Join(dir, ".gitkeep")
		fmt.Printf("Creating file: %s\n", gitkeepPath)
		file, err := os.Create(gitkeepPath)
		checkErr(err)
		file.Close()
	}

	// Step 5: Create main.go file
	fmt.Println("Creating main.go")
	checkErr(os.WriteFile("main.go", []byte(mainContent), 0o644))

	// Step 6: Create .gitignore file
	fmt.Println("Creating .gitignore")
	checkErr(os.WriteFile(".gitignore", []byte(gitignoreContent), 0o644))

	// Step 7: Run go mod tidy
	fmt.Println("Running go mod tidy...")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	checkErr(cmd.Run())

	// Step 8: Initialize git repository
	fmt.Println("Initializing git repository...")
	cmd = exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	checkErr(cmd.Run())

	fmt.Println("Project initialization completed successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  cd", projectDir)
	fmt.Println("  git add .")
	fmt.Println("  git commit -m \"Initial commit\"")
}
