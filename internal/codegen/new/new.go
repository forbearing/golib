package new

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/forbearing/golib/config"
)

// Run initializes a new Go project with the specified project name.
func Run(projectName string) error {
	// Extract project name from module path (e.g., "github.com/user/project" -> "project")
	projectDir := filepath.Base(projectName)

	// Create project directory
	fmt.Printf("Creating project directory: %s\n", projectDir)
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		return err
	}

	// Change to project directory
	fmt.Printf("Entering directory: %s\n", projectDir)
	if err := os.Chdir(projectDir); err != nil {
		return err
	}

	// Initialize Go module
	fmt.Printf("Initializing Go module: %s\n", projectName)
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create directories with .gitkeep files
	dirs := []string{"configx", "cronjobx", "model", "service", "router"}
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

	// Create configx/configx.go
	fmt.Println("Creating configx/configx.go")
	if err := os.WriteFile("configx/configx.go", []byte(configxContent), 0o644); err != nil {
		return err
	}

	// Create cronjobx/cronjobx.go
	fmt.Println("Creating cronjobx/cronjobx.go")
	if err := os.WriteFile("cronjobx/cronjobx.go", []byte(cronjobxContent), 0o644); err != nil {
		return err
	}

	fmt.Println("Creating model/model.go")
	if err := os.WriteFile("model/model.go", []byte(modelContent), 0o644); err != nil {
		return err
	}

	// Create service/service.go
	fmt.Println("Creating service/service.go")
	if err := os.WriteFile("service/service.go", []byte(serviceContent), 0o644); err != nil {
		return err
	}

	// Create router/router.go
	fmt.Println("Creating router/router.go")
	if err := os.WriteFile("router/router.go", []byte(routerContent), 0o644); err != nil {
		return err
	}

	// Create main.go file
	fmt.Println("Creating main.go")
	if err := os.WriteFile("main.go", fmt.Appendf(nil, mainContent, projectName, projectName, projectName, projectName, projectName), 0o644); err != nil {
		return err
	}

	// Create .gitignore file
	fmt.Println("Creating .gitignore")
	if err := os.WriteFile(".gitignore", []byte(gitignoreContent), 0o644); err != nil {
		return err
	}

	// Create template config.ini
	fmt.Println("Creating config.ini")
	if err := createTeplConfig(); err != nil {
		return err
	}

	// Run go mod tidy
	fmt.Println("Running go mod tidy...")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		return err
	}

	// Initialize git repository
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

func createTeplConfig() error {
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()

	null, err := os.Open(os.DevNull)
	if err != nil {
		return err
	}
	os.Stdout = null

	if err := config.Init(); err != nil {
		return err
	}
	defer config.Clean()

	if err := config.Save("config.ini"); err != nil {
		return err
	}
	if err := os.Rename("config.ini", "config.ini.example"); err != nil {
		return err
	}
	return nil
}
