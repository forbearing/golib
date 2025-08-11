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
					if IsServiceType(typeSpec) {
						ApplyServiceType(typeSpec, action)
					}
				}
			}
		}
		if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl != nil {
			if IsServiceMethod1(funcDecl) {
				ApplyServiceMethod1(funcDecl, action)
			}
			if IsServiceMethod2(funcDecl) {
				ApplyServiceMethod2(funcDecl, action)
			}
			if IsServiceMethod3(funcDecl) {
				ApplyServiceMethod3(funcDecl, action)
			}
			if IsServiceMethod4(funcDecl) {
				ApplyServiceMethod4(funcDecl, action)
			}
		}
	}
}
