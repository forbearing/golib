// internal/codegen/walker/walker.go
package walker

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/internal/codegen"
)

// Walker 目录遍历器
type Walker struct {
	config  *codegen.Config
	fileSet *token.FileSet
	parser  *ModelParser
}

// ModelParser model 解析器
type ModelParser struct {
	fileSet    *token.FileSet
	modulePath string
}

// NewWalker 创建新的遍历器
func NewWalker(config *codegen.Config) *Walker {
	fileSet := token.NewFileSet()
	return &Walker{
		config:  config,
		fileSet: fileSet,
		parser: &ModelParser{
			fileSet:    fileSet,
			modulePath: config.ModulePath,
		},
	}
}

// WalkModels 遍历 model 目录并解析所有 model
func (w *Walker) WalkModels() ([]*codegen.ModelInfo, error) {
	var allModels []*codegen.ModelInfo

	err := filepath.Walk(w.config.ModelDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 只处理 .go 文件
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// 跳过需要忽略的文件
		if w.shouldSkipFile(info.Name()) {
			fmt.Printf("Skipping file: %s\n", path)
			return nil
		}

		// 解析文件中的 models
		models, err := w.parser.ParseModelFile(path, w.config.ModelDir)
		if err != nil {
			return fmt.Errorf("failed to parse file %s: %w", path, err)
		}

		allModels = append(allModels, models...)
		return nil
	})

	return allModels, err
}

// shouldSkipFile 检查是否应该跳过文件
func (w *Walker) shouldSkipFile(filename string) bool {
	// 跳过测试文件
	if strings.HasSuffix(filename, "_test.go") {
		return true
	}

	// 跳过以 _ 开头的文件
	if strings.HasPrefix(filename, "_") {
		return true
	}

	// 跳过 doc.go
	if filename == "doc.go" {
		return true
	}

	return false
}

// ParseModelFile 解析单个 model 文件
func (mp *ModelParser) ParseModelFile(filePath, modelRoot string) ([]*codegen.ModelInfo, error) {
	file, err := parser.ParseFile(mp.fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	// 计算相对路径
	relativeDir, err := filepath.Rel(modelRoot, filepath.Dir(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	// 解析导入
	imports := mp.parseImports(file)

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
			if mp.hasBaseEmbedded(structType, imports) {
				model := &codegen.ModelInfo{
					Name:         typeSpec.Name.Name,
					PackageName:  file.Name.Name, // 这是实际的包名
					PackagePath:  mp.getPackagePath(filePath, relativeDir),
					RelativePath: relativeDir,
					FilePath:     filePath,
				}
				models = append(models, model)
			}
		}
	}

	return models, nil
}

// parseImports 解析导入声明
func (mp *ModelParser) parseImports(file *ast.File) map[string]string {
	imports := make(map[string]string)

	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		var alias string

		if imp.Name != nil {
			alias = imp.Name.Name
		} else {
			parts := strings.Split(importPath, "/")
			alias = parts[len(parts)-1]
		}

		imports[alias] = importPath
	}

	return imports
}

// hasBaseEmbedded 检查结构体是否嵌入了 Base 接口
func (mp *ModelParser) hasBaseEmbedded(structType *ast.StructType, imports map[string]string) bool {
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			if mp.isBaseType(field.Type, imports) {
				return true
			}
		}
	}
	return false
}

// isBaseType 检查类型是否是 Base 接口
func (mp *ModelParser) isBaseType(expr ast.Expr, imports map[string]string) bool {
	switch t := expr.(type) {
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			packageAlias := ident.Name
			typeName := t.Sel.Name

			if typeName != "Base" {
				return false
			}

			if importPath, exists := imports[packageAlias]; exists {
				return mp.isGolibModelPackage(importPath)
			}
		}
	case *ast.Ident:
		return t.Name == "Base"
	}
	return false
}

// isGolibModelPackage 检查导入路径是否是 golib 的 model 包
func (mp *ModelParser) isGolibModelPackage(importPath string) bool {
	return strings.Contains(importPath, "golib/model") ||
		importPath == "github.com/forbearing/golib/model"
}

// getPackagePath 获取包路径
func (mp *ModelParser) getPackagePath(filePath, relativePath string) string {
	if relativePath == "." {
		return mp.modulePath + "/model"
	}

	// 转换为 Go 包路径格式
	pkgPath := strings.ReplaceAll(relativePath, string(filepath.Separator), "/")
	return mp.modulePath + "/model/" + pkgPath
}
