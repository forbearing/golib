package codegen

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/forbearing/golib/internal/codegen/gen"
)

// FindModelsInDirectory finds all models in a directory
func FindModelsInDirectory(module, modelDir, serviceDir string, excludes []string) ([]*gen.ModelInfo, error) {
	allModels := make([]*gen.ModelInfo, 0)

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

		models, err := gen.FindModels(module, path)
		if err != nil {
			return nil
		}
		for _, m := range models {
			dir := filepath.Dir(path)
			svcDir := strings.Replace(dir, modelDir, serviceDir, 1)
			svcFile := filepath.Join(svcDir, strings.ToLower(m.ModelName)+".go")
			m.ServiceFilePath = svcFile
			m.ModelFilePath = path
			allModels = append(allModels, m)
		}

		return nil
	})

	return allModels, nil
}
