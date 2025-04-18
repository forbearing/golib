<a name="unreleased"></a>
## [Unreleased]

### Chore
- **deps:** upgrade dependencies to latest version

### Docs
- add CHANGELOG.md and .chglog configuration

### Refactor
- **boostrap:** replace custom initFunc with func() error


<a name="v0.4.2"></a>
## [v0.4.2] - 2025-04-01
### Chg
- **config:** 1.Remove `DB` field from `Server` config and move it to `Database` as field `Type` 2.change config Sqlite.IsMemory default value to true
- **etcd:** replace etcd default loggeer by logger.Etcd
- **pprof:** manually control mutex and block profile rate

### Chore
- update examples
- update example/simple
- update example/demo
- go mod tidy
- update example/myproject
- **deps:** upgrade dependencies to latest version
- **deps:** go mod tidy
- **deps:** add protoc too dependencies in go.mod
- **logger:** move time encoder format to consts package
- **nats:** replace zap.() with logger.nats

### Feat
- **logger:** expose zap.Logger instance via ZapLogger() method
- **logger:** add Clean function ot ensure proper zap logger shutdown
- **provider:** add package provider/rockeetmq to support rocketmq
- **provider:** add package rethinkdb
- **provider:** add package scylla to support scylladb
- **provider:** support memcached.

### Fix
- **config:** add custom ini encoder in latest viper version

### Refactor
- **scylla:** simplify batch statements appending


<a name="v0.4.1"></a>
## [v0.4.1] - 2025-03-12
### Enh
- **bootstrap:** improve application lifecycle management

### Feat
- **middleware:** add Circuit Breaker middleware

### Fix
- **gops:** prevent gops agent capture signal and exit 1

### Refactor
- **grpc:** improve server lifecycle management
- **pprof:** improve server liftcycle management
- **router:** improve server lifecycle management
- **statsviz:** improve server liftcycle management


<a name="v0.4.0"></a>
## [v0.4.0] - 2025-03-11
### Chg
- **bootstrap:** bootstrap feishu
- **bootstrap:** bootstrap influxdb
- **bootstrap:** bootstrap grpc server
- **bootstrap:** bootstrap kafka,nats,etcd,cassandra
- **config:** rename enable_statsviz -> statsviz_enable, enable_pprof -> pprof_enable, enable_gops -> gops_enable
- **config:** change redis config: remove `host`,`port`, add `addr`,`addrs`,`cluster_mode`
- **config:** remove redis config field: `idle_timeout`, `max_conn_age`
- **controller:** Remove operation in ExportFactory, ImportFactory
- **logger:** remove Global,Internal,Job, add Cassandra,Etcd,Feishu,Influxdb,Kafka,Ldap,Minio,Nats
- **redis:** upgrade Redis client from go-redis/v8 to go-redis/v9

### Chore
- go mod tidy
- **deps:** go mod tidy
- **deps:** upgrade dependencies to latest version
- **deps:** go mod tidy
- **deps:** upgrade dependencies to latest version
- **example:** update examples/myproject
- **examples:** replace github.com/pkg/errors with github.com/cockroachdb/errors
- **examples:** update demo using latest golib
- **examples:** update example simple
- **examples:** add example/bench
- **examples:** add examples/myproject
- **examples:** update example myproject
- **examples:** update example myproject
- **examples:** update example myproject
- **logger:** clean up comments and improve function naming
- **logger:** rename initVar to readConf for better clarify
- **logger:** rename logger Visitor -> Runtime for better clarity
- **minio:** remove debug print statement
- **minio:** rename cli -> client
- **redis:** no-op
- **redis:** remove unused comment

### Docs
- **config:** update comment for setDefault method

### Enh
- **config:** update influxdb configuration
- **elastic:** improve Init and New elastic client
- **ldap:** enhance provider ldap
- **minio:** enhance provider/minio
- **redis:** support redis cluster mode
- **redis:** enhance redis configuration and security options

### Feat
- prepare support grpc
- support BatchCreate, BatchDelete, BatchUpdate, BatchUpdatePartial
- support grpc server
- prepare support grpc
- **client:** add ListRaw and GetRaw methods
- **mongo:** enhance MongoDB client configuration and connection handing
- **provider:** support cassandra
- **provider:** support influxdb
- **provider:** support influxdb
- **provider:** prepare support `influxdb`, `feishu`
- **provider:** support nats
- **provider:** support kafka
- **provider:** prepare support cassandra, etcd, kafka, nats
- **provider:** support etcd
- **provider:** support feishu
- **redis:** graceful shutdown for connection cleanup
- **task:** improve collect runtime metrics
- **util:** add TLS configuration builder function `BuildTLSConfig`

