# GoLib - Golang Lightning Backend Framework

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/forbearing/golib)

🚀 **GoLib** 是一个高性能、模块化的 Golang 后端开发框架，专为快速构建企业级 Web 应用而设计。

⚠️ **注意**: 框架正在积极开发中，API 可能会发生重大变化。

## ✨ 核心特性

- 🔥 **代码生成器**: 基于 DSL 自动生成 CRUD 操作代码
- 🏗️ **模块化架构**: 清晰的分层设计（Model-Service-Controller）
- 🚀 **高性能**: 基于 Gin 框架，支持高并发
- 🔧 **丰富组件**: 内置缓存、数据库、认证授权等组件
- 📊 **数据结构**: 提供完整的数据结构库
- 🔐 **安全特性**: 内置 JWT 认证和 RBAC 权限控制
- 📈 **可观测性**: 集成日志、指标和调试工具

## 🚀 快速开始

### 安装代码生成器

```bash
go install github.com/forbearing/golib/cmd/gg@latest
```

### 创建新项目

```bash
# 初始化项目
gg new demo
cd demo

# 新增 model
cat > model/user.go <<'EOF'
package model

import (
  . "github.com/forbearing/golib/dsl"
  "github.com/forbearing/golib/model"
)

type User struct {
  Name string
  Age  string

  model.Base
}

func (User) Design() {
  Migrate(true)

  Create(func() {
    Enabled(true)
  })
  Delete(func() {
    Enabled(true)
  })
  Update(func() {
    Enabled(true)
  })
  Patch(func() {
    Enabled(true)
  })
  List(func() {
    Enabled(true)
  })
  Get(func() {
    Enabled(true)
  })
  CreateMany(func() {
    Enabled(true)
  })
  DeleteMany(func() {
    Enabled(true)
  })
  UpdateMany(func() {
    Enabled(true)
  })
  PatchMany(func() {
    Enabled(true)
  })
}
EOF

# 生成代码
gg gen

# 运行项目
go run .
```

### 基本使用

1. **定义模型**:

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    model.Base
}

