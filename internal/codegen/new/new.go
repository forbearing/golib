package new

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/forbearing/golib/config"
)

var fileContentMap = map[string]string{
	"configx/configx.go":         configxContent,
	"cronjobx/cronjobx.go":       cronjobxContent,
	"middlewarex/middlewarex.go": middlewarexContent,
	"model/model.go":             modelContent,
	"service/service.go":         serviceContent,
	"router/router.go":           routerContent,
}

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

	for file, content := range fileContentMap {
		if err := createFile(file, content); err != nil {
			return err
		}
	}

	// Create main.go file
	if err := createFile("main.go", fmt.Sprintf(mainContent, projectName, projectName, projectName, projectName, projectName, projectName)); err != nil {
		return err
	}
	// Create .gitignore file
	if err := createFile(".gitignore", gitignoreContent); err != nil {
		return err
	}
	// Create template config.ini
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

func EnsureFileExists() {
	for file, content := range fileContentMap {
		if _, err := os.Stat(file); err != nil && errors.Is(err, os.ErrNotExist) {
			createFile(file, content)
		}
	}
}

func createFile(path string, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	fmt.Println("Creating", path)
	return os.WriteFile(path, []byte(content), 0o644)
}

func createTeplConfig() error {
	fmt.Println("Creating config.ini")
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
