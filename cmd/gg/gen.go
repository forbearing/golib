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

	modelStmts := make([]ast.Stmt, 0)
	serviceStmts := make([]ast.Stmt, 0)
	routerStmts := make([]ast.Stmt, 0)
	modelImports := make(map[string]struct{})
	routerImports := make(map[string]struct{})
	sersviceImports := make(map[string]struct{})
	for _, m := range allModels {
		if m.Design.Enabled {
			modelStmts = append(modelStmts, gen.StmtModelRegister(m.ModelName))
		}
		dsl.RangeAction(m.Design, func(s string, a *dsl.Action, p consts.Phase) {
			serviceStmts = append(serviceStmts, gen.StmtServiceRegister(fmt.Sprintf("%s.%s", strings.ToLower(m.ModelName), p.RoleName()), p))
			routerStmts = append(routerStmts, gen.StmtRouterRegister(m.ModelPkgName, m.ModelName, a.Payload, a.Result, s, p.MethodName()))
			modelImport := filepath.Join(m.ModulePath, m.ModelPkgName)
			if !strings.HasSuffix(modelImport, "/model") {
				modelImports[modelImport] = struct{}{}
			}
			routerImports[filepath.Join(m.ModulePath, m.ModelPkgName)] = struct{}{}
			sersviceImports[filepath.Join(m.ModulePath, serviceDir, strings.ToLower(m.ModelName))] = struct{}{}
		})
	}
	modelCode, err := gen.BuildModelFile("model", lo.Keys(modelImports), modelStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(modelDir, "model.go"), []byte(modelCode), 0o644))

	// generate service/service.go
	serviceCode, err := gen.BuildServiceFile("service", lo.Keys(sersviceImports), nil, serviceStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(serviceDir, "service.go"), []byte(serviceCode), 0o644))

	// generate router/router.go
	routerCode, err := gen.BuildRouterFile("router", lo.Keys(routerImports), routerStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(routerDir, "router.go"), []byte(routerCode), 0o644))

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
			checkErr(ensureParentDir(filename))
			checkErr(os.WriteFile(filename, []byte(code), 0o644))
		} else {
			fmt.Printf("generate %s\n", filename)
			checkErr(ensureParentDir(filename))
			checkErr(os.WriteFile(filename, []byte(code), 0o644))
		}
	}

	for _, m := range allModels {
		dsl.RangeAction(m.Design, func(s string, a *dsl.Action, p consts.Phase) {
			if file := gen.GenerateService(m, a, p); file != nil {
				code, err := gen.FormatNodeExtra(file)
				checkErr(err)
				// code = gen.MethodAddComments(code, m.ModelName)
				filename := filepath.Join(serviceDir, strings.ToLower(m.ModelName), strings.ToLower(string(p))+".go")
				applyFile(filename, code, a)
			}
		})
	}
}