// DSL 配置
func (User) Design() {
    Enabled(true)
    Endpoint("users")
    Migrate(true)
    
    Create(func() {
        Enabled(true)
        Service(true)
        Payload[CreateUserRequest]()
        Result[*User]()
    })
    
    List(func() { Enabled(true) })
    Get(func() { Enabled(true) })
    Update(func() { Enabled(true) })
    Delete(func() { Enabled(true) })
}
```

2. **生成代码**:

```bash
gg gen  # 自动生成 Service、Controller、Router 代码
```

3. **启动应用**:

```go
go run .
```



## 🏗️ 架构概览

GoLib 采用分层架构设计，提供清晰的职责分离：

```
┌─────────────────┐
│   HTTP Router   │  ← 路由层：处理 HTTP 请求路由
├─────────────────┤
│   Controller    │  ← 控制层：处理请求/响应，参数验证
├─────────────────┤
│    Service      │  ← 业务层：业务逻辑处理
├─────────────────┤
│    Database     │  ← 数据层：数据库操作抽象
├─────────────────┤
│     Model       │  ← 模型层：数据模型定义
└─────────────────┘
```

### 核心组件

- **DSL**: 领域特定语言，用于声明式 API 设计
- **Code Generator**: 基于 DSL 自动生成样板代码
- **Bootstrap**: 应用启动和初始化管理
- **Middleware**: 中间件支持（认证、日志、CORS 等）
- **Provider**: 第三方服务集成（Redis、MongoDB、Kafka 等）
- **Cache**: 多种缓存实现（内存、Redis、Memcached）
- **Auth**: 认证授权（JWT、RBAC）

## 📚 示例和文档

### 完整示例项目

- [Demo 项目](./examples/demo/) - 展示框架完整功能的示例项目

### 使用指南

- [模型定义](./examples/demo/model/) - 如何定义数据模型
- [服务层开发](./examples/demo/service/) - 业务逻辑实现
- [路由配置](./examples/demo/router/router.go) - API 路由设置自动生成
- [控制器使用](./controller/README.md) - 控制器开发指南
- [客户端使用](./client/client_test.go) - HTTP 客户端使用

### 组件文档

- **数据库**: SQLite, PostgreSQL, MySQL, ClickHouse, SQL Server
- **缓存**: Redis, Memcached, 内存缓存
- **消息队列**: Kafka, RocketMQ, NATS, MQTT
- **存储**: MinIO, MongoDB, Cassandra
- **搜索**: Elasticsearch
- **监控**: InfluxDB, 指标收集
- **调试**: pprof, statsviz, gops



## 📊 数据结构库

GoLib 提供了完整的数据结构实现，适用于各种算法和应用场景：

### 线性数据结构

| 类型 | 实现 | 特点 | 使用场景 |
|------|------|------|----------|
| **List** | [ArrayList](./ds/list/arraylist/) | 动态数组，随机访问 O(1) | 需要频繁随机访问 |
| | [LinkedList](./ds/list/linkedlist/) | 链表，插入删除 O(1) | 频繁插入删除操作 |
| | [SkipList](./ds/list/skiplist/) | 跳表，有序，查找 O(log n) | 有序数据，范围查询 |
| **Stack** | [ArrayStack](./ds/stack/arraystack/) | 数组实现的栈 | 函数调用，表达式求值 |
| | [LinkedStack](./ds/stack/linkedstack/) | 链表实现的栈 | 动态大小的栈 |
| **Queue** | [ArrayQueue](./ds/queue/arrayqueue/) | 数组实现的队列 | 固定大小的缓冲区 |
| | [LinkedQueue](./ds/queue/linkedqueue/) | 链表实现的队列 | 动态大小的队列 |
| | [PriorityQueue](./ds/queue/priorityqueue/) | 优先队列 | 任务调度，算法优化 |
| | [CircularBuffer](./ds/queue/circularbuffer/) | 环形缓冲区 | 流数据处理 |

### 树形数据结构

| 类型 | 特点 | 时间复杂度 | 使用场景 |
|------|------|------------|----------|
| [RedBlack Tree](./ds/tree/rbtree/) | 自平衡二叉搜索树 | O(log n) | 有序映射，数据库索引 |
| [AVL Tree](./ds/tree/avltree/) | 严格平衡二叉搜索树 | O(log n) | 查找密集型应用 |
| [Splay Tree](./ds/tree/splaytree/) | 自调整二叉搜索树 | 摊销 O(log n) | 局部性访问模式 |
| [Trie](./ds/tree/trie/) | 前缀树 | O(m) | 字符串匹配，自动补全 |

### 其他数据结构

| 类型 | 实现 | 特点 | 使用场景 |
|------|------|------|----------|
| **Heap** | [BinaryHeap](./ds/heap/binaryheap/) | 完全二叉树实现 | 优先队列，堆排序 |
| **Set** | [MapSet](./ds/mapset/) | 基于 Map 的集合 | 去重，集合运算 |
| **MultiMap** | [MultiMap](./ds/multimap/) | 一对多映射 | 分组数据，索引 |


## 🔧 核心接口

GoLib 定义了一套清晰的接口体系，支持依赖注入和模块化开发：

### 🚀 应用初始化

```go
// Initalizer - 应用组件初始化接口
type Initalizer interface {
    Init() error  // 初始化组件，在应用启动时调用
}
```

**使用场景**: 数据库连接、缓存初始化、第三方服务配置等。

### 📝 日志系统

```go
// StandardLogger - 标准日志接口
type StandardLogger interface {
    Debug(args ...any)
    Info(args ...any)
    Warn(args ...any)
    Error(args ...any)
    Fatal(args ...any)
    
    Debugf(format string, args ...any)
    Infof(format string, args ...any)
    Warnf(format string, args ...any)
    Errorf(format string, args ...any)
    Fatalf(format string, args ...any)
}

// StructuredLogger - 结构化日志接口
type StructuredLogger interface {
    Debugw(msg string, keysAndValues ...any)
    Infow(msg string, keysAndValues ...any)
    Warnw(msg string, keysAndValues ...any)
    Errorw(msg string, keysAndValues ...any)
    Fatalw(msg string, keysAndValues ...any)
}

// ZapLogger - Zap 专用日志接口
type ZapLogger interface {
    Debugz(msg string, fields ...zap.Field)
    Infoz(msg string, fields ...zap.Field)
    Warnz(msg string, fields ...zap.Field)
    Errorz(msg string, fields ...zap.Field)
    Fatalz(msg string, fields ...zap.Field)
}

