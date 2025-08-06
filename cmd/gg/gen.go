package main

import (
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	"github.com/forbearing/golib/internal/codegen"
	"github.com/forbearing/golib/internal/codegen/gen"
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
	for _, m := range allModels {
		routerStmts = append(routerStmts, gen.StmtRouterRegister(m.ModelPkgName, m.ModelName, m.ModelName, m.ModelName, m.Design.Endpoint))
		serviceStmts = append(serviceStmts, gen.StmtServiceRegister(strings.ToLower(m.ModelName)))
	}

	// generate router/router.go
	routerCode, err := gen.BuildRouterFile("router", routerStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(routerDir, "router.go"), []byte(routerCode), 0o644))

	// generate service/service.go
	serviceCode, err := gen.BuildServiceFile("service", serviceStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(serviceDir, "service.go"), []byte(serviceCode), 0o644))

	// generate main.go
	module, err := gen.GetModulePath()
	checkErr(err)
	mainCode, err := gen.BuildMainFile(module)
	checkErr(err)
	checkErr(os.WriteFile("main.go", []byte(mainCode), 0o644))

	for _, m := range allModels {
		dir := filepath.Dir(m.ServiceFilePath)
		checkErr(os.MkdirAll(dir, 0o755))

		file := gen.GenerateService(m)
		code, err := gen.FormatNodeExtra(file)
		checkErr(err)
		code = gen.MethodAddComments(code, m.ModelName)

		// If service file already exists, skip generate.
		if !fileExists(m.ServiceFilePath) {
			fmt.Printf("generate %s\n", m.ServiceFilePath)
			checkErr(os.WriteFile(m.ServiceFilePath, []byte(code), 0o644))
		} else {
			fmt.Printf("skip %s\n", m.ServiceFilePath)
		}

	}
}
