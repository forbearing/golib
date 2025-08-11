package gen

import (
	"go/ast"
	"go/token"

	"github.com/forbearing/golib/dsl"
)

func ApplyServiceFile(file *ast.File, action *dsl.Action) {
	if file == nil || action == nil {
		return
	}

	for _, decl := range file.Decls {
		// Handle type declarations first to adjust service.Base generics
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if isServiceType(typeSpec) {
						applyServiceType(typeSpec, action)
					}
				}
			}
		}
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl != nil {
			if isServiceMethod1(funcDecl) {
				applyServiceMethod1(funcDecl, action)
			}
			if isServiceMethod2(funcDecl) {
				applyServiceMethod2(funcDecl, action)
			}
			if isServiceMethod3(funcDecl) {
				applyServiceMethod3(funcDecl, action)
			}
			if isServiceMethod4(funcDecl) {
				applyServiceMethod4(funcDecl, action)
			}
		}
	}
}