// Logger - 统一日志接口，组合所有日志功能
type Logger interface {
    With(fields ...string) Logger
    WithObject(name string, obj zapcore.ObjectMarshaler) Logger
    WithArray(name string, arr zapcore.ArrayMarshaler) Logger
    
    // 上下文感知日志
    WithControllerContext(*ControllerContext, consts.Phase) Logger
    WithServiceContext(*ServiceContext, consts.Phase) Logger
    WithDatabaseContext(*DatabaseContext, consts.Phase) Logger
    
    StandardLogger    // 标准日志方法
    StructuredLogger  // 结构化日志方法
    ZapLogger        // Zap 专用方法
}
```

**实现**: `logger/zap` 包提供了完整的日志实现，支持多种输出格式和级别控制。

### 💾 数据库接口

```go
// Database - 泛型数据库操作接口
type Database[M Model] interface {
    // 基础 CRUD 操作
    Create(objs ...M) error
    Delete(objs ...M) error
    Update(objs ...M) error
    UpdateById(id string, key string, value any) error
    
    // 查询操作
    List(dest *[]M, cache ...*[]byte) error
    Get(dest M, id string, cache ...*[]byte) error
    First(dest M, cache ...*[]byte) error
    Last(dest M, cache ...*[]byte) error
    Take(dest M, cache ...*[]byte) error
    Count(*int64) error
    
    // 维护操作
    Cleanup() error
    Health() error
    
    DatabaseOption[M]  // 嵌入查询选项
}

// DatabaseOption - 数据库查询选项接口
type DatabaseOption[M Model] interface {
    WithDB(any) Database[M]                    // 指定数据库连接
    WithTable(name string) Database[M]         // 指定表名
    WithDebug() Database[M]                    // 启用调试模式
    WithQuery(query M, fuzzyMatch ...bool) Database[M]  // 条件查询
    WithQueryRaw(query any, args ...any) Database[M]    // 原生 SQL 查询
    WithTransaction(tx any) Database[M]        // 事务支持
    WithScope(page, size int) Database[M]      // 分页
    WithOrder(order string) Database[M]        // 排序
    WithCache(...bool) Database[M]             // 缓存控制
    // ... 更多选项
}
```

**支持的数据库**: SQLite, PostgreSQL, MySQL, ClickHouse, SQL Server

**实现**: `database/gorm` 包提供了基于 GORM 的完整实现，支持事务、连接池、读写分离等高级功能。

### 📋 数据模型接口

```go
// Model - 数据模型基础接口
type Model interface {
    // 基础字段访问
    GetTableName() string
    GetID() string
    SetID(id ...string)
    
    // 审计字段
    GetCreatedBy() string
    GetUpdatedBy() string
    GetCreatedAt() time.Time
    GetUpdatedAt() time.Time
    SetCreatedBy(string)
    SetUpdatedBy(string)
    SetCreatedAt(time.Time)
    SetUpdatedAt(time.Time)
    
    // 查询控制
    Expands() []string                    // 关联查询字段
    Excludes() map[string][]any          // 排除条件
    MarshalLogObject(zapcore.ObjectEncoder) error  // 日志序列化
    
    Hooker  // 生命周期钩子
}

