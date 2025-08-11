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
	"github.com/forbearing/golib/types/consts"
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
		Returns(ast.NewIdent("nil")),
	)
}

// GenerateServiceMethod2 使用 AST 生成 ListBefore, ListAfter 方法.
func GenerateServiceMethod2(info *ModelInfo, methodName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return ServiceMethod2(info.ModelVarName, info.ModelName, methodName, info.ModelPkgName,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(ast.NewIdent("nil")),
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
		Returns(ast.NewIdent("nil")),
	)
}

// GenerateServiceMethod4 使用 AST 生成 Create,Delete,Update,Patch,List,Get,CreateMany,DeleteMany,UpdateMany,PatchMany 方法.
func GenerateServiceMethod4(info *ModelInfo, methodName, reqName, rspName string) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(methodName), "_", " ")

	return ServiceMethod4(info.ModelVarName, info.ModelName, methodName, info.ModelPkgName, reqName, rspName,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(
			ast.NewIdent("rsp"),
			ast.NewIdent("nil"),
		),
	)
}

func GenerateService(info *ModelInfo, action *dsl.Action, phase consts.Phase) *ast.File {
	if !IsValidModelPackage(info.ModelPkgName) || !action.Enabled {
		return nil
	}

	decls := []ast.Decl{
		Imports(info.ModulePath, info.ModelFileDir, info.ModelPkgName),
		// Inits(info.ModelName),
		// Types(info.ModelPkgName, info.ModelName, info.Design.Create.Payload, info.Design.Create.Result),
	}

	if info.Design.Create.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.Create.Payload, info.Design.Create.Result, false))
	}
	if info.Design.Delete.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.Delete.Payload, info.Design.Delete.Result, false))
	}
	if info.Design.Update.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.Update.Payload, info.Design.Update.Result, false))
	}
	if info.Design.Patch.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.Patch.Payload, info.Design.Patch.Result, false))
	}
	if info.Design.List.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.List.Payload, info.Design.List.Result, false))
	}
	if info.Design.Get.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.Get.Payload, info.Design.Get.Result, false))
	}
	if info.Design.CreateMany.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.CreateMany.Payload, info.Design.CreateMany.Result, false))
	}
	if info.Design.DeleteMany.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.DeleteMany.Payload, info.Design.DeleteMany.Result, false))
	}
	if info.Design.UpdateMany.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.UpdateMany.Payload, info.Design.UpdateMany.Result, false))
	}
	if info.Design.PatchMany.Enabled {
		decls = append(decls, Types(info.ModelPkgName, info.ModelName, info.Design.PatchMany.Payload, info.Design.PatchMany.Result, false))
	}

	switch phase {
	case consts.PHASE_CREATE:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_CREATE_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_CREATE_AFTER.MethodName()))
	case consts.PHASE_DELETE:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_DELETE_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_DELETE_AFTER.MethodName()))
	case consts.PHASE_UPDATE:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_UPDATE_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_UPDATE_AFTER.MethodName()))
	case consts.PHASE_PATCH:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_PATCH_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_PATCH_AFTER.MethodName()))
	case consts.PHASE_LIST: // List method use GenerateServiceMethod2
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod2(info, consts.PHASE_LIST_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod2(info, consts.PHASE_LIST_AFTER.MethodName()))
	case consts.PHASE_GET:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_GET_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod1(info, consts.PHASE_GET_AFTER.MethodName()))
	case consts.PHASE_CREATE_MANY: // XXXMany methods use GenerateServiceMethod3
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_CREATE_MANY_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_CREATE_MANY_AFTER.MethodName()))
	case consts.PHASE_DELETE_MANY:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_DELETE_MANY_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_DELETE_MANY_AFTER.MethodName()))
	case consts.PHASE_UPDATE_MANY:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_UPDATE_MANY_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_UPDATE_MANY_AFTER.MethodName()))
	case consts.PHASE_PATCH_MANY:
		decls = append(decls, GenerateServiceMethod4(info, phase.MethodName(), action.Payload, action.Result))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_PATCH_MANY_BEFORE.MethodName()))
		decls = append(decls, GenerateServiceMethod3(info, consts.PHASE_PATCH_MANY_AFTER.MethodName()))
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