### Fix
- **boostrap:** RegisterExitHandler(cassandra.Close)
- **controller:** use consts package for parameter constants
- **elasticsearch:** check elasticsearch connection in Init, - change config field: Hosts: string -> Addrs []string
- **influxdb:** properly close client on health check failure
- **kafka:** properly close client if no available broker
- **middleware:** replace logger.Global with zap.S()
- **mongo:** prevent potential use of invalid client on connection failure
- **mqtt:** prevent potential use of invalid client on connection failure
- **nats:** properly close client on health check failure
- **provider:** RegisterExitHandler(etcd.Close)
- **redis:**  close client on connection failure to prevent resource leaks

### Perf
- **boostrap:** run handlers concurrently to improve cleanup performance
- **redis:** replace encoding/json with json-iterator for better performance, add benchmark test case

### Refactor
- reorganize cache components into `cache` directory. - move database/redis, lru, cmap into `cache` directory
- reorganize components into `provider` directory. - move elastic, ldap, minio, mongo, mqtt, minio to `provider` directory
- reorganize components into provider directory - move elastic, ldap, minio, mongo, mqtt, minio to 'provider' directory
- **bootstrap:** rename exit handlers to cleanup handlers for clarity
- **config:** split config structs into seperate files
- **config:** modularize configuration defaults and move global constants near to configuration struct
- **config:** simplify config struct names and standardize viper usage


<a name="v0.3.4"></a>
## [v0.3.4] - 2025-03-05
### Chg
- **controller:** createSession -> CreateSession

### Chore
- ignore docs
- **example:** update examples/myproject
- **example:** update example/myproject
- **example:** update examples/demo
- **example:** update examples/simple
- **router:** refine final shutdown log message

### Docs
- update README.md

### Enh
- **config:** support read custom config values from envrionment variables, the priority is: env var > config file > default values

### Feat
- support debug tools: "pprof","gops"
- support debug/statsviz
- **bootstrap:** optimize cpu utilization with automaxprocs
- **bootstrap:** add Run to boostrap server

### Fix
- **config:** support parse default for time.Duration
- **config:** correct statsviz listen address field name
- **debug:** improve pprof shutdown handing
- **debug:** improve gops shutdown handing
- **debug:** improve statsviz shutdown handing

### Refactor
- **boostrap:** replace channel with errgroup for concurrent initialization

### Style
- **router:** standardize server log message format


<a name="v0.3.3"></a>
## [v0.3.3] - 2025-03-03
### Chore
- **example:** update examples/myproject
- **examples:** update example myproject
- **logger:** update comment

### Enh
- wrap errors with stack/context for better debugging.
- **task:** support reigster task before or after bootstrap.Boostrap()

### Feat
- add package cronjob
- add package cronjob
- add package cronjob
- **config:** support for custom config registration and retrieval

### Fix
- **boostrap:** prevent multiple initialization on repeated calls
- **controller:** correct error formatting in logs


<a name="v0.3.2"></a>
## [v0.3.2] - 2025-02-26
### Chg
- **config:** remove automatic domain assignment base on mode
- **config:** rename config.Auth: TokenExpireDuration -> AccessTokenExpireDuration; add RefreshTokenExpireDuration
- **config.Auth:** rename NoneExpireUser -> NoneExpireUsername; NoneExpirePass -> NoneExpirePassword
- **config.Auth:** rename NoneExpireUser -> NoneExpireUsername; NoneExpirePass -> NoneExpirePassword
- **logger:** improve GORM slow query logging with configuable threshold
- **middleware:** delete RequestID middleware, add TraceID middleware
- **response:** remove some Code and the response data add field `request_id`
- **router:** replace RequestID by TraceID

### Chore
- update example/myproject
- update examples/myproject
- nnop
- update examples/myproject
- **deps:** upgrade dependencies to latest version

### Docs
- update README.md
- update README.md

### Enh
- **jwt:** enhance jwt token handling

### Feat
- database support clickhouse
- database support sql server
- propagate tracing context to database layer and gorm with logging support
- **config:** add Clickhouse configuration support
- **config:** add SQL Server configuration support
- **config:** add DatabaseConfig to configures sqlite/postgres/mysql connection
- **config:** add slow_query_threshold configuration for server
- **util:** add TraceID and SpanID generation functions

### Fix
- **config:** set default value for Config to support read config from environment.

### Opt
- **logger:** optimize logger "With" performance

### Refactor
- **binaryheap:** remove redundant cmp parmmeter in downMinHeap and downMaxHeap methods
- **database:** simplify batch processing logic using min()
- **errors:** replace std "errors" with "github.com/cockroachdb/errors" for better error handing.
- **errors:** replace github.com/pkg/errors with github.com/cockroachdb/errors for better error handing
- **jwt:** remove accessTokenCache and refreshTokenCache

### Test
- **splaytree:** add debug print statement in test


<a name="v0.3.1"></a>
## [v0.3.1] - 2025-02-16
### Chg
- **avltree:** update WithNodeFormatter signature
- **avltree:** WithNodeFormat(string) -> WithNodeFormatter(func(*Node[K,V])string)
- **rbtree:** update WithNodeFormatter signature
- **rbtree:** WithNodeFormat(string) -> WithNodeFormatter(func(*Node[K,V])string)
- **trie:** update WithNodeFormatter and WithKeyFormatter signatures

