package gen

import (
	"go/ast"
	"go/token"

	"github.com/forbearing/gst/dsl"
)

// ApplyServiceFile will apply the dsl.Action to the ast.File.
// It will modify the struct type and struct methods if Payload
// or Result is changed, and returns true.
// Otherwise returns false.
// The servicePkgName parameter specifies the expected package name for the service file.
// This should match the package name used in service registration to maintain consistency.
func ApplyServiceFile(file *ast.File, action *dsl.Action, servicePkgName string) bool {
	if file == nil || action == nil {
		return false
	}

	var changed bool

	// Apply package name correction
	if len(servicePkgName) > 0 && file.Name != nil && file.Name.Name != servicePkgName {
		file.Name.Name = servicePkgName
		changed = true
	}

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
//
//	func (r *recv) Method(ctx *types.ServiceContext, req <pkg>.<Req>) (<pkg>.<Rsp>, error)
func applyServiceMethod4(fn *ast.FuncDecl, action *dsl.Action) bool {
	if fn == nil || action == nil {
		return false
	}

	if !isServiceMethod4(fn) {
		return false
	}

	var changed bool

	// Update the second parameter type based on action.Payload
	if fn.Type != nil && fn.Type.Params != nil && len(fn.Type.Params.List) >= 2 {
		param := fn.Type.Params.List[1]
		if action.Payload != "" {
			// Determine if action.Payload should be a pointer type
			payloadIsPointer := len(action.Payload) > 0 && action.Payload[0] == '*'
			payloadName := action.Payload
			if payloadIsPointer {
				payloadName = action.Payload[1:] // Remove the '*' prefix
			}

			// Handle current *pkg.Type case
			if star, ok := param.Type.(*ast.StarExpr); ok {
				if sel, ok := star.X.(*ast.SelectorExpr); ok {
					if pkgIdent, ok := sel.X.(*ast.Ident); ok {
						if payloadIsPointer {
							// Keep as pointer type, just update the name
							if sel.Sel.Name != payloadName {
								changed = true
								newIdent := ast.NewIdent(payloadName)
								newIdent.NamePos = sel.Sel.NamePos
								sel.Sel = newIdent
							}
						} else {
							// Convert from pointer to non-pointer type
							changed = true
							newSel := &ast.SelectorExpr{
								X:   pkgIdent,
								Sel: ast.NewIdent(payloadName),
							}
							param.Type = newSel
						}
					}
				}
				// Handle current pkg.Type case
			} else if sel, ok := param.Type.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := sel.X.(*ast.Ident); ok {
					if payloadIsPointer {
						// Convert from non-pointer to pointer type
						changed = true
						newStar := &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   pkgIdent,
								Sel: ast.NewIdent(payloadName),
							},
						}
						param.Type = newStar
					} else {
						// Keep as non-pointer type, just update the name
						if sel.Sel.Name != payloadName {
							changed = true
							newIdent := ast.NewIdent(payloadName)
							newIdent.NamePos = sel.Sel.NamePos
							sel.Sel = newIdent
						}
					}
				}
			}
		}
	}

	// Update the first result type based on action.Result
	if fn.Type != nil && fn.Type.Results != nil && len(fn.Type.Results.List) >= 1 {
		res := fn.Type.Results.List[0]
		if action.Result != "" {
			// Determine if action.Result should be a pointer type
			resultIsPointer := len(action.Result) > 0 && action.Result[0] == '*'
			resultName := action.Result
			if resultIsPointer {
				resultName = action.Result[1:] // Remove the '*' prefix
			}

			// Handle current *pkg.Type case
			if star, ok := res.Type.(*ast.StarExpr); ok {
				if sel, ok := star.X.(*ast.SelectorExpr); ok {
					if pkgIdent, ok := sel.X.(*ast.Ident); ok {
						if resultIsPointer {
							// Keep as pointer type, just update the name
							if sel.Sel.Name != resultName {
								changed = true
								newIdent := ast.NewIdent(resultName)
								newIdent.NamePos = sel.Sel.NamePos
								sel.Sel = newIdent
							}
						} else {
							// Convert from pointer to non-pointer type
							changed = true
							newSel := &ast.SelectorExpr{
								X:   pkgIdent,
								Sel: ast.NewIdent(resultName),
							}
							res.Type = newSel
						}
					}
				}
				// Handle current pkg.Type case
			} else if sel, ok := res.Type.(*ast.SelectorExpr); ok {
				if pkgIdent, ok := sel.X.(*ast.Ident); ok {
					if resultIsPointer {
						// Convert from non-pointer to pointer type
						changed = true
						newStar := &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   pkgIdent,
								Sel: ast.NewIdent(resultName),
							},
						}
						res.Type = newStar
					} else {
						// Keep as non-pointer type, just update the name
						if sel.Sel.Name != resultName {
							changed = true
							newIdent := ast.NewIdent(resultName)
							newIdent.NamePos = sel.Sel.NamePos
							sel.Sel = newIdent
						}
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
// or:           type user struct { service.Base[*model.User, model.UserReq, model.UserRsp] }
// depending on whether action.Payload/Result starts with '*'
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
						// Handle second parameter (Payload)
						if changed2 := applyServiceTypeParam(indexListExpr, 1, action.Payload); changed2 {
							changed = true
						}
						// Handle third parameter (Result)
						if changed3 := applyServiceTypeParam(indexListExpr, 2, action.Result); changed3 {
							changed = true
						}
					}
				}
			}
		}
	}

	return changed
}