// Hooker - 数据操作生命周期钩子
type Hooker interface {
    CreateBefore() error  // 创建前钩子
    CreateAfter() error   // 创建后钩子
    DeleteBefore() error  // 删除前钩子
    DeleteAfter() error   // 删除后钩子
    UpdateBefore() error  // 更新前钩子
    UpdateAfter() error   // 更新后钩子
    ListBefore() error    // 列表查询前钩子
    ListAfter() error     // 列表查询后钩子
    GetBefore() error     // 单条查询前钩子
    GetAfter() error      // 单条查询后钩子
}
```

**实现**: 所有模型必须嵌入 `model.Base` 结构体，它提供了接口的默认实现。

### 🔧 服务层接口

```go
// Service - 业务逻辑服务接口
type Service[M Model, REQ Request, RSP Response] interface {
    // 基础 CRUD 操作
    Create(*ServiceContext, REQ) (RSP, error)
    Delete(*ServiceContext, REQ) (RSP, error)
    Update(*ServiceContext, REQ) (RSP, error)
    Patch(*ServiceContext, REQ) (RSP, error)
    List(*ServiceContext, REQ) (RSP, error)
    Get(*ServiceContext, REQ) (RSP, error)

    // 批量操作
    CreateMany(*ServiceContext, REQ) (RSP, error)
    DeleteMany(*ServiceContext, REQ) (RSP, error)
    UpdateMany(*ServiceContext, REQ) (RSP, error)
    PatchMany(*ServiceContext, REQ) (RSP, error)

    // 单条记录生命周期钩子
    CreateBefore(*ServiceContext, M) error
    CreateAfter(*ServiceContext, M) error
    DeleteBefore(*ServiceContext, M) error
    DeleteAfter(*ServiceContext, M) error
    UpdateBefore(*ServiceContext, M) error
    UpdateAfter(*ServiceContext, M) error
    PatchBefore(*ServiceContext, M) error
    PatchAfter(*ServiceContext, M) error
    ListBefore(*ServiceContext, *[]M) error
    ListAfter(*ServiceContext, *[]M) error
    GetBefore(*ServiceContext, M) error
    GetAfter(*ServiceContext, M) error

    // 批量操作生命周期钩子
    CreateManyBefore(*ServiceContext, ...M) error
    CreateManyAfter(*ServiceContext, ...M) error
    DeleteManyBefore(*ServiceContext, ...M) error
    DeleteManyAfter(*ServiceContext, ...M) error
    UpdateManyBefore(*ServiceContext, ...M) error
    UpdateManyAfter(*ServiceContext, ...M) error
    PatchManyBefore(*ServiceContext, ...M) error
    PatchManyAfter(*ServiceContext, ...M) error

    // 数据导入导出
    Import(*ServiceContext, io.Reader) ([]M, error)
    Export(*ServiceContext, ...M) ([]byte, error)

    // 数据过滤
    Filter(*ServiceContext, M) M
    FilterRaw(*ServiceContext) string

    Logger  // 嵌入日志接口
}
```

**特点**: Service 层在 Database 层之上，提供业务逻辑、权限控制、数据验证、生命周期管理等功能。

### 🔐 权限控制接口

```go
// RBAC - 基于角色的访问控制接口
type RBAC interface {
    // 角色管理
    AddRole(name string) error      // 添加角色
    RemoveRole(name string) error   // 删除角色

    // 权限管理
    GrantPermission(role string, resource string, action string) error   // 授予权限
    RevokePermission(role string, resource string, action string) error  // 撤销权限

    // 用户角色分配
    AssignRole(subject string, role string) error    // 分配角色
    UnassignRole(subject string, role string) error  // 取消角色分配
}
```

**使用场景**: API 权限控制、资源访问管理、用户权限验证等。

### 💾 缓存接口

```go
// Cache - 泛型缓存接口
type Cache[T any] interface {
    Set(key string, values T, ttl time.Duration)  // 设置缓存项
    Get(key string) (T, bool)                     // 获取缓存项
    Peek(key string) (T, bool)                    // 查看缓存项（不影响 LRU 顺序）
    Delete(key string)                            // 删除缓存项
    Exists(key string) bool                       // 检查缓存项是否存在
    Len() int                                     // 获取缓存项数量
    Clear()                                       // 清空所有缓存
}
```

**实现**: 支持内存缓存、Redis 缓存、Memcached 等多种缓存后端。

### 🔍 搜索引擎接口

```go
// ESDocumenter - Elasticsearch 文档接口
type ESDocumenter interface {
    Document() map[string]any  // 转换为 ES 文档格式
    GetID() string            // 获取文档 ID
}
```

**使用场景**: 全文搜索、数据分析、日志聚合等。实现此接口的模型可以自动同步到 Elasticsearch。

### 🌐 HTTP 客户端接口

```go
// HTTPClient - HTTP 客户端接口
type HTTPClient interface {
    Get(url string, headers ...map[string]string) (*http.Response, error)
    Post(url string, body io.Reader, headers ...map[string]string) (*http.Response, error)
    Put(url string, body io.Reader, headers ...map[string]string) (*http.Response, error)
    Delete(url string, headers ...map[string]string) (*http.Response, error)
    Patch(url string, body io.Reader, headers ...map[string]string) (*http.Response, error)
}
```

**特点**: 提供统一的 HTTP 客户端抽象，支持重试、超时、中间件等功能。

## 🚀 代码生成器

GoLib 提供了强大的代码生成器 `gg`，基于 DSL 自动生成完整的 CRUD API。

### 安装代码生成器

```bash
# 安装 gg 命令
go install github.com/forbearing/golib/cmd/gg@latest

