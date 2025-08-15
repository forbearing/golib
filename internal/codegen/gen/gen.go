package gen

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/types/consts"
	"github.com/stoewer/go-strcase"
)

// ModelInfo stores model information
//
// Examples:
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model", ModelName:"User", ModelVarName:"u", ModelFileDir:"/tmp/model"},
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model", ModelName:"Group", ModelVarName:"g", ModelFileDir:"/tmp/model"},
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model_auth", ModelName:"User", ModelVarName:"u", ModelFileDir:"/tmp/model"},
// {ModulePath:"github.com/forbearing/golib", ModelPkgName:"model_auth", ModelName:"Group", ModelVarName:"g", ModelFileDir:"/tmp/model"},
type ModelInfo struct {
	// module related fields
	ModulePath string // module path parsed from go.mod

	// model related fields
	ModelPkgName  string // model package name, e.g.: model, model_authz, model_log
	ModelName     string // model name, e.g.: User, Group
	ModelVarName  string // lowercase model variable name, e.g.: u, g
	ModelFileDir  string // relative path of model file directory, e.g.: github.com/forbearing/golib/model
	ModelFilePath string // relative path of model file, e.g.: github.com/forbearing/golib/model/user.go

	// custom request and response related fields
	Design *dsl.Design
}

func (m *ModelInfo) ServiceImportPath(modelDir, serviceDir string) string {
	path := strings.Replace(filepath.Join(m.ModulePath, m.ModelFilePath), modelDir, serviceDir, 1)
	path = strings.TrimRight(path, ".go")
	return path
}

func (m *ModelInfo) RouterImportPath() string {
	return filepath.Join(m.ModulePath, m.ModelFileDir)
}

func (m *ModelInfo) ModelImportPath() (string, bool) {
	// If a struct anonymous inherits from model.Base, than the model will be imported in model/model.go using
	// statement such like: "model.Register[*User]()".
	// Imported the model is not determinated by m.Design.Eanbled value.
	path := filepath.Join(m.ModulePath, m.ModelFileDir)
	if !strings.HasSuffix(path, "/model") {
		return path, true
	}
	return "", false
}

// GetModulePath parses go.mod to get module path
func GetModulePath() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// If go command exists, get module path directly through go list -m command
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err == nil {
		return strings.TrimSpace(string(output)), nil
	}

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

// findModelPackageName finds the actual name of the imported model package
// import "github.com/forbearing/golib/model" returns "model"
// import model_auth "github.com/forbearing/golib/model" returns model_auth
func findModelPackageName(file *ast.File) string {
	return file.Name.Name
}

// // isModelBase 检查字段是否是 model.Base
//
//	func isModelBase(file *ast.File, field *ast.Field, modelPkgName string) bool {
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
func isModelBase(file *ast.File, field *ast.Field) bool {
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

// FindModels finds all structs in model files
func FindModels(module string, modelDir string, filename string) ([]*ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	modelPkgName := findModelPackageName(node)
	if len(modelPkgName) == 0 {
		return nil, fmt.Errorf("file %s has no model package", filename)
	}
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	designs := dsl.Parse(f, "")
	for _, design := range designs {
		// The new endpoint value is the model file dir + the endpoint value
		// For example: old endpoint is "order", the model dir is "model/user",
		// then the new endpoint is "user/order"
		newFilename := strings.TrimPrefix(filename, modelDir) // "/user/order.go"
		newFilename = strings.TrimPrefix(newFilename, "/")    // "user/order.go"
		dir := filepath.Dir(newFilename)                      // "user"
		design.Endpoint = filepath.Join(dir, design.Endpoint) // "user/order"
	}

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
				if isModelBase(node, field) {
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

// modelPkg2ServicePkg converts model name to service name.
func modelPkg2ServicePkg(pkgName string) string {
	if pkgName == "model" {
		return "service"
	}
	// For model_xxx format, replace with service_xxx
	if strings.HasPrefix(pkgName, "model_") {
		return strings.Replace(pkgName, "model_", "service_", 1)
	}
	return strings.Replace(pkgName, "model", "service", 1)
}

// genServiceMethod1 uses AST to generate CreateBefore,CreateAfter,UpdateBefore,UpdateAfter,
// DeleteBefore,DeleteAfter,GetBefore,GetAfter,UpdatePartialBefore,UpdatePartialAfter methods.
func genServiceMethod1(info *ModelInfo, phase consts.Phase) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(phase.MethodName()), "_", " ")

	return serviceMethod1(info.ModelVarName, info.ModelName, info.ModelPkgName, phase,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(ast.NewIdent("nil")),
	)
}

// genServiceMethod2 uses AST to generate ListBefore, ListAfter methods.
func genServiceMethod2(info *ModelInfo, phase consts.Phase) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(phase.MethodName()), "_", " ")

	return serviceMethod2(info.ModelVarName, info.ModelName, info.ModelPkgName, phase,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(ast.NewIdent("nil")),
	)
}

