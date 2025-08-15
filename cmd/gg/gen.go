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

	if !fileExists(modelDir) {
		fmt.Fprintf(os.Stderr, "Error: model dir not found: %s\n", modelDir)
		os.Exit(1)
	}

	allModels, err := codegen.FindModels(module, modelDir, serviceDir, excludes)
	checkErr(err)
	if len(allModels) == 0 {
		return
	}

	modelStmts := make([]ast.Stmt, 0)
	serviceStmts := make([]ast.Stmt, 0)
	routerStmts := make([]ast.Stmt, 0)
	modelImportMap := make(map[string]struct{})
	routerImportMap := make(map[string]struct{})
	serviceImportMap := make(map[string]struct{})
	for _, m := range allModels {
		if m.Design.Enabled && m.Design.Migrate {
			// If the ModelFileDir is "model" or "model/", the model package name is the same as the model name,
			// and the statement in model/model.go will be "Register[*Project]()".
			// otherwise, the model package name is the last segment of the model file dir.
			//
			// For example:
			// If the ModelFileDir is "model/setting", the model package name is "setting",
			// then the statement in model/model.go should be "Register[*setting.Project]()"
			if m.ModelPkgName == strings.TrimRight(m.ModelFileDir, "/") {
				modelStmts = append(modelStmts, gen.StmtModelRegister(m.ModelName))
			} else {
				pkgName := strings.Split(m.ModelFileDir, "/")[1]
				modelStmts = append(modelStmts, gen.StmtModelRegister(fmt.Sprintf("%s.%s", pkgName, m.ModelName)))
			}

			if path, shouldImport := m.ModelImportPath(); shouldImport {
				modelImportMap[path] = struct{}{}
			}

		}

		m.Design.Range(func(s string, a *dsl.Action, p consts.Phase) {
			if a.Service {
				serviceImportMap[m.ServiceImportPath(modelDir, serviceDir)] = struct{}{}
			}
			routerImportMap[m.RouterImportPath()] = struct{}{}
		})
	}

	serviceAliasMap := gen.ResolveImportConflicts(lo.Keys(serviceImportMap))
	for _, m := range allModels {
		m.Design.Range(func(s string, a *dsl.Action, p consts.Phase) {
			if a.Service {
				if alias := serviceAliasMap[m.ServiceImportPath(modelDir, serviceDir)]; len(alias) > 0 {
					// alias import pacakge, eg:
					// pkg1_user "service/pkg1/user"
					// pkg2_user "service/pkg2/user"
					serviceStmts = append(serviceStmts, gen.StmtServiceRegister(fmt.Sprintf("%s.%s", alias, p.RoleName()), p))
				} else {
					serviceStmts = append(serviceStmts, gen.StmtServiceRegister(fmt.Sprintf("%s.%s", strings.ToLower(m.ModelName), p.RoleName()), p))
				}
			}
			routerStmts = append(routerStmts, gen.StmtRouterRegister(m.ModelPkgName, m.ModelName, a.Payload, a.Result, s, p.MethodName()))
		})
	}

	modelCode, err := gen.BuildModelFile("model", lo.Keys(modelImportMap), modelStmts...)
	checkErr(err)
	writeFileWithLog(filepath.Join(modelDir, "model.go"), modelCode)

	// generate service/service.go
	serviceCode, err := gen.BuildServiceFile("service", lo.Keys(serviceImportMap), serviceStmts...)
	checkErr(err)
	writeFileWithLog(filepath.Join(serviceDir, "service.go"), serviceCode)

	// generate router/router.go
	routerCode, err := gen.BuildRouterFile("router", lo.Keys(routerImportMap), routerStmts...)
	checkErr(err)
	writeFileWithLog(filepath.Join(routerDir, "router.go"), routerCode)

	// generate main.go
	mainCode, err := gen.BuildMainFile(module)
	checkErr(err)
	writeFileWithLog("main.go", mainCode)

	fset := token.NewFileSet()
	applyFile := func(filename string, code string, action *dsl.Action) {
		if fileExists(filename) {
			// Read original file content to preserve comments and formatting
			src, err := os.ReadFile(filename)
			checkErr(err)
			f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
			checkErr(err)

			// Apply changes
			changed := gen.ApplyServiceFile(f, action)

			if changed {
				// Only reformat and write file when there are changes
				// Use original FileSet to preserve comment positions
				code, err = gen.FormatNodeExtraWithFileSet(f, fset)
				checkErr(err)
				logUpdate(filename)
				checkErr(ensureParentDir(filename))
				checkErr(os.WriteFile(filename, []byte(code), 0o644))
			} else {
				logSkip(filename)
			}
		} else {
			logCreate(filename)
			checkErr(ensureParentDir(filename))
			checkErr(os.WriteFile(filename, []byte(code), 0o644))
		}
	}

	for _, m := range allModels {
		m.Design.Range(func(s string, a *dsl.Action, p consts.Phase) {
			if file := gen.GenerateService(m, a, p); file != nil {
				code, err := gen.FormatNodeExtra(file)
				checkErr(err)
				// code = gen.MethodAddComments(code, m.ModelName)
				dir := strings.Replace(m.ModelFilePath, modelDir, serviceDir, 1)
				dir = strings.TrimRight(dir, ".go")
				filename := filepath.Join(dir, strings.ToLower(string(p))+".go")
				applyFile(filename, code, a)
			}
		})
	}
}
