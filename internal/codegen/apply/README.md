# Apply Package

The `apply` package provides functionality for generating service files from model definitions while preserving existing business logic using Go AST (Abstract Syntax Tree) analysis.

## 工作原理 (Working Principle)

### 1. Model 识别 (Model Recognition)
- 扫描 `model` 目录下的所有 Go 文件
- 识别符合条件的 model 结构体
- 根据 model 包名和结构体名生成对应的 service 文件路径

### 2. Service 文件处理 (Service File Processing)
- **新文件创建**: 如果对应的 service 文件不存在，直接使用 `gen` 包生成完整的 service 文件
- **现有文件更新**: 如果 service 文件已存在，使用 AST 分析保护业务逻辑

### 3. AST 业务逻辑保护 (AST Business Logic Protection)
通过 Go AST 分析实现精确的业务逻辑保护：

1. **识别服务结构体**: 查找继承了 `service.Base[*ModelName]` 的结构体
2. **方法检查**: 检查结构体是否包含所有必需的钩子函数方法
3. **智能补全**: 只添加缺失的钩子函数，不修改现有的业务逻辑代码
4. **代码保护**: 现有方法的业务逻辑代码完全不受影响

### 4. 钩子函数 (Hook Functions)
自动确保以下钩子函数存在：
- `CreateBefore/CreateAfter`
- `UpdateBefore/UpdateAfter`
- `DeleteBefore/DeleteAfter`
- `GetBefore/GetAfter`
- `ListBefore/ListAfter`
- `BatchCreateBefore/BatchCreateAfter`
- `BatchUpdateBefore/BatchUpdateAfter`
- `BatchDeleteBefore/BatchDeleteAfter`

## Features

- **Model Discovery**: Automatically discovers model definitions from specified directories
- **Service Generation**: Generates complete service files with CRUD operations
- **Business Logic Preservation**: Uses AST analysis to preserve existing business logic when regenerating service files
- **Configurable Exclusions**: Allows excluding specific models from generation
- **Path Management**: Handles proper file paths and package naming conventions
- **Smart Updates**: Only adds missing hook methods without touching existing business logic

## Core Components

### Types

- **`ApplyConfig`**: Configuration for service generation operations
- **`ServiceFileInfo`**: Information about existing service files including AST data
- **`BusinessLogicSection`**: Represents sections of business logic code

### Functions

- **`ApplyServiceGeneration`**: Main entry point for applying service generation
- **`generateNewServiceFile`**: Creates new service files from scratch
- **`applyServiceChanges`**: Updates existing service files while preserving business logic
- **`findServiceStruct`**: Locates service structs that inherit from service.Base
- **`hasMethod`**: Checks if a struct has specific methods
- **`generateHookMethod`**: Generates missing hook methods

## Usage Example

```go
package main

import (
    "log"
    "github.com/forbearing/golib/internal/codegen/apply"
)

func main() {
    // Create configuration
    config := apply.NewApplyConfig(
        "example.com/myapp",  // module path
        "./model",            // model directory
        "./service",          // service directory
    ).WithExclusions([]string{"internal_model"}) // exclude specific models

    // Generate services
    if err := apply.ApplyServiceGeneration(config); err != nil {
        log.Fatal(err)
    }
}
```

## 与 Gen 包的区别 (Difference from Gen Package)

- **Gen 包**: 只生成新的 service 文件，如果文件已存在则跳过
- **Apply 包**: 智能更新现有 service 文件，保护业务逻辑的同时补全缺失的钩子函数

## 安全保证 (Safety Guarantees)

1. **业务逻辑不丢失**: 现有的业务逻辑代码永远不会被覆盖或删除
2. **方法签名一致**: 生成的钩子函数签名与 `service.Base` 完全一致
3. **包结构保持**: 维持原有的包结构和导入关系
4. **代码格式化**: 自动格式化生成的代码，保持代码风格一致