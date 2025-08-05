package gen

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

	"github.com/forbearing/golib/dsl"
	"github.com/stoewer/go-strcase"
	fumpt "mvdan.cc/gofumpt/format"
)

// ModelInfo 存储模型信息
//
// 例如:
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model", ModelName:"User", ModelVarName:"u", ModelFileDir:"/tmp/model", ServiceFilePath:""},
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model", ModelName:"Group", ModelVarName:"g", ModelFileDir:"/tmp/model", ServiceFilePath:""},
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model_auth", ModelName:"User", ModelVarName:"u", ModelFileDir:"/tmp/model", ServiceFilePath:""},
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model_auth", ModelName:"Group", ModelVarName:"g", ModelFileDir:"/tmp/model", ServiceFilePath:""},

type ModelInfo struct {
	// module 相关字段
	ModulePath string // 从 go.mod 解析的模块路径

	// model 相关字段
	ModelPkgName  string // model 包名, 例如: model, model_authz, model_log
	ModelName     string // model 名, 例如: User, Group
	ModelVarName  string // 小写的模型变量名, 例如: u, g
	ModelFileDir  string // model 文件所在目录的的相对路径, 例如: github.com/forbearing/golib/model
	ModelFilePath string // model 文件的相对路径, 例如: github.com/forbearing/golib/model/user.go

	// Service 相关字段
	ServiceFilePath string // service 文件的相对路径, 例如: github.com/forbearing/golib/service

	// 自定义请求和相应相关字段
	Design *dsl.Design
}

// GetModulePath 解析 go.mod 获取模块路径
func GetModulePath() (string, error) {
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

// FindModelPackageName 查找包中导入的 model 包的实际名称
// import "github.com/forbearing/golib/model" 则为 "model"
// import model_auth "github.com/forbearing/golib/model", 则为 model_auth
func FindModelPackageName(file *ast.File) string {
	return file.Name.Name
}

// // IsModelBase 检查字段是否是 model.Base
//
//	func IsModelBase(file *ast.File, field *ast.Field, modelPkgName string) bool {
//		if field.Names != nil { // 不是匿名字段
//			return false
//		}
//
//		getAliasName := func(file *ast.File) string {
//			for _, imp := range file.Imports {
//				path := strings.Trim(imp.Path.Value, `"`)
//				if strings.HasSuffix(path, "github.com/forbearing/golib/model") {
//					if imp.Name != nil {
//						return imp.Name.Name // 使用重命名的包名
//					}
//					return "model" // 默认包名
//				}
//			}
//			return ""
//		}
//		aliasName := getAliasName(file)
//
//		switch t := field.Type.(type) {
//		case *ast.SelectorExpr:
//			if ident, ok := t.X.(*ast.Ident); ok {
//				return ident.Name == aliasName && t.Sel.Name == "Base"
//			}
//		case *ast.Ident:
//			// 处理同包的情况
//			return t.Name == "Base"
//		}
//
//		return false
//	}
func IsModelBase(file *ast.File, field *ast.Field) bool {
	// Not anonymouse field.
	if len(field.Names) != 0 {
		return false
	}

	aliasName := "model"
	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}
		if imp.Path.Value == `"github.com/forbearing/golib/model"` {
			if imp.Name != nil {
				aliasName = imp.Name.Name
			}
			break
		}
	}

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name == aliasName && t.Sel.Name == "Base"
		}
	case *ast.Ident:
		return t.Name == "Base"
	}

	return false
}

