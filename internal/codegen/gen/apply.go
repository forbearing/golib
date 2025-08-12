package gen

import (
	"go/ast"
	"go/token"

	"github.com/forbearing/golib/dsl"
)

// ApplyServiceFile will apply the dsl.Action to the ast.File.
// It will modify the struct type and struct methods if Payload
// or Result is changed, and returns true.
// Otherwise returns false.
func ApplyServiceFile(file *ast.File, action *dsl.Action) bool {
	if file == nil || action == nil {
		return false
	}

	var changed bool

	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if isServiceType(typeSpec) {
						if applyServiceType(typeSpec, action) {
							changed = true
						}
					}
				}
			}
		}
	}

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl != nil {
			if isServiceMethod1(funcDecl) {
				if applyServiceMethod1(funcDecl, action) {
					changed = true
				}
			}
			if isServiceMethod2(funcDecl) {
				if applyServiceMethod2(funcDecl, action) {
					changed = true
				}
			}
			if isServiceMethod3(funcDecl) {
				if applyServiceMethod3(funcDecl, action) {
					changed = true
				}
			}
			if isServiceMethod4(funcDecl) {
				if applyServiceMethod4(funcDecl, action) {
					changed = true
				}
			}
		}
	}

	return changed
}

// applyServiceMethod1 updates functions that match the ServiceMethod1 shape according to DSL.
// Currently ServiceMethod1 does not rely on DSL configuration; keep empty for future extension.
func applyServiceMethod1(fn *ast.FuncDecl, action *dsl.Action) bool { return false }

// applyServiceMethod2 updates functions that match the ServiceMethod2 shape according to DSL.
// Currently ServiceMethod2 does not rely on DSL configuration; keep empty for future extension.
func applyServiceMethod2(fn *ast.FuncDecl, action *dsl.Action) bool { return false }

// applyServiceMethod3 updates functions that match the ServiceMethod3 shape according to DSL.
// Currently ServiceMethod3 does not rely on DSL configuration; keep empty for future extension.
func applyServiceMethod3(fn *ast.FuncDecl, action *dsl.Action) bool { return false }

// applyServiceMethod4 updates functions that match the ServiceMethod4 shape based on the DSL.
// It only updates the shape of *ast.FuncDecl (param/return types) and never touches the method body logic.
// Shape: func (r *recv) Method(ctx *types.ServiceContext, req *<pkg>.<Req>) (*<pkg>.<Rsp>, error)
func applyServiceMethod4(fn *ast.FuncDecl, action *dsl.Action) bool {
	if fn == nil || action == nil {
		return false
	}
	if !isServiceMethod4(fn) {
		return false
	}

	var changed bool

	// Update the second parameter type to *pkg.<Payload>
	if fn.Type != nil && fn.Type.Params != nil && len(fn.Type.Params.List) >= 2 {
		param := fn.Type.Params.List[1]
		if star, ok := param.Type.(*ast.StarExpr); ok {
			if sel, ok := star.X.(*ast.SelectorExpr); ok {
				if _, ok := sel.X.(*ast.Ident); ok {
					if action.Payload != "" && sel.Sel.Name != action.Payload {
						changed = true
						sel.Sel = ast.NewIdent(action.Payload)
					}
				}
			}
		}
	}

	// Update the first result type to *pkg.<Result>
	if fn.Type != nil && fn.Type.Results != nil && len(fn.Type.Results.List) >= 1 {
		res := fn.Type.Results.List[0]
		if star, ok := res.Type.(*ast.StarExpr); ok {
			if sel, ok := star.X.(*ast.SelectorExpr); ok {
				if _, ok := sel.X.(*ast.Ident); ok {
					if action.Result != "" && sel.Sel.Name != action.Result {
						changed = true
						sel.Sel = ast.NewIdent(action.Result)
					}
				}
			}
		}
	}

	return changed
}

// applyServiceType updates a service struct type to match DSL Payload/Result generics.
// It transforms: type user struct { service.Base[*model.User, *model.User, *model.User] }
// into:         type user struct { service.Base[*model.User, *model.UserReq, *model.UserRsp] }
func applyServiceType(spec *ast.TypeSpec, action *dsl.Action) bool {
	if spec == nil || action == nil || action.Payload == "" || action.Result == "" {
		return false
	}
	structType, ok := spec.Type.(*ast.StructType)
	if !ok || structType.Fields == nil {
		return false
	}

	var changed bool

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 { // Embedded field
			indexListExpr, ok := field.Type.(*ast.IndexListExpr)
			if !ok {
				continue
			}
			// ensure service.Base
			if sel, ok := indexListExpr.X.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := sel.X.(*ast.Ident); ok && pkgIdent.Name == "service" && sel.Sel.Name == "Base" {
					if len(indexListExpr.Indices) == 3 {
						// second -> Payload, third -> Result
						if star2, ok := indexListExpr.Indices[1].(*ast.StarExpr); ok {
							if sel2, ok := star2.X.(*ast.SelectorExpr); ok {
								// keep package (sel2.X), replace Sel with action.Payload
								if sel2.Sel.Name != action.Payload {
									changed = true
									sel2.Sel = ast.NewIdent(action.Payload)
								}
							}
						}
						if star3, ok := indexListExpr.Indices[2].(*ast.StarExpr); ok {
							if sel3, ok := star3.X.(*ast.SelectorExpr); ok {
								// keep package (sel3.X), replace Sel with action.Result
								if sel3.Sel.Name != action.Result {
									changed = true
									sel3.Sel = ast.NewIdent(action.Result)
								}
							}
						}
					}
				}
			}
		}
	}

	return changed
}
