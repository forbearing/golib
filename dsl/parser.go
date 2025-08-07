package dsl

import (
	"go/ast"
	"go/token"
	"slices"
	"strings"

	"github.com/forbearing/golib/types/consts"
)

// Parse parses the whole file node to find all models with its "Design".
// returns the map that key is model name, value is *Design.
func Parse(file *ast.File) map[string]*Design {
	m := make(map[string]*Design)
	designs := parse(file)
	for name, fnDecl := range designs {
		design := parseDesign(fnDecl)
		// Default endpoint is the lower case of the model name.
		if len(design.Endpoint) == 0 {
			design.Endpoint = strings.ToLower(name)
		}
		if design.Create == nil {
			design.Create = &Action{Payload: name, Result: name}
		}
		if design.Delete == nil {
			design.Delete = &Action{Payload: name, Result: name}
		}
		if design.Update == nil {
			design.Update = &Action{Payload: name, Result: name}
		}
		if design.Patch == nil {
			design.Patch = &Action{Payload: name, Result: name}
		}
		if design.List == nil {
			design.List = &Action{Payload: name, Result: name}
		}
		if design.Get == nil {
			design.Get = &Action{Payload: name, Result: name}
		}
		if design.CreateMany == nil {
			design.CreateMany = &Action{Payload: name, Result: name}
		}
		if design.DeleteMany == nil {
			design.DeleteMany = &Action{Payload: name, Result: name}
		}
		if design.UpdateMany == nil {
			design.UpdateMany = &Action{Payload: name, Result: name}
		}
		if design.PatchMany == nil {
			design.PatchMany = &Action{Payload: name, Result: name}
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
	// Every model should always has a *ast.FuncDecl,
	// If model has no "Design" method, then the value is nil.
	// It's convenient to generate a default design for the model.
	for _, model := range models {
		designs[model] = nil
	}

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
			} else {
				designs[recvName] = nil
			}
		}
	}

	return designs
}

// parseDesign parse the *ast.FuncDecl that represents "Design" method and returns a *Design object.
func parseDesign(fn *ast.FuncDecl) *Design {
	defaults := &Design{Enabled: true}
	// model don't have "Design" method, so returns the default design values.
	if fn == nil || fn.Body == nil || len(fn.Body.List) == 0 {
		return defaults
	}
	stmts := fn.Body.List

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
				defaults.Enabled = arg.Name == "true"
			}
		}

		// Parse "Endpoint" design
		if funcName == "Endpoint" && len(call.Args) == 1 {
			if arg, ok := call.Args[0].(*ast.BasicLit); ok && arg != nil && arg.Kind == token.STRING {
				defaults.Endpoint = trimQuote(arg.Value)
			}
		}

		if payload, result, enabled, exists := parseAction(consts.PHASE_CREATE.MethodName(), funcName, call.Args); exists {
			defaults.Create = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_DELETE.MethodName(), funcName, call.Args); exists {
			defaults.Delete = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_UPDATE.MethodName(), funcName, call.Args); exists {
			defaults.Update = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_PATCH.MethodName(), funcName, call.Args); exists {
			defaults.Patch = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_LIST.MethodName(), funcName, call.Args); exists {
			defaults.List = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_GET.MethodName(), funcName, call.Args); exists {
			defaults.Get = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_CREATE_MANY.MethodName(), funcName, call.Args); exists {
			defaults.CreateMany = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_DELETE_MANY.MethodName(), funcName, call.Args); exists {
			defaults.DeleteMany = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_UPDATE_MANY.MethodName(), funcName, call.Args); exists {
			defaults.UpdateMany = &Action{Payload: payload, Result: result, Enabled: enabled}
		}
		if payload, result, enabled, exists := parseAction(consts.PHASE_PATCH_MANY.MethodName(), funcName, call.Args); exists {
			defaults.PatchMany = &Action{Payload: payload, Result: result, Enabled: enabled}
		}

	}

	return defaults
}

// parseAction parse the "Payload" and "Result" type from Action function.
// The "Action" is represented by function name that already defined in the method list.
func parseAction(name string, funcName string, args []ast.Expr) (string, string, bool, bool) {
	var payload string
	var result string
	var enabled bool

	if funcName == name && len(args) == 1 {
		if flit, ok := args[0].(*ast.FuncLit); ok && flit != nil && flit.Body != nil {
			for _, stmt := range flit.Body.List {
				if expr, ok := stmt.(*ast.ExprStmt); ok && expr != nil {
					if call, ok := expr.X.(*ast.CallExpr); ok && call != nil && call.Fun != nil {

						// Parse Enabled(true)/Enabled(false)
						var isEnabledCall bool
						switch fun := call.Fun.(type) {
						case *ast.Ident:
							// anonymous import: Enabled(true)
							if fun != nil && fun.Name == "Enabled" {
								isEnabledCall = true
								enabled = true
							}
						case *ast.SelectorExpr:
							// non-anonymous import: dsl.Enabled(true)
							if fun != nil && fun.Sel != nil && fun.Sel.Name == "Enabled" {
								isEnabledCall = true
								enabled = true
							}
						}

						if isEnabledCall && enabled && len(call.Args) > 0 && call.Args[0] != nil {
							if identExpr, ok := call.Args[0].(*ast.Ident); ok && identExpr != nil {
								// check the argument of Enabled() is true.
								enabled = enabled && identExpr.Name == "true"
							}
						}

						// Parse Payload[User] or Result[*User].
						if indexExpr, ok := call.Fun.(*ast.IndexExpr); ok && indexExpr != nil {
							var isPayload bool
							var isResult bool
							var funcName string
							switch x := indexExpr.X.(type) {
							case *ast.Ident:
								// anonymous import: Payload[User]
								if x != nil {
									funcName = x.Name
								}
							case *ast.SelectorExpr:
								// non-anonymous import: dsl.Payload[User]
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
								if ident, ok := indexExpr.Index.(*ast.Ident); ok && ident != nil { // Payload[User]
									payload = ident.Name
								} else if starExpr, ok := indexExpr.Index.(*ast.StarExpr); ok && starExpr != nil { // Payload[*User]
									if ident, ok := starExpr.X.(*ast.Ident); ok && ident != nil {
										payload = ident.Name
									}
								}
							}
							if isResult {
								if ident, ok := indexExpr.Index.(*ast.Ident); ok && ident != nil { // Result[User]
									result = ident.Name
								} else if starExpr, ok := indexExpr.Index.(*ast.StarExpr); ok && starExpr != nil { // Result[*User]
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
		return payload, result, enabled, true
	}

	return "", "", false, false
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
