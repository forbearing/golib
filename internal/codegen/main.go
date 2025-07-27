package codegen

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/pflag"
)

var (
	modelDir   string
	serviceDir string
	excludes   []string
)

func init() {
	pflag.StringVarP(&modelDir, "model", "m", "model", "model directory path")
	pflag.StringVarP(&serviceDir, "service", "s", "service", "service directory path")
	pflag.StringSliceVarP(&excludes, "exclude", "e", nil, "exclude files")

	pflag.Parse()
}

func Main() {
	module, err := getModulePath()
	checkErr(err)

	models := make([]*ModelInfo, 0)

	filepath.Join()
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

		_models, err := findModels(module, path)
		if err != nil {
			return nil
		}
		for _, m := range _models {
			dir := filepath.Dir(path)
			svcDir := strings.Replace(dir, modelDir, serviceDir, 1)
			svcFile := filepath.Join(svcDir, strings.ToLower(m.ModelName)+".go")
			m.ServiceFilePath = svcFile
			models = append(models, m)
		}
		// for _, m := range _models {
		// 	fmt.Println(path, m.ServiceFilePath)
		// }

		return nil
	})

	for _, m := range models {
		dir := filepath.Dir(m.ServiceFilePath)
		checkErr(os.MkdirAll(dir, 0o755))

		file := generateServiceFile(m)
		code, err := formatNode(file)
		checkErr(err)
		code = methodAddComments(code, m.ModelName)

		fmt.Printf("Generate %s\n", m.ServiceFilePath)
		checkErr(os.WriteFile(m.ServiceFilePath, []byte(code), 0o644))

	}
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
