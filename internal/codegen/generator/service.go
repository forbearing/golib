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
	config   *codegen.Config
}

// NewServiceGenerator 创建新的 service 生成器
func NewServiceGenerator(config *codegen.Config) (*ServiceGenerator, error) {
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
		config:   config,
	}, nil
}

// GenerateFromModel 从 model 信息生成 service 代码
func (g *ServiceGenerator) GenerateFromModel(modelInfo *codegen.ModelInfo) ([]byte, error) {
	// 生成包名
	packageName := g.config.GetServicePackageName(modelInfo.RelativePath)

	// 获取 model 包的导入信息
	modelPackageAlias, needsAlias := modelInfo.GetModelPackageInfo()

	data := &codegen.ServiceTemplateData{
		PackageName:       packageName,
		ModelPackage:      modelInfo.PackagePath,
		ModelPackageName:  modelInfo.PackageName,
		ModelPackageAlias: modelPackageAlias,
		ModelName:         modelInfo.Name,
		ServiceName:       strings.ToLower(modelInfo.Name),
		ModelVariable:     strings.ToLower(modelInfo.Name) + "s",
		FrameworkPath:     g.config.FrameworkPath,
		NeedsAlias:        needsAlias,
	}

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

// GenerateToFile 生成代码并写入文件
func (g *ServiceGenerator) GenerateToFile(modelInfo *codegen.ModelInfo) error {
	// 计算输出路径
	outputDir := filepath.Join(g.config.ServiceDir, modelInfo.RelativePath)
	filename := filepath.Join(outputDir, strings.ToLower(modelInfo.Name)+".go")

	// 检查文件是否已存在且不为空（仅在非强制覆盖模式下）
	if !g.config.ForceOverwrite && g.fileExistsAndNotEmpty(filename) {
		fmt.Printf("Skipping %s: file %s already exists and is not empty (use --force to overwrite)\n",
			modelInfo.Name, filename)
		return nil
	}

	// 如果是强制覆盖模式且文件存在，显示覆盖信息
	if g.config.ForceOverwrite && g.fileExists(filename) {
		fmt.Printf("Overwriting existing file: %s\n", filename)
	}

	// 生成代码
	code, err := g.GenerateFromModel(modelInfo)
	if err != nil {
		return fmt.Errorf("failed to generate code for model %s: %w", modelInfo.Name, err)
	}

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
	}

	// 写入文件
	if err := os.WriteFile(filename, code, 0o644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	fmt.Printf("Generated service code for %s -> %s\n", modelInfo.Name, filename)
	return nil
}

// GenerateAll 批量生成所有 models 的 service 代码
func (g *ServiceGenerator) GenerateAll(models []*codegen.ModelInfo) error {
	var errors []string
	generatedCount := 0
	skippedCount := 0

	for _, model := range models {
		// 检查是否会跳过
		outputDir := filepath.Join(g.config.ServiceDir, model.RelativePath)
		filename := filepath.Join(outputDir, strings.ToLower(model.Name)+".go")

		willSkip := !g.config.ForceOverwrite && g.fileExistsAndNotEmpty(filename)

		if err := g.GenerateToFile(model); err != nil {
			errors = append(errors, fmt.Sprintf("model %s: %v", model.Name, err))
		} else {
			if willSkip {
				skippedCount++
			} else {
				generatedCount++
			}
		}
	}

	// 显示统计信息
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Generated: %d files\n", generatedCount)
	if skippedCount > 0 {
		fmt.Printf("  Skipped: %d files (already exist)\n", skippedCount)
	}
	if len(errors) > 0 {
		fmt.Printf("  Errors: %d files\n", len(errors))
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to generate some service files:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// fileExists 检查文件是否存在
func (g *ServiceGenerator) fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// fileExistsAndNotEmpty 检查文件是否存在且不为空
func (g *ServiceGenerator) fileExistsAndNotEmpty(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return info.Size() > 0
}
