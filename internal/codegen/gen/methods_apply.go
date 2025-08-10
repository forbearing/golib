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