### Chore
- update ds/tree/READMD.md
- **binaryheap:** remove comments
- **deps:** upgrade dependencies to latest version

### Docs
- **binaryheap:** fix function comments for heap operations
- **circularbuffer:** change "NewFromSlice" comments
- **rbtree:** add comments for Inorder and Postorder traversal methods
- **trie:** add comments for WithNodeFomatter,WithKeyFormatter

### Feat
- **arraylist:** add NewFromOrderedSlice, rename NewWithOrderedElements -> NewOrdered
- **arraylist:** add NewWithOrderedElements
- **avltree:** add String method for tree visualization
- **ds:** add avltree implement in package ds/tree/avltree
- **ds:** add binary heap implement on package ds/heap/binaryheap
- **ds:** add priority queue implementation on package ds/queue/priorityqueue
- **ds:** add splay tree implement in package ds/tree/splaytree
- **ds:** add trie implementation on package ds/tree/trie
- **ds:** add circular buffer implementation in package ds/queue/circularbuffer
- **ds:** add read-black tree implement in package ds/tree/rbtree
- **ds:** add skip list implementation in package ds/list/skiplist
- **rbtree:** rename Inorder -> InorerChan, add Inorder; Preorder,Postorder,LevelOrder same like Inorder
- **rbtree:** add GetNode to retrieve tree node by key
- **types:** add ErrFuncNil for nil function error handling

### Fix
- **avltree:** fix data race condition in String()
- **avltree:** add nil check for traversal functions
- **rbtree:** fix data race condition in String()
- **rbtree:** pass options to NewWithOrderedKeys in NewFromMapWithOrderedKeys
- **rbtree:** check comparsion function in New
- **rbtree:** initial default FackLocker
- **splaytree:** add nil check for traversal functions
- **trie:** fix data race condition in String()

### Refactor
- **arrayqueue:** use IsEmpty() instead of Len() == 0 for clarity
- **avltree:** change the avltree's method return type: *Node[K,V] -> (K,V,bool)
- **ds:** centralize error variables in ds/types/errors.go
- **multimap:** replace EqualFn with cmp function for value comparsion
- **priorityqueue:** simplify Clone function by removing redundant variable
- **rbtree:** change the rbtree's method return type: *Node[K,V] -> (K,V,bool)
- **rbtree:** reuse New in NewWithOrderedKeys
- **splaytree:** rename iter -> fn in MarshalJSON for clarity

### Style
- **avltree:** simplify AVL tree constructor name
- **rbtree:** simplify red-black tree constructor name
- **splaytree:** simplify splay tree constructor name

### Test
- **avltree:** optimize benchmark tests
- **avltree:** refactor compartor usage and add TestAVLTree_String
- **rbtree:** optimize benchmark tests
- **rbtree:** adjust benchmark size from {100, 100000} to {10, 100000}

### Tests
- **circularbuffer:** add test case for json encoding


<a name="v0.3.0"></a>
## [v0.3.0] - 2025-01-29
### Chore
- rename ds/multimap/multimap_bechmark_test.go -> ds/multimap/multimap_benchmark_test.go
- **deps:** upgrade dependencies to latest version
- **linkedlist:** update comment
- **linkedlist:** clarify MergeSorted doc
- **mapset:** move MarshalJSON and UnmarshalJSON to set_encoding.go

### Docs
- **arraystack:** fix typo in NewFromMapValues
- **mapset:** fix typo in UnmarshalJSON commit

### Feat
- **arraylist:** add options method to clone arraylist properties
- **arraylist:** support concurrent safe.
- **arraystack:** support concurrency safety
- **arraystack:** add a stack based on arraylist
- **ds:** add package ds/mapset that implement datastructre "set"
- **ds:** add a queue based on linkedlist in package ds/queue/linkdqueue
- **ds:** add a queue based on array list in package ds/queue/arrayqueue
- **ds:** add a stack based on linkedlist in ds/stack/linkdstack
- **ds:** add linkedlist package under ds/list
- **ds:** add arraylist package under ds/list
- **linkedlist:** add options method to clone linkedlist properties
- **linkedstack:** support concurrency safe
- **mapset:** provides WithSorted to support makeup sorted internal element, affect method: `Slice`,`String`,`MarshalJSON`, `Range`, `Iter`
- **mapset:** support concurrent safe

### Fix
- **arraylist:** ensure the underlying array capacity is always greater than 0
- **arraystack:** not use new array stack to avoid sync.RWMutex leak
- **linkedlist:** call internal "pushBackNode" to avoid deadlock in concurrent mode
- **types:** correct spelling of FakeLocker, FackeLocker -> FakeLocker

