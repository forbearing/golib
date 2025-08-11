package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/dsl"
	"github.com/forbearing/golib/internal/codegen"
	"github.com/forbearing/golib/internal/codegen/gen"
	"github.com/forbearing/golib/types/consts"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate service code",
	Run: func(cmd *cobra.Command, args []string) {
		genRun()
	},
}

func genRun() {
	if len(module) == 0 {
		var err error
		module, err = gen.GetModulePath()
		checkErr(err)
	}

	allModels, err := codegen.FindModels(module, modelDir, serviceDir, excludes)
	checkErr(err)

	routerStmts := make([]ast.Stmt, 0)
	serviceStmts := make([]ast.Stmt, 0)
	importsMap := make(map[string]struct{})
	for _, m := range allModels {
		dsl.RangeAction(m.Design, func(s string, a *dsl.Action, p consts.Phase) {
			routerStmts = append(routerStmts, gen.StmtRouterRegister(m.ModelPkgName, m.ModelName, a.Payload, a.Result, s, p.MethodName()))
			serviceStmts = append(serviceStmts, gen.StmtServiceRegister(fmt.Sprintf("%s%s", strings.ToLower(m.ModelName), p.RoleName())))
			importsMap[fmt.Sprintf("%s/%s", m.ModelPkgName, m.ModelName)] = struct{}{}
		})
	}

	// generate router/router.go
	routerCode, err := gen.BuildRouterFile("router", routerStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(routerDir, "router.go"), []byte(routerCode), 0o644))

	// generate service/service.go
	serviceCode, err := gen.BuildServiceFile("service", lo.Keys(importsMap), nil, serviceStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(serviceDir, "service.go"), []byte(serviceCode), 0o644))

	// generate main.go
	mainCode, err := gen.BuildMainFile(module)
	checkErr(err)
	checkErr(os.WriteFile("main.go", []byte(mainCode), 0o644))

	fset := token.NewFileSet()
	applyFile := func(filename string, code string, action *dsl.Action) {
		// If service file already exists, skip generate.
		if fileExists(filename) {
			f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
			checkErr(err)
			gen.ApplyServiceFile(f, action)
			code, err = gen.FormatNodeExtra(f)
			checkErr(err)
			fmt.Printf("update %s\n", filename)
			checkErr(os.WriteFile(filename, []byte(code), 0o644))
		} else {
			fmt.Printf("generate %s\n", filename)
			checkErr(os.WriteFile(filename, []byte(code), 0o644))
		}
	}

	for _, m := range allModels {
		dir := filepath.Dir(m.ServiceFilePath)
		checkErr(os.MkdirAll(dir, 0o755))

		dsl.RangeAction(m.Design, func(s string, a *dsl.Action, p consts.Phase) {
			if file := gen.GenerateService(m, a, p); file != nil {
				code, err := gen.FormatNodeExtra(file)
				checkErr(err)
				// code = gen.MethodAddComments(code, m.ModelName)
				applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_"+strings.ToLower(string(p))+".go", code, a)
			}
		})

	}
}
