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
	pkgnew "github.com/forbearing/golib/internal/codegen/new"
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

var tsCmd = &cobra.Command{
	Use:   "ts",
	Short: "generate typescript interface code",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("not implemented")
	},
}

func init() {
	genCmd.AddCommand(tsCmd)
}

func genRun() {
	if len(module) == 0 {
		var err error
		module, err = gen.GetModulePath()
		checkErr(err)
	}

	// Ensure required files exist
	logSection("Ensure Required Files")
	pkgnew.EnsureFileExists()

	if !fileExists(modelDir) {
		logError(fmt.Sprintf("model dir not found: %s", modelDir))
		os.Exit(1)
	}

	// Scan all models
	logSection("Scan Models")
	allModels, err := codegen.FindModels(module, modelDir, serviceDir, excludes)
	checkErr(err)
	if len(allModels) == 0 {
		fmt.Println(gray("  No models found, nothing to do"))
		return
	}
	fmt.Printf("  %s %d models found\n", green("✔"), len(allModels))

	// Record old service files list (if prune option is enabled)
	var oldServiceFiles []string
	if prune {
		oldServiceFiles = scanExistingServiceFiles(serviceDir, allModels)
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

	// Resolve import conflicts
	serviceAliasMap := gen.ResolveImportConflicts(lo.Keys(serviceImportMap))
	for _, m := range allModels {
		m.Design.Range(func(edp string, action *dsl.Action, phase consts.Phase) {
			if action.Service {
				if alias := serviceAliasMap[m.ServiceImportPath(modelDir, serviceDir)]; len(alias) > 0 {
					// alias import pacakge, eg:
					// pkg1_user "service/pkg1/user"
					// pkg2_user "service/pkg2/user"
					serviceStmts = append(serviceStmts, gen.StmtServiceRegister(fmt.Sprintf("%s.%s", alias, phase.RoleName()), phase))
				} else {
					serviceStmts = append(serviceStmts, gen.StmtServiceRegister(fmt.Sprintf("%s.%s", strings.ToLower(m.ModelName), phase.RoleName()), phase))
				}
			}
			route := "Auth"
			if action.Public {
				route = "Pub"
			}
			// If the phase is matched, the model endpoint will append the param, eg:
			// Endpoint: tenant, param is ":tenant", new endpoint is "tenant/:tenant"
			// Endpoint: tenant, param is ":id", new endpoint is "tenant/:id"
			switch phase {
			case consts.PHASE_DELETE, consts.PHASE_UPDATE, consts.PHASE_PATCH, consts.PHASE_GET:
				if len(m.Design.Param) == 0 {
					edp = filepath.Join(edp, ":id") // empty param will append default ":id" to endpoint.
				} else {
					edp = filepath.Join(edp, m.Design.Param)
				}
			case consts.PHASE_CREATE_MANY, consts.PHASE_DELETE_MANY, consts.PHASE_UPDATE_MANY, consts.PHASE_PATCH_MANY:
				edp = filepath.Join(edp, "batch")
			case consts.PHASE_IMPORT:
				edp = filepath.Join(edp, "import")
			case consts.PHASE_EXPORT:
				edp = filepath.Join(edp, "export")

			}

			routerStmts = append(routerStmts, gen.StmtRouterRegister(m.ModelPkgName, m.ModelName, action.Payload, action.Result, route, edp, phase.MethodName()))
		})
	}

	// ============================================================
	// Generate model/service/router/main files
	// ============================================================
	logSection("Generate Files")

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

	// ============================================================
	// Apply actions to services
	// ============================================================
	logSection("Apply Actions To Services")

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
				fset := token.NewFileSet()
				code, err := gen.FormatNodeExtraWithFileSet(file, fset)
				// pretty.Println(file)
				checkErr(err)
				// code = gen.MethodAddComments(code, m.ModelName)
				dir := strings.Replace(m.ModelFilePath, modelDir, serviceDir, 1)
				dir = strings.TrimSuffix(dir, ".go")
				filename := filepath.Join(dir, strings.ToLower(string(p))+".go")
				applyFile(filename, code, a)
			}
		})
	}

	// ============================================================
	// Prune disabled service files
	// ============================================================
	if prune && len(oldServiceFiles) > 0 {
		pruneServiceFiles(oldServiceFiles, allModels)
	}

	// ============================================================
	// Completion message
	// ============================================================
	logSection("Done")
	fmt.Printf("\n%s Code generation completed successfully!\n", green("🎉"))
}

// scanExistingServiceFiles scans existing service files in the service directory
// Only includes files that match phase names (e.g., create.go, update.go, etc.)
func scanExistingServiceFiles(serviceDir string, allModels []*gen.ModelInfo) []string {
	var files []string

	// Check if service directory exists
	if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
		return files
	}

	// Get all valid phase names
	validPhases := map[string]bool{
		consts.PHASE_CREATE.Filename():      true,
		consts.PHASE_DELETE.Filename():      true,
		consts.PHASE_UPDATE.Filename():      true,
		consts.PHASE_PATCH.Filename():       true,
		consts.PHASE_LIST.Filename():        true,
		consts.PHASE_GET.Filename():         true,
		consts.PHASE_CREATE_MANY.Filename(): true,
		consts.PHASE_DELETE_MANY.Filename(): true,
		consts.PHASE_UPDATE_MANY.Filename(): true,
		consts.PHASE_PATCH_MANY.Filename():  true,
		consts.PHASE_IMPORT.Filename():      true,
		consts.PHASE_EXPORT.Filename():      true,
	}

	// Walk through the service directory
	err := filepath.Walk(serviceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			fileName := filepath.Base(path)
			// Only include files that match phase names
			if validPhases[fileName] {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Warning: failed to scan existing service files: %v\n", err)
	}
	return files
}

// pruneServiceFiles prunes disabled service files
func pruneServiceFiles(oldServiceFiles []string, allModels []*gen.ModelInfo) {
	// Get list of service files that should currently exist
	currentServiceFiles := make(map[string]bool)
	for _, m := range allModels {
		m.Design.Range(func(s string, a *dsl.Action, p consts.Phase) {
			if a.Enabled && a.Service {
				dir := strings.Replace(m.ModelFilePath, modelDir, serviceDir, 1)
				dir = strings.TrimSuffix(dir, ".go")
				filename := filepath.Join(dir, strings.ToLower(string(p))+".go")
				currentServiceFiles[filename] = true
			}
		})
	}

	// Find files to delete (exist in old list but not in current list)
	filesToDelete := make([]string, 0)
	for _, oldFile := range oldServiceFiles {
		if !currentServiceFiles[oldFile] {
			filesToDelete = append(filesToDelete, oldFile)
		}
	}

	if len(filesToDelete) == 0 {
		fmt.Printf("  %s No disabled service files to prune\n", green("✔"))
		return
	}

	// Display list of files to be deleted
	fmt.Printf("\n%s Files to be deleted:\n", yellow("⚠"))
	for _, file := range filesToDelete {
		fmt.Printf("  %s %s\n", red("✘"), file)
	}

	// Ask user for confirmation
	fmt.Printf("\n%s Do you want to delete these files? (y/N): ", cyan("?"))
	var response string
	fmt.Scanln(&response)

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Printf("  %s Deletion cancelled\n", gray("→"))
		return
	}

	// Execute deletion operation
	for _, file := range filesToDelete {
		if err := os.Remove(file); err != nil {
			fmt.Printf("  %s Failed to delete %s: %v\n", red("✘"), file, err)
		} else {
			fmt.Printf("  %s Deleted %s\n", green("✔"), file)
		}
	}
}
