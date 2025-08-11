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
	for _, m := range allModels {
		routerStmts = append(routerStmts, gen.StmtRouterRegister(m.ModelPkgName, m.ModelName, m.ModelName, m.ModelName, m.Design.Endpoint))
		serviceStmts = append(serviceStmts, gen.StmtServiceRegister(strings.ToLower(m.ModelName)))
	}

	// generate router/router.go
	routerCode, err := gen.BuildRouterFile("router", routerStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(routerDir, "router.go"), []byte(routerCode), 0o644))

	// generate service/service.go
	importsMap := make(map[string]struct{})
	// types := make([]*ast.GenDecl, 0)
	for _, m := range allModels {
		if m.Design.Create.Enabled {
			// types = append(types, gen.Types(m.ModelPkgName, m.ModelName, m.Design.Create.Payload, m.Design.Create.Result, false))
			importsMap[fmt.Sprintf("%s/%s", m.ModelPkgName, m.ModelName)] = struct{}{}
		}
	}
	// serviceCode, err := gen.BuildServiceFile("service", lo.Keys(importsMap), types, serviceStmts...)
	serviceCode, err := gen.BuildServiceFile("service", lo.Keys(importsMap), nil, serviceStmts...)
	checkErr(err)
	checkErr(os.WriteFile(filepath.Join(serviceDir, "service.go"), []byte(serviceCode), 0o644))

	// generate main.go
	mainCode, err := gen.BuildMainFile(module)
	checkErr(err)
	checkErr(os.WriteFile("main.go", []byte(mainCode), 0o644))

	// for _, m := range allModels {
	// 	pretty.Println(m)
	// }
	// return

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

		// service create
		if file := gen.GenerateService(m, m.Design.Create, consts.PHASE_CREATE); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_create.go", code, m.Design.Create)
		}

		// service delete
		if file := gen.GenerateService(m, m.Design.Delete, consts.PHASE_DELETE); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_delete.go", code, m.Design.Delete)
		}

		// service update
		if file := gen.GenerateService(m, m.Design.Update, consts.PHASE_UPDATE); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_update.go", code, m.Design.Update)
		}

		// service patch
		if file := gen.GenerateService(m, m.Design.Patch, consts.PHASE_PATCH); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_patch.go", code, m.Design.Patch)
		}

		// service list
		if file := gen.GenerateService(m, m.Design.List, consts.PHASE_LIST); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_list.go", code, m.Design.List)
		}

		// service get
		if file := gen.GenerateService(m, m.Design.Get, consts.PHASE_GET); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_get.go", code, m.Design.Get)
		}

		// service create many
		if file := gen.GenerateService(m, m.Design.CreateMany, consts.PHASE_CREATE_MANY); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_create_many.go", code, m.Design.CreateMany)
		}

		// service delete many
		if file := gen.GenerateService(m, m.Design.DeleteMany, consts.PHASE_DELETE_MANY); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_delete_many.go", code, m.Design.DeleteMany)
		}

		// service update many
		if file := gen.GenerateService(m, m.Design.UpdateMany, consts.PHASE_UPDATE_MANY); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_update_many.go", code, m.Design.UpdateMany)
		}

		// service patch many
		if file := gen.GenerateService(m, m.Design.PatchMany, consts.PHASE_PATCH_MANY); file != nil {
			code, err := gen.FormatNodeExtra(file)
			checkErr(err)
			// code = gen.MethodAddComments(code, m.ModelName)
			applyFile(strings.TrimRight(m.ServiceFilePath, ".go")+"_patch_many.go", code, m.Design.PatchMany)
		}

	}
}
