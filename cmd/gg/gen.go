package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/forbearing/golib/internal/codegen/gen"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "generate service code",
	Run: func(cmd *cobra.Command, args []string) {
		if len(module) == 0 {
			var err error
			module, err = gen.GetModulePath()
			checkErr(err)
		}

		models := make([]*gen.ModelInfo, 0)
		filepath.Walk(modelDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			base := filepath.Base(path)
			if path != modelDir && (base == "vendor" || base == "testdata") {
				return filepath.SkipDir
			}
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), ".go") ||
				strings.HasSuffix(info.Name(), "_test.go") ||
				strings.HasPrefix(info.Name(), "_") ||
				slices.Contains(excludes, info.Name()) {
				return nil
			}

			_models, err := gen.FindModels(module, path)
			if err != nil {
				return nil
			}
			for _, m := range _models {
				dir := filepath.Dir(path)
				svcDir := strings.Replace(dir, modelDir, serviceDir, 1)
				svcFile := filepath.Join(svcDir, strings.ToLower(m.ModelName)+".go")
				m.ServiceFilePath = svcFile
				m.ModelFilePath = path
				models = append(models, m)
			}

			return nil
		})

		for _, m := range models {
			dir := filepath.Dir(m.ServiceFilePath)
			checkErr(os.MkdirAll(dir, 0o755))

			file := gen.GenerateServiceFile(m)
			code, err := gen.FormatNode(file)
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
	},
}
