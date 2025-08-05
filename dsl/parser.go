package dsl

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"github.com/forbearing/golib/types/consts"
)

func Parse(file *ast.File) map[string]*Design {
	m := make(map[string]*Design)
	designs := parse(file)
	for name, fnDecl := range designs {
		if fnDecl == nil {
			continue
		}
		design := parseDesign(fnDecl)
		if len(design.Endpoint) == 0 {
			// Default endpoint is the lower case of the model name.
			design.Endpoint = strings.ToLower(name)
		}
		m[name] = design
	}

	return m
}

// parse will parse the whole file node to find all models with its "Design" method node.
//
// key is model name, value is ast node that represents the model "Design" method.
func parse(file *ast.File) map[string]*ast.FuncDecl {
	designs := make(map[string]*ast.FuncDecl)
	if file == nil {
		return designs
	}

	models := findAllModelNames(file)

	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok && fn != nil {
			if fn.Name == nil || len(fn.Name.Name) == 0 {
				continue
			}

			// Check if the model has method "Design"
			if fn.Name.Name != "Design" {
				continue
			}
			// Check if the method receiver name is the model name.
			if fn.Recv == nil || len(fn.Recv.List) == 0 {
				continue
			}
			var recvName string
			switch t := fn.Recv.List[0].Type.(type) {
			case *ast.Ident:
				if t != nil {
					recvName = t.Name
				}
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok && ident != nil {
					recvName = ident.Name
				}
			}
			if slices.Contains(models, recvName) {
				designs[recvName] = fn
			}
		}
	}

	return designs
}

// parseDesign parse the *ast.FuncDecl that represents "Design" method and returns a *Design object.
func parseDesign(fn *ast.FuncDecl) *Design {
	design := &Design{}
	if fn == nil || fn.Body == nil || len(fn.Body.List) == 0 {
		return design
	}
	stmts := fn.Body.List

	// Set default values
	design.Enabled = true // default enabled.

	for _, stmt := range stmts {
		callExpr, ok := stmt.(*ast.ExprStmt)
		if !ok || callExpr == nil {
			continue
		}
		call, ok := callExpr.X.(*ast.CallExpr)
		if !ok || call == nil || call.Fun == nil || len(call.Args) == 0 {
			continue
		}
		var funcName string
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			if fun == nil {
				continue
			}
			funcName = fun.Name
		case *ast.SelectorExpr:
			if fun == nil || fun.Sel == nil {
				continue
			}
			funcName = fun.Sel.Name
		default:
			continue
		}
		if !is(funcName) {
			continue
		}

		// Parse "Enabled" design.
		if funcName == "Enabled" && len(call.Args) == 1 {
			arg, ok := call.Args[0].(*ast.Ident)
			if ok && arg != nil {
				design.Enabled = arg.Name == "true"
			}
		}

		// Parse "Endpoint" design
		if funcName == "Endpoint" && len(call.Args) == 1 {
			if arg, ok := call.Args[0].(*ast.BasicLit); ok && arg != nil && arg.Kind == token.STRING {
				design.Endpoint = trimQuote(arg.Value)
			}
		}

		if payload, result, exists := parseAction(consts.PHASE_CREATE.MethodName(), funcName, call.Args); exists {
			design.Create = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_DELETE.MethodName(), funcName, call.Args); exists {
			design.Delete = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_UPDATE.MethodName(), funcName, call.Args); exists {
			design.Update = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_PATCH.MethodName(), funcName, call.Args); exists {
			design.Patch = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_LIST.MethodName(), funcName, call.Args); exists {
			design.List = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_GET.MethodName(), funcName, call.Args); exists {
			design.Get = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_CREATE_MANY.MethodName(), funcName, call.Args); exists {
			design.CreateMany = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_DELETE_MANY.MethodName(), funcName, call.Args); exists {
			design.DeleteMany = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_UPDATE_MANY.MethodName(), funcName, call.Args); exists {
			design.UpdateMany = &Action{Payload: payload, Result: result}
		}
		if payload, result, exists := parseAction(consts.PHASE_PATCH_MANY.MethodName(), funcName, call.Args); exists {
			design.PatchMany = &Action{Payload: payload, Result: result}
		}

	}

	return design
}

func parseAction(name string, funcName string, args []ast.Expr) (string, string, bool) {
	var payload string
	var result string

	if funcName == name && len(args) == 1 {
		if flit, ok := args[0].(*ast.FuncLit); ok && flit != nil && flit.Body != nil {
			// Payload or Result.
			for _, stmt := range flit.Body.List {
				if expr, ok := stmt.(*ast.ExprStmt); ok && expr != nil {
					if call, ok := expr.X.(*ast.CallExpr); ok && call != nil && call.Fun != nil {
						if identExpr, ok := call.Fun.(*ast.IndexExpr); ok && identExpr != nil {
							var isPayload bool
							var isResult bool
							var funcName string
							switch x := identExpr.X.(type) {
							case *ast.Ident:
								if x != nil {
									funcName = x.Name
								}
							case *ast.SelectorExpr:
								if x != nil && x.Sel != nil {
									funcName = x.Sel.Name
								}
							}
							switch funcName {
							case "Payload":
								isPayload = true
							case "Result":
								isResult = true
							}
							if isPayload {
								if ident, ok := identExpr.Index.(*ast.Ident); ok && ident != nil { // Payload[User]
									payload = ident.Name
								} else if starExpr, ok := identExpr.Index.(*ast.StarExpr); ok && starExpr != nil { // Payload[*User]
									if ident, ok := starExpr.X.(*ast.Ident); ok && ident != nil {
										payload = ident.Name
									}
								}
							}
							if isResult {
								if ident, ok := identExpr.Index.(*ast.Ident); ok && ident != nil { // Result[User]
									result = ident.Name
								} else if starExpr, ok := identExpr.Index.(*ast.StarExpr); ok && starExpr != nil { // Result[*User]
									if ident, ok := starExpr.X.(*ast.Ident); ok && ident != nil {
										result = ident.Name
									}
								}
							}
						}
					}
				}
			}
		}
		return payload, result, true
	}

	return "", "", false
}

// findAllModelNames finds all model names in the ast File Node.
func findAllModelNames(file *ast.File) []string {
	names := make([]string, 0)
	if file == nil {
		return names
	}
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl == nil || genDecl.Tok != token.TYPE {
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
			for _, field := range structType.Fields.List {
				if isModelBase(file, field) {
					names = append(names, typeSpec.Name.Name)
					break
				}
			}
		}
	}

	return names
}

// isModelBase check if the struct field is "model.Base" or "aliasmodelname.Base"
//
/*
import (
	"github.com/forbearing/golib/dsl"
	. "github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/model"
	pkgmodel "github.com/forbearing/golib/model"
)
*/
func isModelBase(file *ast.File, field *ast.Field) bool {
	// Not anonymouse field.
	if file == nil || field == nil || len(field.Names) != 0 {
		return false
	}

	aliasNames := []string{"model"}
	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}
		if imp.Path.Value == consts.IMPORT_PATH_MODEL {
			if imp.Name != nil && !slices.Contains(aliasNames, imp.Name.Name) {
				aliasNames = append(aliasNames, imp.Name.Name)
			}
		}
	}

	switch t := field.Type.(type) {
	case *ast.SelectorExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return slices.Contains(aliasNames, ident.Name) && t.Sel.Name == "Base"
		}
	case *ast.Ident:
		return t.Name == "Base"
	}

	return false
}
