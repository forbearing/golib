package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// ModelInfo 存储模型信息
type ModelInfo struct {
	PackageName   string
	ModelName     string
	ModelVarName  string // 小写的模型变量名
	ModulePath    string // 从 go.mod 解析的模块路径
	ModelFilePath string // model 文件的相对路径
}

// 解析 go.mod 获取模块路径
func getModulePath() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}
	return "", fmt.Errorf("module path not found in go.mod")
}

// 查找包中导入的 model 包的实际名称
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

// 检查字段是否是 model.Base
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

// 查找 model 文件中的所有结构体
func findModels(filename string, modulePath string) ([]ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var models []ModelInfo
	packageName := node.Name.Name

	// 获取 model 包的实际名称
	modelPkgName := findModelPackageName(node)
	if modelPkgName == "" && packageName == "model" {
		modelPkgName = "" // 同包情况
	}

	// 遍历所有声明
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			// 检查是否嵌入了 model.Base
			hasModelBase := false
			for _, field := range structType.Fields.List {
				if isModelBase(field, modelPkgName) {
					hasModelBase = true
					break
				}
			}

			if hasModelBase {
				modelName := typeSpec.Name.Name
				models = append(models, ModelInfo{
					PackageName:   packageName,
					ModelName:     modelName,
					ModelVarName:  strings.ToLower(modelName[:1]),
					ModulePath:    modulePath,
					ModelFilePath: filepath.Dir(filename),
				})
			}
		}
	}

	return models, nil
}

// 使用 AST 生成 CreateBefore 方法
func generateCreateBeforeAST(model ModelInfo) *ast.FuncDecl {
	// 接收者名称
	receiverName := model.ModelVarName

	// 创建方法声明
	funcDecl := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{ast.NewIdent(receiverName)},
					Type: &ast.StarExpr{
						X: ast.NewIdent(strings.ToLower(model.ModelName)),
					},
				},
			},
		},
		Name: ast.NewIdent("CreateBefore"),
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{ast.NewIdent("ctx")},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("types"),
								Sel: ast.NewIdent("ServiceContext"),
							},
						},
					},
					{
						Names: []*ast.Ident{ast.NewIdent(strings.ToLower(model.ModelName))},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   ast.NewIdent("model"),
								Sel: ast.NewIdent(model.ModelName),
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: ast.NewIdent("error"),
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				// log := u.WithServiceContext(ctx, ctx.GetPhase())
				&ast.AssignStmt{
					Lhs: []ast.Expr{ast.NewIdent("log")},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent(receiverName),
								Sel: ast.NewIdent("WithServiceContext"),
							},
							Args: []ast.Expr{
								ast.NewIdent("ctx"),
								&ast.CallExpr{
									Fun: &ast.SelectorExpr{
										X:   ast.NewIdent("ctx"),
										Sel: ast.NewIdent("GetPhase"),
									},
								},
							},
						},
					},
				},
				// log.Info("user create before")
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("log"),
							Sel: ast.NewIdent("Info"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: fmt.Sprintf(`"%s create before"`, strings.ToLower(model.ModelName)),
							},
						},
					},
				},
				// 空行（通过空语句实现）
				&ast.EmptyStmt{},
				// return nil
				&ast.ReturnStmt{
					Results: []ast.Expr{
						ast.NewIdent("nil"),
					},
				},
			},
		},
	}

	return funcDecl
}

// 使用 AST 生成完整的 service 文件
func generateServiceFileAST(model ModelInfo) (*ast.File, error) {
	// 创建 import 路径
	modelImportPath := filepath.Join(model.ModulePath, model.ModelFilePath)

	// service 类型名
	serviceName := strings.ToLower(model.ModelName)

	// 创建文件 AST
	file := &ast.File{
		Name: ast.NewIdent("service"),
		Decls: []ast.Decl{
			// import 声明
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: fmt.Sprintf(`"%s"`, modelImportPath),
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/forbearing/golib/service"`,
						},
					},
					&ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: `"github.com/forbearing/golib/types"`,
						},
					},
				},
			},
			// init 函数
			&ast.FuncDecl{
				Name: ast.NewIdent("init"),
				Type: &ast.FuncType{
					Params: &ast.FieldList{},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.IndexExpr{
									X: &ast.SelectorExpr{
										X:   ast.NewIdent("service"),
										Sel: ast.NewIdent("Register"),
									},
									Index: &ast.StarExpr{
										X: ast.NewIdent(serviceName),
									},
								},
							},
						},
					},
				},
			},
			// service 结构体声明（带注释）
			&ast.GenDecl{
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{
							Text: fmt.Sprintf("// %s implements the types.Service[*model.%s] interface.",
								serviceName, model.ModelName),
						},
					},
				},
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: ast.NewIdent(serviceName),
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Type: &ast.IndexExpr{
											X: &ast.SelectorExpr{
												X:   ast.NewIdent("service"),
												Sel: ast.NewIdent("Base"),
											},
											Index: &ast.StarExpr{
												X: &ast.SelectorExpr{
													X:   ast.NewIdent("model"),
													Sel: ast.NewIdent(model.ModelName),
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			// CreateBefore 方法
			generateCreateBeforeAST(model),
		},
	}

	return file, nil
}

// 将 AST 转换为格式化的代码，并手动处理注释
func astToCode(file *ast.File, model ModelInfo) (string, error) {
	var buf strings.Builder
	fset := token.NewFileSet()

	err := format.Node(&buf, fset, file)
	if err != nil {
		return "", err
	}

	// 手动添加业务逻辑注释
	code := buf.String()

	// 在 log.Info 之后添加注释
	searchStr := fmt.Sprintf(`log.Info("%s create before")`, strings.ToLower(model.ModelName))
	replaceStr := fmt.Sprintf(`log.Info("%s create before")
	// =============================
	// Add your business logic here.
	// =============================
`, strings.ToLower(model.ModelName))

	code = strings.Replace(code, searchStr, replaceStr, -1)

	return code, nil
}

func main() {
	// 获取模块路径
	modulePath, err := getModulePath()
	if err != nil {
		fmt.Printf("Error getting module path: %v\n", err)
		return
	}
	fmt.Printf("Module path: %s\n", modulePath)

	// 解析 model 文件
	models, err := findModels("model/user.go", modulePath)
	if err != nil {
		fmt.Printf("Error parsing model file: %v\n", err)
		return
	}

	// 为每个 model 生成 service
	for _, model := range models {
		fmt.Printf("Found model: %s\n", model.ModelName)

		// 使用 AST 生成 service 文件
		file, err := generateServiceFileAST(model)
		if err != nil {
			fmt.Printf("Error generating AST for %s: %v\n", model.ModelName, err)
			continue
		}

		// 转换为代码
		code, err := astToCode(file, model)
		if err != nil {
			fmt.Printf("Error converting AST to code for %s: %v\n", model.ModelName, err)
			continue
		}

		fmt.Printf("\nGenerated service for %s:\n%s\n", model.ModelName, code)

		// TODO: 写入文件
		filename := filepath.Join("service", strings.ToLower(model.ModelName)+".go")
		os.WriteFile(filename, []byte(code), 0o644)
	}
}
