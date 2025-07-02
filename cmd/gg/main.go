// cmd/gg/main.go
package main

import (
	"fmt"
	"log"

	"github.com/forbearing/golib/internal/codegen"
	"github.com/forbearing/golib/internal/codegen/generator"
	"github.com/forbearing/golib/internal/codegen/walker"
	"github.com/spf13/pflag"
)

func main() {
	var (
		modelDir   = pflag.String("model", "./model", "model directory path")
		serviceDir = pflag.String("service", "./service", "service output directory")
		force      = pflag.BoolP("force", "f", false, "force overwrite existing service files")
		debug      = pflag.BoolP("debug", "d", false, "enable debug output")
	)
	pflag.Parse()

	// 加载配置（从 go.mod 读取模块路径）
	config, err := codegen.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 应用命令行参数
	config.ModelDir = *modelDir
	config.ServiceDir = *serviceDir
	config.ForceOverwrite = *force

	fmt.Printf("Scanning model directory: %s\n", config.ModelDir)
	fmt.Printf("Output service directory: %s\n", config.ServiceDir)
	fmt.Printf("Module path: %s\n", config.ModulePath)
	if *force {
		fmt.Printf("Force overwrite: enabled\n")
	} else {
		fmt.Printf("Force overwrite: disabled (existing files will be skipped)\n")
	}

	// 创建遍历器和生成器
	walker := walker.NewWalker(config)
	generator, err := generator.NewServiceGenerator(config)
	if err != nil {
		log.Fatalf("failed to create generator: %v", err)
	}

	// 遍历并解析所有 model 文件
	models, err := walker.WalkModels()
	if err != nil {
		log.Fatalf("failed to walk model directory: %v", err)
	}

	if len(models) == 0 {
		fmt.Printf("No models found in directory: %s\n", config.ModelDir)
		return
	}

	if *debug {
		fmt.Printf("\n=== Found %d models ===\n", len(models))
		for _, model := range models {
			packageName := config.GetServicePackageName(model.RelativePath)
			fmt.Printf("Model: %s\n", model.Name)
			fmt.Printf("  File: %s\n", model.FilePath)
			fmt.Printf("  Package: %s\n", model.PackageName)
			fmt.Printf("  Package Path: %s\n", model.PackagePath)
			fmt.Printf("  Relative Path: %s\n", model.RelativePath)
			fmt.Printf("  Service Package: %s\n", packageName)
			fmt.Println()
		}
	} else {
		fmt.Printf("Found %d models to process\n", len(models))
	}

	// 如果是强制覆盖模式，给出额外的确认信息
	if *force {
		fmt.Println("\n⚠️  WARNING: Force overwrite mode is enabled!")
		fmt.Println("   This will overwrite any existing service files.")
	}

	// 生成所有 service 代码
	fmt.Println("\nGenerating service files...")
	if err := generator.GenerateAll(models); err != nil {
		log.Fatalf("failed to generate service files: %v", err)
	}

	fmt.Println("Done!")
}