# 验证安装
gg version
```

### DSL 语法

在模型文件中使用 `//go:generate gg gen` 注释和 DSL 配置：

```go
//go:generate gg gen

package model

import (
    "github.com/forbearing/golib/dsl"
    "github.com/forbearing/golib/model"
)

// User 用户模型
type User struct {
    model.Base
    Name     string `json:"name" gorm:"column:name"`
    Email    string `json:"email" gorm:"column:email"`
    Password string `json:"password" gorm:"column:password"`
}

// DSL 配置
var UserDesign = dsl.Design{
    Enabled:  true,                    // 启用代码生成
    Endpoint: "/api/v1/users",        // API 端点
    Migrate:  true,                   // 启用数据库迁移
    
    // CRUD 操作配置
    Create: dsl.Action{
        Enabled: true,
        Service: true,  // 生成 Service 层
        Public:  false, // 需要认证
    },
    Update: dsl.Action{
        Enabled: true,
        Service: true,
        Payload: "UserUpdateRequest",  // 自定义请求结构
    },
    Delete: dsl.Action{Enabled: true, Service: true},
    List:   dsl.Action{Enabled: true, Service: true, Public: true},
    Get:    dsl.Action{Enabled: true, Service: true, Public: true},
}
```

### 生成代码

```bash
# 或者直接使用 gg 命令
gg gen
```

### 生成的文件结构

```
.
├── config.ini.example
├── configx
│   └── configx.go
├── cronjob
│   └── cronjob.go
├── dao
├── go.mod
├── go.sum
├── main.go
├── middleware
│   └── middleware.go
├── model
│   ├── config
│   │   ├── namespace
│   │   │   ├── app
│   │   │   │   ├── env
│   │   │   │   │   ├── file.go
│   │   │   │   │   └── item.go
│   │   │   │   ├── env.go
│   │   │   │   └── filetemplate.go
│   │   │   └── app.go
│   │   └── namespace.go
│   ├── iam
│   │   ├── group.go
│   │   └── user.go
│   ├── model.go
│   └── setting
│       ├── project.go
│       ├── region.go
│       ├── tenant.go
│       └── vendor.go
├── provider
├── router
│   └── router.go
└── service
    ├── config
    │   └── namespace
    │       └── app
    │           ├── env
    │           │   └── item
    │           │       └── list.go
    │           └── list.go
    └── service.go
```

### DSL 配置选项

| 字段 | 类型 | 说明 |
|------|------|------|
| `Enabled` | bool | 是否启用代码生成 |
| `Endpoint` | string | API 基础路径 |
| `Param` | string | 路径参数名 |
| `Migrate` | bool | 是否生成数据库迁移 |
| `Service` | bool | 是否生成 Service 层 |
| `Public` | bool | 是否为公开 API（无需认证）|
| `Payload` | string | 自定义请求结构体名称 |
| `Result` | string | 自定义响应结构体名称 |

### 支持的操作类型

- **Create**: 创建单条记录
- **Update**: 更新单条记录（完整更新）
- **Patch**: 部分更新单条记录
- **Delete**: 删除单条记录
- **List**: 分页查询列表
- **Get**: 获取单条记录
- **CreateMany**: 批量创建
- **UpdateMany**: 批量更新
- **DeleteMany**: 批量删除

## 🔧 配置管理

GoLib 使用分层配置系统，支持多种配置源和环境。

### 配置文件结构

```yaml
# config/config.yaml
server:
  host: "0.0.0.0"
  port: 8080
  mode: "prod"

database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "user"
  password: "password"
  database: "myapp"
  
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

logger:
  level: "info"
  format: "json"  # json, text
  output: "stdout" # stdout, file
```

### 环境变量覆盖

```bash
# 环境变量会自动覆盖配置文件
export SERVER_HOST=0.0.0.0
export SERVER_PORT=8080
export DATABASE_HOST=postgres.example.com
export REDIS_HOST=redis.example.com
```

### 配置加载