// FindModels 查找 model 文件中的所有结构体
// TODO: 支持自定义 Request、Response 不和 model 同一个包位置
func FindModels(module string, filename string) ([]*ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	modelPkgName := FindModelPackageName(node)
	if len(modelPkgName) == 0 {
		return nil, fmt.Errorf("file %s has no model package", filename)
	}
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	designs := dsl.Parse(f)

	var models []*ModelInfo
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl == nil || genDecl.Tok != token.TYPE {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec == nil || typeSpec.Type == nil {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok || structType == nil || structType.Fields == nil {
				continue
			}
			hasModelBase := false
			for _, field := range structType.Fields.List {
				if IsModelBase(node, field) {
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
				ModelFileDir:  filepath.Dir(filename),
				ModelFilePath: filename,
				ModelPkgName:  modelPkgName,
				ModelName:     modelName,
				ModelVarName:  strings.ToLower(modelName[:1]),
				ModulePath:    module,
				Design:        designs[modelName],
			})

		}
	}

	return models, nil
}

func IsValidModelPackage(packageName string) bool {
	return strings.HasPrefix(packageName, "model")
}

// ModelPkg2ServicePkg 根据 model name 转换成 service name.
func ModelPkg2ServicePkg(pkgName string) string {
	if pkgName == "model" {
		return "service"
	}
	// 对于 model_xxx 格式，替换为 service_xxx
	if strings.HasPrefix(pkgName, "model_") {
		return strings.Replace(pkgName, "model_", "service_", 1)
	}
	return strings.Replace(pkgName, "model", "service", 1)
}

// GenerateServiceMethod1 使用 AST 生成 CreateBefore,CreateAfter,UpdateBefore,UpdateAfter,
// DeleteBefore,DeleteAfter,GetBefore,GetAfter,UpdatePartialBefore,UpdatePartialAfter 方法.
func GenerateServiceMethod1(info *ModelInfo, methodName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return ServiceMethod1(info.ModelVarName, info.ModelName, methodName, info.ModelPkgName,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns("nil"),
	)
}

// GenerateServiceMethod2 使用 AST 生成 ListBefore, ListAfter 方法.
func GenerateServiceMethod2(info *ModelInfo, methodName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return ServiceMethod2(info.ModelVarName, info.ModelName, methodName, info.ModelPkgName,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns("nil"),
	)
}

// GenerateServiceMethod3 使用 AST 生成 CreateManyBefore, CreateManyAfter,
// DeleteManyBefore, DeleteManyAfter, UpdateManyBefore, UpdateManyAfter, PatchManyBefore, PatchManyAfter.
func GenerateServiceMethod3(info *ModelInfo, methodName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return ServiceMethod3(info.ModelVarName, info.ModelName, methodName, info.ModelPkgName,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns("nil"),
	)
}

func GenerateService(info *ModelInfo) *ast.File {
	if !IsValidModelPackage(info.ModelPkgName) {
		return nil
	}

	decls := []ast.Decl{
		Imports(info.ModulePath, info.ModelFileDir, info.ModelPkgName),
		Inits(info.ModelName),
		Types(info.ModelName, info.ModelPkgName),
	}

	for _, method := range Methods {
		if strings.HasPrefix(method, "List") {
			decls = append(decls, GenerateServiceMethod2(info, method))
		} else if strings.Contains(method, "Many") {
			decls = append(decls, GenerateServiceMethod3(info, method))
		} else {
			decls = append(decls, GenerateServiceMethod1(info, method))
		}
	}

	return &ast.File{
		Name:  ast.NewIdent(ModelPkg2ServicePkg(info.ModelPkgName)),
		Decls: decls,
	}
}

// FormatNode use go standard lib "go/format" to format ast.Node into code.
func FormatNode(node ast.Node) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}

	formated, err := format.Source(buf.Bytes())
	if err != nil {
		return "", err
	}
	return string(formated), nil
}

// FormatNodeExtra use "https://github.com/mvdan/gofumpt" to format ast.Node into code.
func FormatNodeExtra(node ast.Node) (string, error) {
	var buf bytes.Buffer
	fset := token.NewFileSet()

	if err := format.Node(&buf, fset, node); err != nil {
		return "", err
	}

	formatted, err := fumpt.Source(buf.Bytes(), fumpt.Options{
		LangVersion: "",
		ExtraRules:  true,
	})
	return string(formatted), err
}

func MethodAddComments(code string, modelName string) string {
	for _, method := range Methods {
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