// genServiceMethod3 uses AST to generate CreateManyBefore, CreateManyAfter,
// DeleteManyBefore, DeleteManyAfter, UpdateManyBefore, UpdateManyAfter, PatchManyBefore, PatchManyAfter.
func genServiceMethod3(info *ModelInfo, phase consts.Phase) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(phase.MethodName()), "_", " ")

	return serviceMethod3(info.ModelVarName, info.ModelName, info.ModelPkgName, phase,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(ast.NewIdent("nil")),
	)
}

// genServiceMethod4 uses AST to generate Create,Delete,Update,Patch,List,Get,CreateMany,DeleteMany,UpdateMany,PatchMany methods.
func genServiceMethod4(info *ModelInfo, reqName, rspName string, phase consts.Phase) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(phase.MethodName()), "_", " ")

	return serviceMethod4(info.ModelVarName, info.ModelName, info.ModelPkgName, reqName, rspName, phase,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(
			ast.NewIdent("rsp"),
			ast.NewIdent("nil"),
		),
	)
}

// genServiceMethod5 uses AST to generate Import method.
func genServiceMethod5(info *ModelInfo, phase consts.Phase) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(phase.MethodName()), "_", " ")

	return serviceMethod5(info.ModelVarName, info.ModelName, info.ModelPkgName, phase,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(ast.NewIdent(pluralizeCli.Plural(strings.ToLower(info.ModelName))), ast.NewIdent("err")),
	)
}

// genServiceMethod6 uses AST to generate Export method.
func genServiceMethod6(info *ModelInfo, phase consts.Phase) *ast.FuncDecl {
	str := strings.ReplaceAll(strcase.SnakeCase(phase.MethodName()), "_", " ")

	return serviceMethod6(info.ModelVarName, info.ModelName, info.ModelPkgName, phase,
		StmtLogWithServiceContext(info.ModelVarName),
		StmtLogInfo(fmt.Sprintf(`"%s %s"`, strings.ToLower(info.ModelName), str)),
		EmptyLine(),
		Returns(ast.NewIdent("data"), ast.NewIdent("err")),
	)
}

func GenerateService(info *ModelInfo, action *dsl.Action, phase consts.Phase) *ast.File {
	if !action.Enabled || !action.Service {
		return nil
	}

	otherPkgs := []string{}
	if phase == consts.PHASE_IMPORT {
		otherPkgs = append(otherPkgs, "io")
	}

	decls := []ast.Decl{
		imports(info.ModulePath, info.ModelFileDir, info.ModelPkgName, otherPkgs...),
		// Inits(info.ModelName),
		// Types(info.ModelPkgName, info.ModelName, info.Design.Create.Payload, info.Design.Create.Result),
	}

	// add types
	if action.Enabled {
		decls = append(decls, types(info.ModelPkgName, info.ModelName, action.Payload, action.Result, phase, false))
	}

	// add methods
	switch phase {
	case consts.PHASE_CREATE:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod1(info, phase.Before()))
		decls = append(decls, genServiceMethod1(info, phase.After()))
	case consts.PHASE_DELETE:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod1(info, phase.Before()))
		decls = append(decls, genServiceMethod1(info, phase.After()))
	case consts.PHASE_UPDATE:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod1(info, phase.Before()))
		decls = append(decls, genServiceMethod1(info, phase.After()))
	case consts.PHASE_PATCH:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod1(info, phase.Before()))
		decls = append(decls, genServiceMethod1(info, phase.After()))
	case consts.PHASE_LIST: // List method use GenerateServiceMethod2
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod2(info, phase.Before()))
		decls = append(decls, genServiceMethod2(info, phase.After()))
	case consts.PHASE_GET:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod1(info, phase.Before()))
		decls = append(decls, genServiceMethod1(info, phase.After()))
	case consts.PHASE_CREATE_MANY: // XXXMany methods use GenerateServiceMethod3
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod3(info, phase.Before()))
		decls = append(decls, genServiceMethod3(info, phase.After()))
	case consts.PHASE_DELETE_MANY:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod3(info, phase.Before()))
		decls = append(decls, genServiceMethod3(info, phase.After()))
	case consts.PHASE_UPDATE_MANY:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod3(info, phase.Before()))
		decls = append(decls, genServiceMethod3(info, phase.After()))
	case consts.PHASE_PATCH_MANY:
		decls = append(decls, genServiceMethod4(info, action.Payload, action.Result, phase))
		decls = append(decls, genServiceMethod3(info, phase.Before()))
		decls = append(decls, genServiceMethod3(info, phase.After()))
	case consts.PHASE_IMPORT:
		decls = append(decls, genServiceMethod5(info, phase))
	case consts.PHASE_EXPORT:
		decls = append(decls, genServiceMethod6(info, phase))
	}

	return &ast.File{
		Name:  ast.NewIdent(strings.ToLower(info.ModelName)),
		Decls: decls,
	}
}
