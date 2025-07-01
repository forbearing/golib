// internal/codegen/types.go
package codegen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ServiceTemplateData 包含生成 service 代码所需的所有数据
type ServiceTemplateData struct {
	PackageName   string // service 包名
	ModelPackage  string // model 包的导入路径
	ModelName     string // Model 结构体名称 (如: User)
	ServiceName   string // service 结构体名称 (如: user)
	ModelVariable string // model 变量名 (如: users)
	FrameworkPath string // 框架路径
}

// ModelInfo 包含解析的 model 信息
type ModelInfo struct {
	Name        string
	PackageName string
	PackagePath string
}

// Config 代码生成配置
type Config struct {
	ServicePackage string // service 包名，默认 "service"
	FrameworkPath  string // 框架路径，默认 "github.com/forbearing/golib"
	OutputDir      string // 输出目录
	ModulePath     string // 从 go.mod 读取的模块路径
}

// LoadConfig 从当前目录加载配置
func LoadConfig() (*Config, error) {
	modulePath, err := getModulePath()
	if err != nil {
		return nil, fmt.Errorf("failed to read module path: %w", err)
	}

	return &Config{
		ServicePackage: "service",
		FrameworkPath:  "github.com/forbearing/golib",
		OutputDir:      "./service",
		ModulePath:     modulePath,
	}, nil
}

// getModulePath 从 go.mod 文件读取模块路径
func getModulePath() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("go.mod file not found in current directory. Please run this command from project root")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading go.mod: %w", err)
	}

	return "", fmt.Errorf("module declaration not found in go.mod")
}
