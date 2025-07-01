// cmd/gg/main.go
package main

import (
	"flag"
	"log"

	"github.com/forbearing/golib/internal/codegen"
	"github.com/forbearing/golib/internal/codegen/ast"
	"github.com/forbearing/golib/internal/codegen/generator"
)

func main() {
	var (
		modelFile   = flag.String("model", "", "model file path")
		outputDir   = flag.String("output", "./service", "output directory")
		packageName = flag.String("package", "service", "service package name")
		debug       = flag.Bool("debug", false, "enable debug output")
	)
	flag.Parse()

	if *modelFile == "" {
		log.Fatal("model file is required")
	}

	// 加载配置（从 go.mod 读取模块路径）
	config, err := codegen.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 应用命令行参数
	config.ServicePackage = *packageName
	config.OutputDir = *outputDir

	// 创建解析器和生成器
	parser := ast.NewParser(config.ModulePath)
	generator, err := generator.NewServiceGenerator()
	if err != nil {
		log.Fatalf("failed to create generator: %v", err)
	}

	// 解析 model 文件
	var models []*codegen.ModelInfo
	if *debug {
		models, err = parser.ParseModelFileDebug(*modelFile)
	} else {
		models, err = parser.ParseModelFile(*modelFile)
	}

	if err != nil {
		log.Fatalf("failed to parse model file: %v", err)
	}

	if len(models) == 0 {
		log.Printf("No models found in file: %s", *modelFile)
		return
	}

	// 生成代码
	for _, model := range models {
		if err := generator.GenerateToFile(model, config); err != nil {
			log.Printf("failed to generate code for model %s: %v", model.Name, err)
			continue
		}
	}
}
