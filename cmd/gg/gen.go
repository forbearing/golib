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

	allModels, err := codegen.FindModelInfos(module, modelDir, serviceDir, excludes)
	checkErr(err)

	stmts := make([]ast.Stmt, 0)
	for _, m := range allModels {
		stmts = append(stmts, gen.StmtRouterRegister(m.PackageName, m.ModelName, m.ModelName, m.ModelName, strings.ToLower(m.ModelName)))
	}

	routerCode, err := gen.BuildRouterFile("router", stmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(routerDir, "router.go"), []byte(routerCode), 0o644))

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
