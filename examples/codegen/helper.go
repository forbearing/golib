package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ModelInfo 存储模型信息
// 例如:
//	{PackageName:"model", ModelName:"User", ModelVarName:"u", ModulePath:"codegen", ModelFilePath:"model"},
//	{PackageName:"model", ModelName:"Group", ModelVarName:"g", ModulePath:"codegen", ModelFilePath:"model"},

type ModelInfo struct {
	PackageName   string
	ModelName     string
	ModelVarName  string // 小写的模型变量名
	ModulePath    string // 从 go.mod 解析的模块路径
	ModelFilePath string // model 文件的相对路径
}

// getModulePath 解析 go.mod 获取模块路径
func getModulePath() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var moduleName string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				moduleName = parts[1]
			}
		}
	}

	return moduleName, scanner.Err()
}

// findModelPackageName 查找包中导入的 model 包的实际名称
func findModelPackageName(file *ast.File) string {
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if strings.HasSuffix(path, "/model") {
			if imp.Name != nil {
				return imp.Name.Name // 使用重命名的包名
			}
			return "model" // 默认包名
		}
	}
	return ""
}

// isModelBase 检查字段是否是 model.Base
func isModelBase(field *ast.Field, modelPkgName string) bool {
	if field.Names != nil { // 不是匿名字段
		return false
	}

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name == modelPkgName && t.Sel.Name == "Base"
		}
	case *ast.Ident:
		// 处理同包的情况
		return t.Name == "Base"
	}

	return false
}

// findModels 查找 model 文件中的所有结构体
func findModels(filename string, modulePath string) ([]ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// if node.Name == nil {
	// 	return nil, fmt.Errorf("file %s has no package name", filename)
	// }
	// packageName := node.Name.Name
	modelPkgName := findModelPackageName(node)
	// fmt.Println(packageName, modelPkgName)
	if len(modelPkgName) == 0 {
		return nil, fmt.Errorf("file %s has no model package", filename)
	}

	var models []ModelInfo
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl == nil {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			if typeSpec == nil {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if structType.Fields == nil {
				continue
			}
			hasModelBase := false
			for _, field := range structType.Fields.List {
				if isModelBase(field, modelPkgName) {
					hasModelBase = true
					break
				}
			}
			if !hasModelBase {
				continue
			}
			if typeSpec.Name == nil {
				continue
			}
			modelName := typeSpec.Name.Name
			if len(modelName) == 0 {
				continue
			}
			models = append(models, ModelInfo{
				PackageName:   modelPkgName,
				ModelName:     modelName,
				ModelVarName:  strings.ToLower(modelName[:1]),
				ModulePath:    modulePath,
				ModelFilePath: filepath.Dir(filename),
			})

		}
	}

	return models, nil
}

// 使用 AST 生成 CreateBefore 方法
func generateServiceMethod(model ModelInfo, methodName string) *ast.FuncDecl {
	return service_method_1(model.ModelVarName, model.ModelName, methodName,
		assign_with_service_context(model.ModelVarName),
		expr_log_info(fmt.Sprintf(`"%s create before"`, strings.ToLower(model.ModelName))),
		empty_line(),
		returns("nil"),
	)
}