### Refactor
- **arraylist:** replace paramater "equal" with "cmp" in List[E any]
- **arraylist:** replace manual slice construction with s.list.Values()
- **arraylist:** rename parameters: values -> elements, v -> e
- **arrayqueue:** simplify Queue initialization in New function
- **arraystack:** rename `slices` parameter to `slice` in NewFromSlice
- **ds:** move ds interface and types from types to dedicated package ds/types
- **linkedlist:** replace manual slice with s.list.Slice()
- **linkedlist:** rename `slices` to `slice` in benchmark tests
- **linkedqueue:** use IsEmpty() instead of Len() == 0 for clarity
- **linkedstack:** rename `slices` parameter to `slice` in NewFromSlice
- **mapset:** rename `slices` parameter to `slice` in NewFromSlice
- **mapset:** rename mapset.go -> set.go; rename mapset_test.go -> set_test.go
- **mapset:** rename file: set.go -> mapset.go; rename package: set -> mapset

### Style
- **arraylist:** rename type parameter T -> E

### Test
- **arraylist:** refactor arraylist benchmark test case
- **linkedlist:** rename test case name
- **linkedlist:** refactor becnhark test units
- **linkedlist:** refactor and improve benchmark tests

### Tests
- **mapset:** add unit tests for mapset

### Pull Requests
- Merge pull request [#1](https://github.com/forbearing/golib/issues/1) from forbearing/feat/ds


<a name="v0.2.3"></a>
## [v0.2.3] - 2025-01-18
### Chore
- **deps:** upgrade dependencies to latest version
- **examples:** update example myproject
- **examples:** update example demo
- **examples:** update example simple
- **router:** disply exit signal in shutdown log

### Enh
- **router:** enhance server with graceful shutdown handling

### Feat
- **config:** add configurations constants for environment variables
- **logger:** support custom console encoder for better log formatting

### Refactor
- **logger:** remove 'log_' prefix from logger config fields


<a name="v0.2.2"></a>
## [v0.2.2] - 2025-01-07
### Refactor
- move context conversion functions to types/helper package


<a name="v0.2.1"></a>
## [v0.2.1] - 2025-01-07
### Chore
- **deps:** bump go packages
- **docs:** update READMD.md
- **docs:** update READMD.md
- **example:** update examples/demo
- **example:** update examples/simple
- **examples:** update myproject

### Feat
- **looger:** add WithControllerContext and WithServiceContext methods - WithControllerContext: for controller layer context fields - WithServiceControtext: for service layer context fields
- **model:** add tag "url" for query parameter that used by Client package

### Refactor
- **controller:** simplify logging context setup


<a name="v0.2.0"></a>
## [v0.2.0] - 2025-01-06
### Chore
- **deps:** bump go packages
- **deps:** upgrade dependencies to latest version
- **docs:** update READMD.md
- **docs:** update READMD.md
- **docs:** update READMD.md
- **docs:** update READMD.md
- **docs:** update READMD.md
- **examples:** update examples
- **examples:** update examples
- **examples:** remove unused example file
- **gitignore:** add *.db to ignore list
- **testdata:** add restart policy for docker-compose.yaml

### Enh
- **config:** auto create empty config file in temp directory
- **logger:** add Protocol and Binary logger
- **tunnel:** 1. add DecodePayload. 2. update testcase. 3. update docs.

### Feat
- add tunnel package
- **config:** support read config from environment and env has more priority than config file.
- **ldap:** add ldap authentication package
- **ldap:** add ldap authentication package
- **pkg:** add version package
- **tunnel:** add NewCmd, add more testcase
- **types:** add consts package
- **util:** add IsConnClosed
- **util:** add net utility functions

### Refactor
- **mongo:** rename makeURI -> buildURI
- **tunnel:** improve DecodePayload function signature
- **tunnel:** simplify CMD
- **tunnel:** cleanup code and repalce json to msgpack
- **types:** move constants into dedicated consts package


<a name="v0.1.1"></a>
## [v0.1.1] - 2024-12-27
### Chore
- update examples/simple
- update examples/demo
- update READMD.md
- **deps:** upgrade dependencies to latest version

### Docs
- **example:** update myproject code
- **example:** update demo code

### Feat
- **elastic:** add New function to create seperate elasticsearch client.
- **minio:** add `New` function to create minio seperate client.
- **mongo:** add `New` function to create seperate mongo client instance
- **mqtt:** add `New` function to create seperate client
- **mysql:** add `New` to create seperate instance.
- **postgres:** add `New` to create seperate instance.
- **router:** add function `RegisterWithConfig` to custom controller configuration
- **sqlite:** add `New` to create seperate instance.

### Refactor
- **mysql:** rename makeDSN to buildDSN
- **postgres:** rename makeDSN to buildDSN
- **sqlite:** rename makeDSN to buildDSN

### Style
- **mongo:** format log statement in a single line


<a name="v0.1.0"></a>
## [v0.1.0] - 2024-12-26
### Chore
- update examples/demo
- update README.md
- update READMD.md
- update READMD.md
- update examples/demo
- add READMD.md for controller
- update READMD.md
- update READMD.md
- update examples/demo
- update READMD.md
- update READMD.md
- update examples/demo
- update READMD.md
- update example/demo
- update examples/demo
- update examples/simple
- update README.md
- update README.md
- update examples/simple
- update README.md
- update example/demo
- bump go pkg version to latest
- **model:** add doc for `Register` and `Register`, deprecated `RegisterRoutes`

### Enh
- **model:** model.Register() will set id before insert table records


<a name="v0.0.66"></a>
## [v0.0.66] - 2024-12-18
### Chg
- **model:** remove RegisterRoutes
- **model:** rename GetTablename -> GetTableName
- **model:** move GormScannerWrapper from `model.go` to `util.go`, add function `GetTablename`
- **model:** remove model.Base field `Error`

### Chore
- update examples
- update README.md
- update README.md
- update README.md
- update exmaples/simple
- update README.md
- disinfect
- **database:** add more logger

### Feat
- **router:** add function `Register()` to quickly register routes

### Fix
- **model:** change tablename to SnakeCase

### Rename
- model.Verb -> types.HTTPVerb


<a name="v0.0.65"></a>
## [v0.0.65] - 2024-12-09
### Add
- **model:** `SysInfo`
- **util:** `Round` make float to specified percision.

### Chg
- change `model.CreatedAt`, `model.UpdatedAt` type to *time.Time

### Chore
- update examples
- update examples
- remove old comment
- use go1.23
- **controller:** use new filetype pkg path
- **database:** update testcase
- **model:** update doc
- **model:** rename: common.go -> user-agent.go
- **model:** remove field InternalMark

### Fix
- **controller:** using default database
- **model:** sysinfo

### Opt
- **util:** replace `satori/go.uuid` by `google/uuid`

### Rename
- pkg/http_wrapper -> pkg/httpwrapper
- sizedbufferpool -> pkg/sizedbufferpool
- filetype -> pkg/filetype
- sftp -> pkg/sftp
- bufferpool -> pkg/bufferpool
- cache/bigcache -> pkg/bigcache
- net/wrapper -> pkg/http_wrapper


<a name="v0.0.64"></a>
## [v0.0.64] - 2024-12-03
### Fix
- import `context`


<a name="v0.0.63"></a>
## [v0.0.63] - 2024-11-30
### Feat
- **database:** add `WithTryRun` to skip database operation but only invoke model layer hook.

### Fix
- **controller:** `GetFactory`: model set id to support invoke `GetBefore` hook


<a name="v0.0.62"></a>
## [v0.0.62] - 2024-11-26
### Chg
- **elastic:** param add context


<a name="v0.0.61"></a>
## [v0.0.61] - 2024-11-25
### Fix
- **databasee:** WithSelect


<a name="v0.0.60"></a>
## [v0.0.60] - 2024-11-15
### Fix
- **task:** error cause exit


<a name="v0.0.59"></a>
## [v0.0.59] - 2024-11-12
### Chore
- bump go package version
- **database:** remove debug
- **elastic:** add more case

### Fix
- **elasitc:** allow size to 0 to support DSL `aggs`


<a name="v0.0.58"></a>
## [v0.0.58] - 2024-11-11
### Chore
- bump go packages
- **elastic:** add more testcase

### Enh
- **elastic:** QueryBuilder support `aggs`
- **elastic:** Document.Search support `aggs`


<a name="v0.0.57"></a>
## [v0.0.57] - 2024-11-09
### Chore
- **elastic:** add more testcase

### Enh
- **elastic:** QueryBuilder add more doc, add method `BuildForce`
- **elastic:** add `SearchNext` to searches for N next hits, add `SearchPrev` to searchs for N previous hits.
- **elastic:** improve QueryBuilder to suport complex bool query


<a name="v0.0.56"></a>
## [v0.0.56] - 2024-11-07
### Feat
- **database:** add `WithJoinRaw`, `WithSelectRaw`
- **database:** add `WithJoinRaw`, `WithSelectRaw`


<a name="v0.0.55"></a>
## [v0.0.55] - 2024-11-07
### Bugfix
- Create/Update will remove/invalide cache, feat: trace database operation cost, feat: add `WithTransaction`, `WithLock`, `WithJoin`,`WithGroup`, `WithHaving`

### Chg
- interface `Database`, `DatebaseOption` Database: add `Health` DatabaseOption: add `WithTransaction`, `WithLock`

### Chore
- bump go package
- **logger:** -

### Enh
- **elastic:** support QueryBuilder


<a name="v0.0.54"></a>
## [v0.0.54] - 2024-11-04
### Add
- **logger:** mongo logger

### Chg
- **bootstrap:** bootstrap mongo
- **config:** mqtt config

### Chore
- update example
- bump go package
- update examples

### Enh
- **mqtt:** reimplement package mqtt

### Feat
- mongo package

### Update
- **config:** add mongo config


<a name="v0.0.53"></a>
## [v0.0.53] - 2024-11-02
### Chg
- **config:** config add `enable`
- **minio:** check `enable`
- **mqtt:** check `enable`
- **rbac:** check `enable`
- **util.RunOrDie:** error exit with context


<a name="v0.0.52"></a>
## [v0.0.52] - 2024-11-02
### Chg
- bootstrap mqtt
- **boostrap:** boostrap all database
- **config:** `server` config add `db` to specific which database should use
- **database:** database only boostrap when `server.db` is meet current database
- **example:** update

### Chore
- update README.md
- update READMD.md
- bump go package

### Opt
- **logger:** more check


<a name="v0.0.51"></a>
## [v0.0.51] - 2024-11-01
### Fix
- add recover for task


<a name="v0.0.50"></a>
## [v0.0.50] - 2024-10-24
### Chore
- update StringAny


<a name="v0.0.49"></a>
## [v0.0.49] - 2024-10-24
### Chg
- replace cmap to lru

### Chore
- bump go package
- remove comment


<a name="v0.0.48"></a>
## [v0.0.48] - 2024-10-22
### Chg
- set default log to console; set controller access log to access.log
- BulkIndex -> (*document).BulkIndex

### Chore
- update examples
- change logger position
- add more log
- update example
- update examples
- update README.md
- update README.md
- update README.md
- update README.md
- update README.md
- update README.md
- update README.md
- update README.md

### Fix
- support query parameter _select


<a name="v0.0.47"></a>
## [v0.0.47] - 2024-10-16
### Add
- Contains

### Chore
- update README.md
- update examples


<a name="v0.0.46"></a>
## [v0.0.46] - 2024-10-16
### Feat
- controller support `_select` query params

### Fix
- WithSelect


<a name="v0.0.45"></a>
## [v0.0.45] - 2024-10-13
### Feat
- support using custom index to query
- support using custom index to query


<a name="v0.0.44"></a>
## [v0.0.44] - 2024-10-13
### Feat
- database support WithSelect


<a name="v0.0.43"></a>
## [v0.0.43] - 2024-10-13
### Chore
- update README.md


<a name="v0.0.42"></a>
## [v0.0.42] - 2024-10-11
### Chg
- write `access_token`, `refresh_token`

### Chore
- bump go package version
- clean code

### Feat
- support refresh token; upgrade jwt to v5


<a name="v0.0.41"></a>
## [v0.0.41] - 2024-10-10
### Chg
- change tinyint -> smallint to support postgresql
- using helper
- remove router `/api/ping`

### Chore
- upgrade gorm drivers and plugins
- update documents
- update README.md
- bump go packages
- update examples for database postgresql
- change log
- change log
- update README.md
- update example

### Feat
- support database/postgresql
- database support sqlite


<a name="v0.0.40"></a>
## [v0.0.40] - 2024-10-09
### Chore
- update README.md


<a name="v0.0.39"></a>
## [v0.0.39] - 2024-10-07
### Chore
- using const
- add doc
- remove debug output

### Feat
- support _nototal in controller layer

### Fix
- nil Rows cause panic

### Opt
- add table index for `updated_at`, `created_by`,`updated_by`


<a name="v0.0.38"></a>
## [v0.0.38] - 2024-09-30

<a name="v0.0.37"></a>
## [v0.0.37] - 2024-09-30
### Opt
- concurrently query column data from database


<a name="v0.0.36"></a>
## [v0.0.36] - 2024-09-29

<a name="v0.0.35"></a>
## [v0.0.35] - 2024-09-29

<a name="v0.0.34"></a>
## [v0.0.34] - 2024-09-28
### Fix
- using new session for batch size.

### Task
- logger add cost field


<a name="v0.0.33"></a>
## [v0.0.33] - 2024-09-22
### Fix
- use logger.Task in task package


<a name="v0.0.32"></a>
## [v0.0.32] - 2024-09-04
### Chg
- remove default middleware.RateLimiter


<a name="v0.0.31"></a>
## [v0.0.31] - 2024-08-30

<a name="v0.0.30"></a>
## [v0.0.30] - 2024-08-25

<a name="v0.0.29"></a>
## [v0.0.29] - 2024-08-24
### Chg
- default base-auth and token using config


<a name="v0.0.28"></a>
## [v0.0.28] - 2024-08-24

<a name="v0.0.27"></a>
## [v0.0.27] - 2024-08-24

<a name="v0.0.26"></a>
## [v0.0.26] - 2024-08-23

<a name="v0.0.25"></a>
## [v0.0.25] - 2024-08-23

<a name="v0.0.24"></a>
## [v0.0.24] - 2024-08-22

<a name="v0.0.23"></a>
## [v0.0.23] - 2024-08-22

<a name="v0.0.22"></a>
## [v0.0.22] - 2024-08-02

<a name="v0.0.21"></a>
## [v0.0.21] - 2024-07-24

<a name="v0.0.20"></a>
## [v0.0.20] - 2024-06-28

<a name="v0.0.19"></a>
## [v0.0.19] - 2024-06-17
### Add
- SetConfigFile


<a name="v0.0.18"></a>
## [v0.0.18] - 2024-06-17
### Add
- Cache.Init

### Feat
- add SetConfigName, SetConfigType


<a name="v0.0.17"></a>
## [v0.0.17] - 2024-06-17

<a name="v0.0.16"></a>
## [v0.0.16] - 2024-05-25

<a name="v0.0.15"></a>
## [v0.0.15] - 2024-04-05
### Fix
- using default mysql instance.

### Opt
- upgrade boostrap package.


<a name="v0.0.14"></a>
## [v0.0.14] - 2024-04-03

<a name="v0.0.13"></a>
## [v0.0.13] - 2024-03-19
### Fix
- database.WithQuery
- util.Depointer


<a name="v0.0.12"></a>
## [v0.0.12] - 2024-03-04

<a name="v0.0.11"></a>
## [v0.0.11] - 2024-03-04

<a name="v0.0.10"></a>
## [v0.0.10] - 2024-03-04
### Fix
- register -> Register


<a name="v0.0.9"></a>
## [v0.0.9] - 2024-03-04
### Fix
- service.base[M types.Model] -> service.Base[M types.Model]


<a name="v0.0.8"></a>
## [v0.0.8] - 2024-03-02

<a name="v0.0.7"></a>
## [v0.0.7] - 2024-03-02

<a name="v0.0.6"></a>
## [v0.0.6] - 2024-03-02

<a name="v0.0.5"></a>
## [v0.0.5] - 2024-03-02
### Fix
- If structure field not contains json tags, structure lowercase field name as database query condition


<a name="v0.0.4"></a>
## [v0.0.4] - 2024-03-02
### Fix
- If structure field not contains json tags, structure lowercase field name as database query condition.


<a name="v0.0.3"></a>
## [v0.0.3] - 2024-02-21
### Fix
- disable automigrate model User


<a name="v0.0.2"></a>
## [v0.0.2] - 2024-02-16

<a name="v0.0.1"></a>
## v0.0.1 - 2024-02-15

[Unreleased]: https://github.com/forbearing/golib/compare/v0.4.2...HEAD
[v0.4.2]: https://github.com/forbearing/golib/compare/v0.4.1...v0.4.2
[v0.4.1]: https://github.com/forbearing/golib/compare/v0.4.0...v0.4.1
[v0.4.0]: https://github.com/forbearing/golib/compare/v0.3.4...v0.4.0
[v0.3.4]: https://github.com/forbearing/golib/compare/v0.3.3...v0.3.4
[v0.3.3]: https://github.com/forbearing/golib/compare/v0.3.2...v0.3.3
[v0.3.2]: https://github.com/forbearing/golib/compare/v0.3.1...v0.3.2
[v0.3.1]: https://github.com/forbearing/golib/compare/v0.3.0...v0.3.1
[v0.3.0]: https://github.com/forbearing/golib/compare/v0.2.3...v0.3.0
[v0.2.3]: https://github.com/forbearing/golib/compare/v0.2.2...v0.2.3
[v0.2.2]: https://github.com/forbearing/golib/compare/v0.2.1...v0.2.2
[v0.2.1]: https://github.com/forbearing/golib/compare/v0.2.0...v0.2.1
[v0.2.0]: https://github.com/forbearing/golib/compare/v0.1.1...v0.2.0
[v0.1.1]: https://github.com/forbearing/golib/compare/v0.1.0...v0.1.1
[v0.1.0]: https://github.com/forbearing/golib/compare/v0.0.66...v0.1.0
[v0.0.66]: https://github.com/forbearing/golib/compare/v0.0.65...v0.0.66
[v0.0.65]: https://github.com/forbearing/golib/compare/v0.0.64...v0.0.65
[v0.0.64]: https://github.com/forbearing/golib/compare/v0.0.63...v0.0.64
[v0.0.63]: https://github.com/forbearing/golib/compare/v0.0.62...v0.0.63
[v0.0.62]: https://github.com/forbearing/golib/compare/v0.0.61...v0.0.62
[v0.0.61]: https://github.com/forbearing/golib/compare/v0.0.60...v0.0.61
[v0.0.60]: https://github.com/forbearing/golib/compare/v0.0.59...v0.0.60
[v0.0.59]: https://github.com/forbearing/golib/compare/v0.0.58...v0.0.59
[v0.0.58]: https://github.com/forbearing/golib/compare/v0.0.57...v0.0.58
[v0.0.57]: https://github.com/forbearing/golib/compare/v0.0.56...v0.0.57
[v0.0.56]: https://github.com/forbearing/golib/compare/v0.0.55...v0.0.56
[v0.0.55]: https://github.com/forbearing/golib/compare/v0.0.54...v0.0.55
[v0.0.54]: https://github.com/forbearing/golib/compare/v0.0.53...v0.0.54
[v0.0.53]: https://github.com/forbearing/golib/compare/v0.0.52...v0.0.53
[v0.0.52]: https://github.com/forbearing/golib/compare/v0.0.51...v0.0.52
[v0.0.51]: https://github.com/forbearing/golib/compare/v0.0.50...v0.0.51
[v0.0.50]: https://github.com/forbearing/golib/compare/v0.0.49...v0.0.50
[v0.0.49]: https://github.com/forbearing/golib/compare/v0.0.48...v0.0.49
[v0.0.48]: https://github.com/forbearing/golib/compare/v0.0.47...v0.0.48
[v0.0.47]: https://github.com/forbearing/golib/compare/v0.0.46...v0.0.47
[v0.0.46]: https://github.com/forbearing/golib/compare/v0.0.45...v0.0.46
[v0.0.45]: https://github.com/forbearing/golib/compare/v0.0.44...v0.0.45
[v0.0.44]: https://github.com/forbearing/golib/compare/v0.0.43...v0.0.44
[v0.0.43]: https://github.com/forbearing/golib/compare/v0.0.42...v0.0.43
[v0.0.42]: https://github.com/forbearing/golib/compare/v0.0.41...v0.0.42
[v0.0.41]: https://github.com/forbearing/golib/compare/v0.0.40...v0.0.41
[v0.0.40]: https://github.com/forbearing/golib/compare/v0.0.39...v0.0.40
[v0.0.39]: https://github.com/forbearing/golib/compare/v0.0.38...v0.0.39
[v0.0.38]: https://github.com/forbearing/golib/compare/v0.0.37...v0.0.38
[v0.0.37]: https://github.com/forbearing/golib/compare/v0.0.36...v0.0.37
[v0.0.36]: https://github.com/forbearing/golib/compare/v0.0.35...v0.0.36
[v0.0.35]: https://github.com/forbearing/golib/compare/v0.0.34...v0.0.35
[v0.0.34]: https://github.com/forbearing/golib/compare/v0.0.33...v0.0.34
[v0.0.33]: https://github.com/forbearing/golib/compare/v0.0.32...v0.0.33
[v0.0.32]: https://github.com/forbearing/golib/compare/v0.0.31...v0.0.32
[v0.0.31]: https://github.com/forbearing/golib/compare/v0.0.30...v0.0.31
[v0.0.30]: https://github.com/forbearing/golib/compare/v0.0.29...v0.0.30
[v0.0.29]: https://github.com/forbearing/golib/compare/v0.0.28...v0.0.29
[v0.0.28]: https://github.com/forbearing/golib/compare/v0.0.27...v0.0.28
[v0.0.27]: https://github.com/forbearing/golib/compare/v0.0.26...v0.0.27
[v0.0.26]: https://github.com/forbearing/golib/compare/v0.0.25...v0.0.26
[v0.0.25]: https://github.com/forbearing/golib/compare/v0.0.24...v0.0.25
[v0.0.24]: https://github.com/forbearing/golib/compare/v0.0.23...v0.0.24
[v0.0.23]: https://github.com/forbearing/golib/compare/v0.0.22...v0.0.23
[v0.0.22]: https://github.com/forbearing/golib/compare/v0.0.21...v0.0.22
[v0.0.21]: https://github.com/forbearing/golib/compare/v0.0.20...v0.0.21
[v0.0.20]: https://github.com/forbearing/golib/compare/v0.0.19...v0.0.20
[v0.0.19]: https://github.com/forbearing/golib/compare/v0.0.18...v0.0.19
[v0.0.18]: https://github.com/forbearing/golib/compare/v0.0.17...v0.0.18
[v0.0.17]: https://github.com/forbearing/golib/compare/v0.0.16...v0.0.17
[v0.0.16]: https://github.com/forbearing/golib/compare/v0.0.15...v0.0.16
[v0.0.15]: https://github.com/forbearing/golib/compare/v0.0.14...v0.0.15
[v0.0.14]: https://github.com/forbearing/golib/compare/v0.0.13...v0.0.14
[v0.0.13]: https://github.com/forbearing/golib/compare/v0.0.12...v0.0.13
[v0.0.12]: https://github.com/forbearing/golib/compare/v0.0.11...v0.0.12
[v0.0.11]: https://github.com/forbearing/golib/compare/v0.0.10...v0.0.11
[v0.0.10]: https://github.com/forbearing/golib/compare/v0.0.9...v0.0.10
[v0.0.9]: https://github.com/forbearing/golib/compare/v0.0.8...v0.0.9
[v0.0.8]: https://github.com/forbearing/golib/compare/v0.0.7...v0.0.8
[v0.0.7]: https://github.com/forbearing/golib/compare/v0.0.6...v0.0.7
[v0.0.6]: https://github.com/forbearing/golib/compare/v0.0.5...v0.0.6
[v0.0.5]: https://github.com/forbearing/golib/compare/v0.0.4...v0.0.5
[v0.0.4]: https://github.com/forbearing/golib/compare/v0.0.3...v0.0.4
[v0.0.3]: https://github.com/forbearing/golib/compare/v0.0.2...v0.0.3
[v0.0.2]: https://github.com/forbearing/golib/compare/v0.0.1...v0.0.2
