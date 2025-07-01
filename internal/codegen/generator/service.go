// internal/codegen/generator/service.go
package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/forbearing/golib/internal/codegen"
)

// ServiceGenerator 生成 service 层代码
type ServiceGenerator struct {
	template *template.Template
}

// NewServiceGenerator 创建新的 service 生成器
func NewServiceGenerator() (*ServiceGenerator, error) {
	funcMap := template.FuncMap{
		"firstChar": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return strings.ToLower(string(s[0]))
		},
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
	}

	tmpl, err := template.New("service").Funcs(funcMap).Parse(ServiceTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse service template: %w", err)
	}

	return &ServiceGenerator{
		template: tmpl,
	}, nil
}

// Generate 生成 service 代码
func (g *ServiceGenerator) Generate(data *codegen.ServiceTemplateData) ([]byte, error) {
	var buf bytes.Buffer

	if err := g.template.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	// 格式化生成的代码
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated code: %w", err)
	}

	return formatted, nil
}

// GenerateFromModel 从 model 信息生成 service 代码
func (g *ServiceGenerator) GenerateFromModel(modelInfo *codegen.ModelInfo, config *codegen.Config) ([]byte, error) {
	data := &codegen.ServiceTemplateData{
		PackageName:   config.ServicePackage,
		ModelPackage:  modelInfo.PackagePath,
		ModelName:     modelInfo.Name,
		ServiceName:   strings.ToLower(modelInfo.Name),
		ModelVariable: strings.ToLower(modelInfo.Name) + "s", // users, products, etc.
		FrameworkPath: config.FrameworkPath,
	}

	return g.Generate(data)
}

// GenerateToFile 生成代码并写入文件
func (g *ServiceGenerator) GenerateToFile(modelInfo *codegen.ModelInfo, config *codegen.Config) error {
	// 生成文件名（小写）
	filename := filepath.Join(config.OutputDir, strings.ToLower(modelInfo.Name)+".go")

	// 检查文件是否已存在且不为空
	if g.fileExistsAndNotEmpty(filename) {
		fmt.Printf("Skipping %s: file %s already exists and is not empty\n", modelInfo.Name, filename)
		return nil
	}

	// 生成代码
	code, err := g.GenerateFromModel(modelInfo, config)
	if err != nil {
		return fmt.Errorf("failed to generate code for model %s: %w", modelInfo.Name, err)
	}

	// 创建输出目录
	if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(filename, code, 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	fmt.Printf("Generated service code for %s -> %s\n", modelInfo.Name, filename)
	return nil
}

// fileExistsAndNotEmpty 检查文件是否存在且不为空
func (g *ServiceGenerator) fileExistsAndNotEmpty(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		// 文件不存在
		return false
	}

	// 文件存在且大小大于0
	return info.Size() > 0
}
