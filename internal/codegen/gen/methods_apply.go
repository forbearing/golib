package gen

import (
	"go/ast"

	"github.com/forbearing/golib/dsl"
)

// ApplyServiceMethod1 updates functions that match the ServiceMethod1 shape according to DSL.
// Currently ServiceMethod1 does not rely on DSL configuration; keep empty for future extension.
func ApplyServiceMethod1(fn *ast.FuncDecl, action *dsl.Action) {}

// ApplyServiceMethod2 updates functions that match the ServiceMethod2 shape according to DSL.
// Currently ServiceMethod2 does not rely on DSL configuration; keep empty for future extension.
func ApplyServiceMethod2(fn *ast.FuncDecl, action *dsl.Action) {}

// ApplyServiceMethod3 updates functions that match the ServiceMethod3 shape according to DSL.
// Currently ServiceMethod3 does not rely on DSL configuration; keep empty for future extension.
func ApplyServiceMethod3(fn *ast.FuncDecl, action *dsl.Action) {}

// ApplyServiceMethod4 updates functions that match the ServiceMethod4 shape based on the DSL.
// It only updates the shape of *ast.FuncDecl (param/return types) and never touches the method body logic.
// Shape: func (r *recv) Method(ctx *types.ServiceContext, req *<pkg>.<Req>) (*<pkg>.<Rsp>, error)
func ApplyServiceMethod4(fn *ast.FuncDecl, action *dsl.Action) {
	if fn == nil || action == nil {
		return
	}
	if !IsServiceMethod4(fn) {
		return
	}

	// Update the second parameter type to *pkg.<Payload>
	if fn.Type != nil && fn.Type.Params != nil && len(fn.Type.Params.List) >= 2 {
		param := fn.Type.Params.List[1]
		if star, ok := param.Type.(*ast.StarExpr); ok {
			if sel, ok := star.X.(*ast.SelectorExpr); ok {
				if _, ok := sel.X.(*ast.Ident); ok {
					if action.Payload != "" {
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
					if action.Result != "" {
						sel.Sel = ast.NewIdent(action.Result)
					}
				}
			}
		}
	}
}

// ApplyServiceType updates a service struct type to match DSL Payload/Result generics.
// It transforms: type user struct { service.Base[*model.User, *model.User, *model.User] }
// into:         type user struct { service.Base[*model.User, *model.UserReq, *model.UserRsp] }
func ApplyServiceType(spec *ast.TypeSpec, action *dsl.Action) {
	if spec == nil || action == nil || action.Payload == "" || action.Result == "" {
		return
	}
	structType, ok := spec.Type.(*ast.StructType)
	if !ok || structType.Fields == nil {
		return
	}
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
								sel2.Sel = ast.NewIdent(action.Payload)
							}
						}
						if star3, ok := indexListExpr.Indices[2].(*ast.StarExpr); ok {
							if sel3, ok := star3.X.(*ast.SelectorExpr); ok {
								// keep package (sel3.X), replace Sel with action.Result
								sel3.Sel = ast.NewIdent(action.Result)
							}
						}
					}
				}
			}
		}
	}
}
