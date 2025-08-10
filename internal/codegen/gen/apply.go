package gen

import (
	"go/ast"

	"github.com/forbearing/golib/dsl"
)

func ApplyServiceFile(file *ast.File, action *dsl.Action) {
	if file == nil || action == nil {
		return
	}

	for _, decl := range file.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			ApplyServiceMethod1(funcDecl, action)
			ApplyServiceMethod2(funcDecl, action)
			ApplyServiceMethod3(funcDecl, action)
			ApplyServiceMethod4(funcDecl, action)
		}
	}
}