// applyServiceTypeParam updates a specific type parameter in service.Base[T1, T2, T3]
// based on whether the actionType starts with '*' (pointer) or not (non-pointer)
func applyServiceTypeParam(indexListExpr *ast.IndexListExpr, paramIndex int, actionType string) bool {
	if paramIndex >= len(indexListExpr.Indices) {
		return false
	}

	// Determine if actionType should be a pointer type
	actionIsPointer := len(actionType) > 0 && actionType[0] == '*'
	actionName := actionType
	if actionIsPointer {
		actionName = actionType[1:] // Remove the '*' prefix
	}

	currentParam := indexListExpr.Indices[paramIndex]

	// Handle current *pkg.Type case
	if star, ok := currentParam.(*ast.StarExpr); ok {
		if sel, ok := star.X.(*ast.SelectorExpr); ok {
			if actionIsPointer {
				// Keep as pointer type, just update the name
				if sel.Sel.Name != actionName {
					newIdent := ast.NewIdent(actionName)
					newIdent.NamePos = sel.Sel.NamePos
					sel.Sel = newIdent
					return true
				}
			} else {
				// Convert from pointer to non-pointer type
				newIdent := ast.NewIdent(actionName)
				newIdent.NamePos = sel.Sel.NamePos
				sel.Sel = newIdent
				// Replace the StarExpr with SelectorExpr
				indexListExpr.Indices[paramIndex] = sel
				return true
			}
		}
	}

	// Handle current pkg.Type case (non-pointer)
	if sel, ok := currentParam.(*ast.SelectorExpr); ok {
		if actionIsPointer {
			// Convert from non-pointer to pointer type
			newIdent := ast.NewIdent(actionName)
			newIdent.NamePos = sel.Sel.NamePos
			// Create a new SelectorExpr with updated name
			newSel := &ast.SelectorExpr{
				X:   sel.X, // Keep the same package identifier
				Sel: newIdent,
			}
			// Wrap new SelectorExpr with StarExpr
			starExpr := &ast.StarExpr{
				Star: sel.Pos() - 1, // Position the * just before the selector
				X:    newSel,
			}
			indexListExpr.Indices[paramIndex] = starExpr
			return true
		} else {
			// Keep as non-pointer type, just update the name
			if sel.Sel.Name != actionName {
				newIdent := ast.NewIdent(actionName)
				newIdent.NamePos = sel.Sel.NamePos
				sel.Sel = newIdent
				return true
			}
		}
	}

	return false
}
