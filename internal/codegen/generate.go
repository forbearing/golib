package codegen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stoewer/go-strcase"
)

// ModelInfo 存储模型信息
//
// 例如:
// {ModulePath:"github.com/forbearing/golib", PackageName:"model", ModelName:"User", ModelVarName:"u", ModelFileDir:"/tmp/model", ServiceFilePath:""},
// {ModulePath:"github.com/forbearing/golib", PackageName:"model", ModelName:"Group", ModelVarName:"g", ModelFileDir:"/tmp/model", ServiceFilePath:""},
// {ModulePath:"github.com/forbearing/golib", PackageName:"model_auth", ModelName:"User", ModelVarName:"u", ModelFileDir:"/tmp/model", ServiceFilePath:""},
// {ModulePath:"github.com/forbearing/golib", PackageName:"model_auth", ModelName:"Group", ModelVarName:"g", ModelFileDir:"/tmp/model", ServiceFilePath:""},

type ModelInfo struct {
	ModulePath      string // 从 go.mod 解析的模块路径
	PackageName     string // model 包名, 例如: model_authz, model_log
	ModelName       string // model 名, 例如: User, Group
	ModelVarName    string // 小写的模型变量名, 例如: u, g
	ModelFileDir    string // model 文件所在目录的的相对路径, 例如: github.com/forbearing/golib/model
	ServiceFilePath string // service 文件的相对路径, 例如: github.com/forbearing/golib/service
}

// getModulePath 解析 go.mod 获取模块路径
func getModulePath() (string, error) {
	// 如果存在 go 命令直接通过 go list -m 命令获取模块路径
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

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
// import "github.com/forbearing/golib/model" 则为 "model"
// import model_auth "github.com/forbearing/golib/model", 则为 model_auth
func findModelPackageName(file *ast.File) string {
	return file.Name.Name
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if strings.HasSuffix(path, "github.com/forbearing/golib/model") {
			// pretty.Println("-----", imp.Name)
			if imp.Name != nil {
				return imp.Name.Name // 使用重命名的包名
			}
			return "model" // 默认包名
		}
	}
	return ""
}

// isModelBase 检查字段是否是 model.Base
func isModelBase(file *ast.File, field *ast.Field, modelPkgName string) bool {
	if field.Names != nil { // 不是匿名字段
		return false
	}

	getAliasName := func(file *ast.File) string {
		for _, imp := range file.Imports {
			path := strings.Trim(imp.Path.Value, `"`)
			if strings.HasSuffix(path, "github.com/forbearing/golib/model") {
				if imp.Name != nil {
					return imp.Name.Name // 使用重命名的包名
				}
				return "model" // 默认包名
			}
		}
		return ""
	}
	aliasName := getAliasName(file)

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name == aliasName && t.Sel.Name == "Base"
		}
	case *ast.Ident:
		// 处理同包的情况
		return t.Name == "Base"
	}

	return false
}

// findModels 查找 model 文件中的所有结构体
func findModels(modulePath string, filename string) ([]*ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	modelPkgName := findModelPackageName(node)
	if len(modelPkgName) == 0 {
		return nil, fmt.Errorf("file %s has no model package", filename)
	}

	var models []*ModelInfo
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl == nil {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec == nil {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok || structType == nil || structType.Fields == nil {
				continue
			}
			hasModelBase := false
			for _, field := range structType.Fields.List {
				if isModelBase(node, field, modelPkgName) {
					hasModelBase = true
					break
				}
			}
			if !hasModelBase || typeSpec.Name == nil {
				continue
			}
			modelName := typeSpec.Name.Name
			if len(modelName) == 0 {
				continue
			}
			models = append(models, &ModelInfo{
				PackageName:  modelPkgName,
				ModelName:    modelName,
				ModelVarName: strings.ToLower(modelName[:1]),
				ModulePath:   modulePath,
				ModelFileDir: filepath.Dir(filename),
			})

		}
	}

	return models, nil
}

func isValidModelPackage(packageName string) bool {
	return strings.HasPrefix(packageName, "model")
}

// modelPkg2ServicePkg 根据 model name 转换成 service name.
func modelPkg2ServicePkg(pkgName string) string {
	return strings.Replace(pkgName, "model", "service", 1)
}

// generateServiceMethod1 使用 AST 生成 CreateBefore 等方法
func generateServiceMethod1(info *ModelInfo, methodName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return service_method_1(info.ModelVarName, info.ModelName, methodName, info.PackageName,
		assign_with_service_context(info.ModelVarName),
		expr_log_info(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		empty_line(),
		returns("nil"),
	)
}

// 使用 AST 生成 ListBefore, ListAfter 方法.
func generateServiceMethod2(info *ModelInfo, methodName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return service_method_2(info.ModelVarName, info.ModelName, methodName, info.PackageName,
		assign_with_service_context(info.ModelVarName),
		expr_log_info(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		empty_line(),
		returns("nil"),
	)
}

// 使用 AST 生成 CreateManyBefore, CreateManyAfter, DeleteManyBefore, DeleteManyAfter
// UpdateManyBefore, UpdateManyAfter, PatchManyBefore, PatchManyAfter.
func generateServiceMethod3(info *ModelInfo, methodName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return service_method_3(info.ModelVarName, info.ModelName, methodName, info.PackageName,
		assign_with_service_context(info.ModelVarName),
		expr_log_info(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		empty_line(),
		returns("nil"),
	)
}

func generateServiceFile(info *ModelInfo) *ast.File {
	if !isValidModelPackage(info.PackageName) {
		return nil
	}

	decls := []ast.Decl{
		imports(info.ModulePath, info.ModelFileDir, info.PackageName),
		inits(info.ModelName),
		types(info.ModelName, info.PackageName),
	}

	for _, method := range methods {
		if strings.HasPrefix(method, "List") {
			decls = append(decls, generateServiceMethod2(info, method))
		} else if strings.Contains(method, "Many") {
			decls = append(decls, generateServiceMethod3(info, method))
		} else {
			decls = append(decls, generateServiceMethod1(info, method))
		}
	}

	return &ast.File{
		Name:  ast.NewIdent(modelPkg2ServicePkg(info.PackageName)),
		Decls: decls,
	}
}

func formatNode(node ast.Node) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}

	// TODO: 使用 gofumpt
	formated, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}
	return string(formated), nil
}

func methodAddComments(code string, modelName string) string {
	for _, method := range methods {
		str := strings.ReplaceAll(strcase.SnakeCase(method), "_", " ")
		// 在 log.Info 之后添加注释
		searchStr := fmt.Sprintf(`log.Info("%s %s")`, strings.ToLower(modelName), str)
		replaceStr := fmt.Sprintf(`log.Info("%s %s")
	// =============================
	// Add your business logic here.
	// =============================
`, strings.ToLower(modelName), str)

		code = strings.ReplaceAll(code, searchStr, replaceStr)
	}

	return code
}
