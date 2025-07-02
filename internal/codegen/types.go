// internal/codegen/types.go
package codegen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ServiceTemplateData 包含生成 service 代码所需的所有数据
type ServiceTemplateData struct {
	PackageName   string // service 包名 (如: service_keycloak)
	ModelPackage  string // model 包的导入路径
	ModelName     string // Model 结构体名称 (如: User)
	ServiceName   string // service 结构体名称 (如: user)
	ModelVariable string // model 变量名 (如: users)
	FrameworkPath string // 框架路径
}

// ModelInfo 包含解析的 model 信息
type ModelInfo struct {
	Name         string // 结构体名称
	PackageName  string // 原包名
	PackagePath  string // 完整包路径
	RelativePath string // 相对于 model 根目录的路径
	FilePath     string // 文件路径
}

// Config 代码生成配置
type Config struct {
	ModelDir       string // model 根目录
	ServiceDir     string // service 根目录
	FrameworkPath  string // 框架路径
	ModulePath     string // 从 go.mod 读取的模块路径
	ForceOverwrite bool   // 强制覆盖已存在的文件
}

// LoadConfig 从当前目录加载配置
func LoadConfig() (*Config, error) {
	modulePath, err := getModulePath()
	if err != nil {
		return nil, fmt.Errorf("failed to read module path: %w", err)
	}

	return &Config{
		ModelDir:       "./model",
		ServiceDir:     "./service",
		FrameworkPath:  "github.com/forbearing/golib",
		ModulePath:     modulePath,
		ForceOverwrite: false, // 默认不覆盖
	}, nil
}

// GetServicePackageName 根据相对路径生成 service 包名
func (c *Config) GetServicePackageName(relativePath string) string {
	if relativePath == "" || relativePath == "." {
		return "service"
	}

	// 将路径分隔符替换为下划线
	packageSuffix := strings.ReplaceAll(relativePath, string(filepath.Separator), "_")
	packageSuffix = strings.ReplaceAll(packageSuffix, "/", "_")

	return "service_" + packageSuffix
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
