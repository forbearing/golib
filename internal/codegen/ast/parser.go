// internal/codegen/ast/parser.go
package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/internal/codegen"
)

// Parser AST 解析器
type Parser struct {
	fileSet    *token.FileSet
	modulePath string
}

// NewParser 创建新的解析器
func NewParser(modulePath string) *Parser {
	return &Parser{
		fileSet:    token.NewFileSet(),
		modulePath: modulePath,
	}
}

// ParseModelFile 解析 model 文件
func (p *Parser) ParseModelFile(filename string) ([]*codegen.ModelInfo, error) {
	file, err := parser.ParseFile(p.fileSet, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filename, err)
	}

	// 首先解析导入，找到 model 包的别名
	imports := p.parseImports(file)

	var models []*codegen.ModelInfo

	// 遍历文件中的所有声明
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		// 遍历类型声明
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			// 检查是否是结构体
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 检查是否嵌入了 Base 接口
			if p.hasBaseEmbedded(structType, imports) {
				model := &codegen.ModelInfo{
					Name:        typeSpec.Name.Name,
					PackageName: file.Name.Name,
					PackagePath: p.getPackagePath(filename, file.Name.Name),
				}
				models = append(models, model)
			}
		}
	}

	return models, nil
}

// ParseModelFileDebug 带调试信息的解析方法
func (p *Parser) ParseModelFileDebug(filename string) ([]*codegen.ModelInfo, error) {
	fmt.Printf("=== Parsing file: %s ===\n", filename)
	fmt.Printf("Module path: %s\n", p.modulePath)

	models, err := p.ParseModelFile(filename)
	if err != nil {
		return nil, err
	}

	fmt.Printf("=== Total models found: %d ===\n", len(models))
	for _, model := range models {
		fmt.Printf("  Model: %s, Package: %s, Path: %s\n",
			model.Name, model.PackageName, model.PackagePath)
	}

	return models, nil
}

// parseImports 解析导入声明，返回别名到导入路径的映射
func (p *Parser) parseImports(file *ast.File) map[string]string {
	imports := make(map[string]string)

	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		var alias string

		if imp.Name != nil {
			// 显式别名
			alias = imp.Name.Name
		} else {
			// 默认别名是包名
			parts := strings.Split(importPath, "/")
			alias = parts[len(parts)-1]
		}

		imports[alias] = importPath
	}

	return imports
}

// hasBaseEmbedded 检查结构体是否嵌入了 Base 接口
func (p *Parser) hasBaseEmbedded(structType *ast.StructType, imports map[string]string) bool {
	for _, field := range structType.Fields.List {
		// 只检查嵌入字段（匿名字段）
		if len(field.Names) == 0 {
			if p.isBaseType(field.Type, imports) {
				return true
			}
		}
	}
	return false
}

// isBaseType 检查类型是否是 Base 接口
func (p *Parser) isBaseType(expr ast.Expr, imports map[string]string) bool {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		// 处理 package.Type 形式，如 model.Base
		if ident, ok := t.X.(*ast.Ident); ok {
			packageAlias := ident.Name
			typeName := t.Sel.Name

			// 检查是否是 Base 类型
			if typeName != "Base" {
				return false
			}

			// 检查包是否是 golib 的 model 包
			if importPath, exists := imports[packageAlias]; exists {
				return p.isGolibModelPackage(importPath)
			}
		}
	case *ast.Ident:
		// 处理同包内的类型，如 Base
		typeName := t.Name
		return typeName == "Base"
	}
	return false
}

// isGolibModelPackage 检查导入路径是否是 golib 的 model 包
func (p *Parser) isGolibModelPackage(importPath string) bool {
	// 检查是否是 golib 框架的 model 包
	return strings.Contains(importPath, "golib/model") ||
		importPath == "github.com/forbearing/golib/model"
}

// getPackagePath 获取包路径
func (p *Parser) getPackagePath(filename, packageName string) string {
	dir := filepath.Dir(filename)

	// 转换为绝对路径
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return p.modulePath + "/" + packageName
	}

	// 获取当前工作目录
	wd, err := filepath.Abs(".")
	if err != nil {
		return p.modulePath + "/" + packageName
	}

	// 计算相对路径
	relPath, err := filepath.Rel(wd, absDir)
	if err != nil {
		return p.modulePath + "/" + packageName
	}

	// 构建完整的包路径
	if relPath == "." {
		return p.modulePath
	}

	// 转换为 Go 包路径格式
	pkgPath := strings.ReplaceAll(relPath, string(filepath.Separator), "/")
	return p.modulePath + "/" + pkgPath
}
