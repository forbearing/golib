package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/creasty/defaults"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	APP_NAME         = "golib"
	noneExpireToken  = `fake_token`
	noneExpireUser   = "admin"
	noneExpirePass   = "admin"
	baseAuthUsername = "admin"
	baseAuthPassword = "admin"
)

var (
	App = new(Config)

	configPaths = []string{}
	configFile  = ""
	configName  = "config"
	configType  = "ini"

	registeredConfigs = make(map[string]any)
	registeredTypes   = make(map[string]reflect.Type)

	inited bool
	mu     sync.RWMutex
)

type (
	Mode string
	DB   string
)

const (
	ModeProd = "prod"
	ModeStg  = "stg"
	ModeDev  = "dev"
)

const (
	DBSqlite     = "sqlite"
	DBPostgre    = "postgres"
	DBMySQL      = "mysql"
	DBSQLServer  = "sqlserver"
	DBClickHouse = "clickhouse"
)

// Init initializes the application configuration
//
// Configuration priority (from highest to lowest):
// 1. Environment variables
// 2. Configuration file
// 3. Default values
func Init() (err error) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	SetDefaultValue()

	if len(configFile) > 0 {
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/")
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			var tempdir string
			if tempdir, err = os.MkdirTemp("", "golib_"); err != nil {
				return errors.Wrap(err, "failed to create temp dir")
			}
			if err = os.WriteFile(filepath.Join(tempdir, fmt.Sprintf("%s.%s", configName, configType)), nil, 0o644); err != nil {
				return errors.Wrap(err, "failed to create config file")
			}
		} else {
			return errors.Wrap(err, "failed to read config file")
		}
	}
	if err = viper.Unmarshal(App); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	for name, typ := range registeredTypes {
		registerType(name, typ)
	}
	inited = true

	return nil
}

// Register registers a configuration type with the given name.
// The type parameter T can be either struct type or pointer to struct type.
// The registered type will be used to create and initialize the configuration
// instance when loading configuration from file or environment variables.
//
// Configuration values are loaded in the following priority order (from highest to lowest):
// 1. Environment variables (format: SECTION_FIELD, e.g., NATS_USERNAME)
// 2. Configuration file values
// 3. Default values from struct tags
//
// The struct tag "default" can be used to set default values for fields.
// For time.Duration fields, you can use duration strings like "5s", "1m", etc.
//
// Register can be called before or after `Init`. If called before `Init`,
// the registration will be processed during initialization.
//
// Example usage:
//
//	type WechatConfig struct {
//		AppID     string `json:"app_id" mapstructure:"app_id" default:"myappid"`
//		AppSecret string `json:"app_secret" mapstructure:"app_secret" default:"myappsecret"`
//		Enable    bool   `json:"enable" mapstructure:"enable"`
//	}
//
//	type NatsConfig struct {
//		URL      string        `json:"url" mapstructure:"url" default:"nats://127.0.0.1:4222"`
//		Username string        `json:"username" mapstructure:"username" default:"nats"`
//		Password string        `json:"password" mapstructure:"password" default:"nats"`
//		Timeout  time.Duration `json:"timeout" mapstructure:"timeout" default:"5s"`
//		Enable   bool          `json:"enable" mapstructure:"enable"`
//	}
//
//	// Register with struct type
//	config.Register[WechatConfig]("wechat")
//
//	// Register with pointer type (equivalent to above)
//	config.Register[*NatsConfig]("nats")
//
// After registration, you can access the config using Get:
//
//	natsCfg := config.Get[NatsConfig]("nats")
//	// or with pointer
//	natsPtr := config.Get[*NatsConfig]("nats")
func Register[T any](name string) {
	mu.Lock()
	defer mu.Unlock()

	var t T
	typ := reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if inited {
		registerType(name, typ)
	} else {
		registeredTypes[name] = typ
	}
}

func registerType(name string, typ reflect.Type) {
	// Set default value from struct tag "default".
	cfg := reflect.New(typ).Interface()
	if err := defaults.Set(cfg); err != nil {
		zap.S().Warnw("failed to set default value", "name", name, "type", typ, "error", err)
	}
	// NOTE: package "defaults" not support set default value for time.Duration, so we should set it manually.
	setDefaultDurationFields(typ, reflect.ValueOf(cfg).Elem())

	// Set config value from config file.
	if err := viper.UnmarshalKey(name, cfg); err != nil {
		zap.S().Warnw("failed to unmarshal config", "name", name, "type", typ, "error", err)
	}

	// Set config value from environment variables.
	envCfg := reflect.New(typ).Interface()
	envPrefix := strings.ToUpper(name) + "_"
	v := reflect.ValueOf(envCfg).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		mapstructureTag := field.Tag.Get("mapstructure")
		if len(mapstructureTag) == 0 {
			continue
		}
		envKey := envPrefix + strings.ToUpper(mapstructureTag)
		if envVal, exists := os.LookupEnv(envKey); exists {
			fieldVal := v.Field(i)
			switch fieldVal.Kind() {
			case reflect.String:
				fieldVal.SetString(envVal)
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(envVal)
				if err == nil {
					fieldVal.SetBool(boolVal)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if field.Type == reflect.TypeOf(time.Duration(0)) {
					// handle time.Duration
					if duration, err := time.ParseDuration(envVal); err == nil {
						fieldVal.SetInt(int64(duration))
					}
				} else {
					if intVal, err := strconv.ParseInt(envVal, 10, 64); err == nil {
						fieldVal.SetInt(intVal)
					}
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if uintVal, err := strconv.ParseUint(envVal, 10, 64); err == nil {
					fieldVal.SetUint(uintVal)
				}
			case reflect.Float32, reflect.Float64:
				if floatVal, err := strconv.ParseFloat(envVal, 64); err == nil {
					fieldVal.SetFloat(floatVal)
				}

			}
		}
	}
	mergeNonZeroFields(reflect.ValueOf(cfg).Elem(), v)

	registeredConfigs[name] = cfg
}

func setDefaultDurationFields(typ reflect.Type, val reflect.Value) {
	if typ.Kind() != reflect.Struct {
		return
	}
	for i := range typ.NumField() {
		fieldTyp := typ.Field(i)
		fieldVal := val.Field(i)

		// Handle embedded structs
		if fieldTyp.Anonymous && fieldTyp.Type.Kind() == reflect.Struct {
			setDefaultDurationFields(fieldTyp.Type, fieldVal)
			continue
		}

		// Handle time.Duration field
		if fieldTyp.Type == reflect.TypeOf(time.Duration(0)) {
			// Check if the field has a default tag and its current value is zero
			if defaultValue, ok := fieldTyp.Tag.Lookup("default"); ok && fieldVal.Interface().(time.Duration) == 0 {
				// Parse the duration string
				if duration, err := time.ParseDuration(defaultValue); err == nil {
					fieldVal.Set(reflect.ValueOf(duration))
				} else {
					zap.S().Warnw("failed to parse duration default value",
						"field", fieldTyp.Name,
						"default", defaultValue,
						"error", err)
				}
			}
		}

		// Recursively process nested structs (if not embedded)
		if fieldTyp.Type.Kind() == reflect.Struct && !fieldTyp.Anonymous {
			setDefaultDurationFields(fieldTyp.Type, fieldVal)
		}

		// Handle pointer to struct
		if fieldTyp.Type.Kind() == reflect.Ptr && fieldTyp.Type.Elem().Kind() == reflect.Struct {
			// If the pointer is nil, initialize it
			if fieldVal.IsNil() {
				fieldVal.Set(reflect.New(fieldTyp.Type.Elem()))
			}
			setDefaultDurationFields(fieldTyp.Type.Elem(), fieldVal.Elem())
		}
	}
}

func mergeNonZeroFields(dst, src reflect.Value) {
	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		if !isZeroValue(srcField) {
			dst.Field(i).Set(srcField)
		}
	}
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// Get returns the registered configuration by name.
// The type parameter T must match the registered type or be a pointer to it,
// othherwise a zero value or nil pointer will be returned.
//
// Example usage:
//
//	config.Register[WechatConfig]("wechat")
//
//	// Get by value type - returns value
//	cfg1 := config.Get[WechatConfig]("wechat")
//
//	// Get by pointer type - returns pointer
//	cfg2 := config.Get[*WechatConfig]("wechat")
//
//	// Type mismatch - returns zero value
//	cfg3 := config.Get[OtherConfig]("wechat")
//
//	// Type mismatch - returns nil
//	cfg4 := config.Get[*OtherConfig]("wechat")
func Get[T any](name string) (t T) {
	mu.RLock()
	defer mu.RUnlock()

	config, exists := registeredConfigs[name]
	if !exists {
		zap.S().Warnw("config not found", "name", name)
		return
	}

	storedVal := reflect.ValueOf(config)
	storedTyp := storedVal.Elem().Type()
	destTyp := reflect.TypeOf(t)

	if storedTyp == destTyp {
		return storedVal.Elem().Interface().(T)
	}
	if destTyp.Kind() == reflect.Ptr {
		if storedTyp == destTyp.Elem() {
			return storedVal.Interface().(T)
		}
	}

	zap.S().Warnw("config type mismatch", "name", name, "stored", storedTyp.Name(), "dest", destTyp.Name())
	return t
}

// SetDefaultValue will set config default value
func SetDefaultValue() {
	viper.SetDefault("server.mode", ModeDev)
	viper.SetDefault("server.listen", "")
	viper.SetDefault("server.port", 9000)
	viper.SetDefault("server.db", DBSqlite)
	viper.SetDefault("server.domain", "")
	viper.SetDefault("server.enable_rbac", false)

	viper.SetDefault("auth.none_expire_token", noneExpireToken)
	viper.SetDefault("auth.none_expire_username", noneExpireUser)
	viper.SetDefault("auth.none_expire_passord", noneExpirePass)
	viper.SetDefault("auth.base_auth_username", baseAuthUsername)
	viper.SetDefault("auth.base_auth_password", baseAuthPassword)
	viper.SetDefault("auth.access_token_expire_duration", "2h")
	viper.SetDefault("auth.refresh_token_expire_duration", "168h")

	viper.SetDefault("logger.dir", "./logs")
	viper.SetDefault("logger.prefix", "")
	viper.SetDefault("logger.file", "")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("logger.format", "json")
	viper.SetDefault("logger.encoder", "json")
	viper.SetDefault("logger.max_age", 30)
	viper.SetDefault("logger.max_size", 100)
	viper.SetDefault("logger.max_backups", 1)

	viper.SetDefault("database.slow_query_threshold", 500*time.Millisecond)
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", 1*time.Hour)
	viper.SetDefault("database.conn_max_idle_time", 10*time.Minute)

	viper.SetDefault("sqlite.path", "./data.db")
	viper.SetDefault("sqlite.database", "main")
	viper.SetDefault("sqlite.is_memory", false)
	viper.SetDefault("sqlite.enable", true)

	viper.SetDefault("postgres.host", "127.0.0.1")
	viper.SetDefault("postgres.port", 5432)
	viper.SetDefault("postgres.database", "postgres")
	viper.SetDefault("postgres.username", "postgres")
	viper.SetDefault("postgres.password", "")
	viper.SetDefault("postgres.sslmode", "disable")
	viper.SetDefault("postgres.timezone", "UTC")
	viper.SetDefault("postgres.enable", true)

	viper.SetDefault("mysql.host", "127.0.0.1")
	viper.SetDefault("mysql.port", 3306)
	viper.SetDefault("mysql.database", "")
	viper.SetDefault("mysql.username", "root")
	viper.SetDefault("mysql.password", "toor")
	viper.SetDefault("mysql.charset", "utf8mb4")
	viper.SetDefault("mysql.enable", true)

	viper.SetDefault("clickhouse.host", "127.0.0.1")
	viper.SetDefault("clickhouse.port", 9000)
	viper.SetDefault("clickhouse.database", "default")
	viper.SetDefault("clickhouse.username", "default")
	viper.SetDefault("clickhouse.password", "")
	viper.SetDefault("clickhouse.dial_timeout", "5s")
	viper.SetDefault("clickhouse.read_timeout", "30s")
	viper.SetDefault("clickhouse.write_timeout", "30s")
	viper.SetDefault("clickhouse.compress", false)
	viper.SetDefault("clickhouse.debug", false)
	viper.SetDefault("clickhouse.enable", false)

	viper.SetDefault("sqlserver.host", "127.0.0.1")
	viper.SetDefault("sqlserver.port", 1433)
	viper.SetDefault("sqlserver.database", "")
	viper.SetDefault("sqlserver.username", "sa")
	viper.SetDefault("sqlserver.password", "")
	viper.SetDefault("sqlserver.encrypt", false)
	viper.SetDefault("sqlserver.trust_server", true)
	viper.SetDefault("sqlserver.app_name", "golib")
	viper.SetDefault("sqlserver.enable", false)

	viper.SetDefault("redis.addr", "127.0.0.1:6379")
	viper.SetDefault("redis.addrs", []string{"127.0.0.1:6379"})
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.pool_size", runtime.NumCPU())
	viper.SetDefault("redis.namespace", APP_NAME)
	viper.SetDefault("redis.expiration", 0)
	viper.SetDefault("redis.cluster_mode", false)
	viper.SetDefault("redis.dial_timeout", 0)
	viper.SetDefault("redis.read_timeout", 0)
	viper.SetDefault("redis.write_timeout", 0)
	viper.SetDefault("redis.min_idle_conns", 0)
	viper.SetDefault("redis.max_retries", 0)
	viper.SetDefault("redis.min_retry_backoff", 0)
	viper.SetDefault("redis.max_retry_backoff", 0)
	viper.SetDefault("redis.enable_tls", false)
	viper.SetDefault("redis.cert_file", "")
	viper.SetDefault("redis.key_file", 0)
	viper.SetDefault("redis.ca_file", "")
	viper.SetDefault("redis.insecure_skip_verify", false)
	viper.SetDefault("redis.enable", false)

	viper.SetDefault("elasticsearch.hosts", "127.0.0.1")
	viper.SetDefault("elasticsearch.username", "")
	viper.SetDefault("elasticsearch.password", "")
	viper.SetDefault("elasticsearch.cloud_id", "")
	viper.SetDefault("elasticsearch.api_key", "")
	viper.SetDefault("elasticsearch.service_token", "")
	viper.SetDefault("elasticsearch.certificate_fingerprint", "")
	viper.SetDefault("elasticsearch.enable", false)

	viper.SetDefault("mongo.host", "127.0.0.1")
	viper.SetDefault("mongo.port", 27017)
	viper.SetDefault("mongo.username", "")
	viper.SetDefault("mongo.password", "")
	viper.SetDefault("mongo.database", "")
	viper.SetDefault("mongo.auth_source", "admin")
	viper.SetDefault("mongo.max_pool_size", 0)
	viper.SetDefault("mongo.min_pool_size", 0)
	viper.SetDefault("mongo.connect_timeout", 0)
	viper.SetDefault("mongo.server_selection_timeout", 0)
	viper.SetDefault("mongo.max_conn_idle_time", 0)
	viper.SetDefault("mongo.max_connecting", 0)
	viper.SetDefault("mongo.read_concern", "")
	viper.SetDefault("mongo.write_concern", "")
	viper.SetDefault("mongo.enable_tls", false)
	viper.SetDefault("mongo.cert_file", "")
	viper.SetDefault("mongo.key_file", "")
	viper.SetDefault("mongo.ca_file", "")
	viper.SetDefault("mongo.insecure_skip_verify", false)
	viper.SetDefault("mongo.enable", false)

	viper.SetDefault("ldap.host", "127.0.0.1")
	viper.SetDefault("ldap.port", 389)
	viper.SetDefault("ldap.use_ssl", false)
	viper.SetDefault("ldap.bind_dn", "")
	viper.SetDefault("ldap.bind_password", "")
	viper.SetDefault("ldap.base_dn", "")
	viper.SetDefault("ldap.search_filter", "")
	viper.SetDefault("ldap.enable", false)

	viper.SetDefault("influxdb.host", "127.0.0.1")
	viper.SetDefault("influxdb.port", 8086)
	viper.SetDefault("influxdb.admin_password", "")
	viper.SetDefault("influxdb.admin_token", "")
	viper.SetDefault("influxdb.admin_org", "")
	viper.SetDefault("influxdb.bucket", "")
	viper.SetDefault("influxdb.write_interval", 1*time.Minute)
	viper.SetDefault("influxdb.enable", false)

	viper.SetDefault("minio.endpoint", "127.0.0.1:9000")
	viper.SetDefault("minio.region", "")
	viper.SetDefault("minio.access_key", "")
	viper.SetDefault("minio.secret_key", "")
	viper.SetDefault("minio.bucket", "")
	viper.SetDefault("minio.use_ssl", false)
	viper.SetDefault("minio.enable", false)

	viper.SetDefault("s3.endpoint", "")
	viper.SetDefault("s3.region", "")
	viper.SetDefault("s3.access_key_id", "")
	viper.SetDefault("s3.secret_access_key", "")
	viper.SetDefault("s3.bucket", "")
	viper.SetDefault("s3.use_ssl", false)
	viper.SetDefault("s3.enable", false)

	viper.SetDefault("mqtt.addr", "127.0.0.1:1883")
	viper.SetDefault("mqtt.username", "")
	viper.SetDefault("mqtt.password", "")
	viper.SetDefault("mqtt.client_prefix", "")
	viper.SetDefault("mqtt.connect_timeout", 10*time.Second)
	viper.SetDefault("mqtt.keepalive", 1*time.Minute)
	viper.SetDefault("mqtt.clean_session", true)
	viper.SetDefault("mqtt.auto_reconnect", true)
	viper.SetDefault("mqtt.use_tls", false)
	viper.SetDefault("mqtt.cert_file", "")
	viper.SetDefault("mqtt.key_file", "")
	viper.SetDefault("mqtt.insecure_skip_verify", true)
	viper.SetDefault("mqtt.enable", false)

	viper.SetDefault("feishu.app_id", "")
	viper.SetDefault("feishu.app_secret", "")
	viper.SetDefault("feishu.msg_app_id", "")
	viper.SetDefault("feishu.msg_app_secret", "")
	viper.SetDefault("feishu.enable", false)

	viper.SetDefault("debug.enable_statsviz", false)
	viper.SetDefault("debug.enable_pprof", false)
	viper.SetDefault("debug.enable_gops", false)
	viper.SetDefault("debug.statsviz_listen", "127.0.0.1")
	viper.SetDefault("debug.pprof_listen", "127.0.0.1")
	viper.SetDefault("debug.gops_listen", "127.0.0.1")
	viper.SetDefault("debug.statsviz_port", 10000)
	viper.SetDefault("debug.pprof_port", 10001)
	viper.SetDefault("debug.gops_port", 10002)
}

// SetConfigFile set the config file path.
// You should always call this funtion before `Init`.
func SetConfigFile(file string) {
	mu.Lock()
	defer mu.Unlock()
	configFile = file
}

// SetConfigName set the config file name, default to 'config'.
// NOTE: any suffix will be ignored and the default file type is ini.
// You should always call this funtion before `Init`.
func SetConfigName(name string) {
	mu.Lock()
	defer mu.Unlock()
	configName = name
}

// SetConfigType set the config file type, default to 'ini'.
// You should always call this funtion before `Init`.
func SetConfigType(typ string) {
	mu.Lock()
	defer mu.Unlock()
	configType = typ
}

// AddPath add custom config path. default: ./config, /etc
// You should always call this funtion before `Init`.
func AddPath(paths ...string) {
	mu.Lock()
	defer mu.Unlock()
	configPaths = append(configPaths, paths...)
}

// Save config instance to file.
func Save(filename string) error {
	return viper.WriteConfigAs(filename)
}

type Config struct {
	ServerConfig        `json:"server" mapstructure:"server" ini:"server" yaml:"server"`
	AuthConfig          `json:"auth" mapstructure:"auth" ini:"auth" yaml:"auth"`
	DatabaseConfig      `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	SqliteConfig        `json:"sqlite" mapstructure:"sqlite" ini:"sqlite" yaml:"sqlite"`
	PostgreConfig       `json:"postgres" mapstructure:"postgres" ini:"postgres" yaml:"postgres"`
	MySQLConfig         `json:"mysql" mapstructure:"mysql" ini:"mysql" yaml:"mysql"`
	SQLServerConfig     `json:"sqlserver" mapstructure:"sqlserver" ini:"sqlserver" yaml:"sqlserver"`
	ClickhouseConfig    `json:"clickhouse" mapstructure:"clickhouse" ini:"clickhouse" yaml:"clickhouse"`
	RedisConfig         `json:"redis" mapstructure:"redis" ini:"redis" yaml:"redis"`
	ElasticsearchConfig `json:"elasticsearch" mapstructure:"elasticsearch" ini:"elasticsearch" yaml:"elasticsearch"`
	MongoConfig         `json:"mongo" mapstructure:"mongo" ini:"mongo" yaml:"mongo"`
	MinioConfig         `json:"minio" mapstructure:"minio" ini:"minio" yaml:"minio"`
	S3Config            `json:"s3" mapstructure:"s3" ini:"s3" yaml:"s3"`
	LoggerConfig        `json:"logger" mapstructure:"logger" ini:"logger" yaml:"logger"`
	LdapConfig          `json:"ldap" mapstructure:"ldap" ini:"ldap" yaml:"ldap"`
	InfluxdbConfig      `json:"influxdb" mapstructure:"influxdb" ini:"influxdb" yaml:"influxdb"`
	MqttConfig          `json:"mqtt" mapstructure:"mqtt" ini:"mqtt" yaml:"mqtt"`
	FeishuConfig        `json:"feishu" mapstructure:"feishu" ini:"feishu" yaml:"feishu"`
	DebugConfig         `json:"debug" mapstructure:"debug" ini:"debug" yaml:"debug"`
}

type ServerConfig struct {
	Mode       Mode   `json:"mode" mapstructure:"mode" ini:"mode" yaml:"mode"`
	Listen     string `json:"listen" mapstructure:"listen" ini:"listen" yaml:"listen"`
	Port       int    `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	DB         DB     `json:"db" mapstructure:"db" ini:"db" yaml:"db"`
	Domain     string `json:"domain" mapstructure:"domain" ini:"domain" yaml:"domain"`
	EnableRBAC bool   `json:"enable_rbac" mapstructure:"enable_rbac" ini:"enable_rbac" yaml:"enable_rbac"`
}

type AuthConfig struct {
	NoneExpireToken            string        `json:"none_expire_token" mapstructure:"none_expire_token" ini:"none_expire_token" yaml:"none_expire_token"`
	NoneExpireUsername         string        `json:"none_expire_username" mapstructure:"none_expire_username" ini:"none_expire_username" yaml:"none_expire_username"`
	NoneExpirePassword         string        `json:"none_expire_passord" mapstructure:"none_expire_passord" ini:"none_expire_passord" yaml:"none_expire_passord"`
	BaseAuthUsername           string        `json:"base_auth_username" mapstructure:"base_auth_username" ini:"base_auth_username" yaml:"base_auth_username"`
	BaseAuthPassword           string        `json:"base_auth_password" mapstructure:"base_auth_password" ini:"base_auth_password" yaml:"base_auth_password"`
	AccessTokenExpireDuration  time.Duration `json:"access_token_expire_duration" mapstructure:"access_token_expire_duration" ini:"access_token_expire_duration" yaml:"access_token_expire_duration"`
	RefreshTokenExpireDuration time.Duration `json:"refresh_token_expire_duration" mapstructure:"refresh_token_expire_duration" ini:"refresh_token_expire_duration" yaml:"refresh_token_expire_duration"`
}

// LoggerConfig represents section "logger" for client-side or server-side configuration,
// and there is only one copy during the application entire lifetime.
type LoggerConfig struct {
	// Dir specifies which direcotory log to.
	Dir string `json:"dir" ini:"dir" yaml:"dir" mapstructure:"dir"`

	// Prefix specifies the log prefix.
	// You can set the prefix name to your project name.
	Prefix string `json:"prefix" ini:"prefix" yaml:"prefix" mapstructure:"prefix"`

	// File specifies the which file log to.
	// If value is "/dev/stdout", log to os.Stdout.
	// If value is "/dev/stderr", log to os.Stderr.
	// If value is empty(length is zero), log to os.Stdout.
	File string `json:"file" ini:"file" yaml:"file" mapstructure:"file"`

	// Level specifies the log level,  supported values are: (error|warn|warning|info|debug).
	// The value default to "info" and ignore case.
	Level string `json:"level" ini:"level" yaml:"level" mapstructure:"level"`

	// Format specifies the log format, supported values are: (json|text).
	// The Value default to "text" and ignore case.
	Format string `json:"format" ini:"format" yaml:"format" mapstructure:"format"`

	// Encoder is the same as LogFormat.
	Encoder string `json:"encoder" ini:"encoder" yaml:"encoder" mapstructure:"encoder"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.
	// uint is "day" and default to 7.
	MaxAge uint `json:"max_age" ini:"max_age" yaml:"max_age" mapstructure:"max_age"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated, default to 1MB.
	MaxSize uint `json:"max_size" ini:"max_size" yaml:"max_size" mapstructure:"max_size"`

	// MaxBackups is the maximum number of old log files to retain.
	// The value default to 3.
	MaxBackups uint `json:"max_backups" ini:"max_backups" yaml:"max_backups" mapstructure:"max_backups"`
}

type DatabaseConfig struct {
	SlowQueryThreshold time.Duration `json:"slow_query_threshold" mapstructure:"slow_query_threshold" ini:"slow_query_threshold" yaml:"slow_query_threshold"`
	MaxIdleConns       int           `json:"max_idle_conns" mapstructure:"max_idle_conns" ini:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns       int           `json:"max_open_conns" mapstructure:"max_open_conns" ini:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime    time.Duration `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime" ini:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `json:"conn_max_idle_time" mapstructure:"conn_max_idle_time" ini:"conn_max_idle_time" yaml:"conn_max_idle_time"`
}

type SqliteConfig struct {
	Path     string `json:"path" mapstructure:"path" ini:"path" yaml:"path"`
	Database string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	IsMemory bool   `json:"is_memory" mapstructure:"is_memory" ini:"is_memory" yaml:"is_memory"`
	Enable   bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type PostgreConfig struct {
	Host     string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port     uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Database string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	Username string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	SSLMode  string `json:"sslmode" mapstructure:"sslmode" ini:"sslmode" yaml:"sslmode"`
	TimeZone string `json:"timezone" mapstructure:"timezone" ini:"timezone" yaml:"timezone"`
	Enable   bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type MySQLConfig struct {
	Host     string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port     uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Database string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	Username string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	Charset  string `json:"charset" mapstructure:"charset" ini:"charset" yaml:"charset"`
	Enable   bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type ClickhouseConfig struct {
	Host         string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port         uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Database     string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	Username     string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password     string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	DialTimeout  string `json:"dial_timeout" mapstructure:"dial_timeout" ini:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  string `json:"read_timeout" mapstructure:"read_timeout" ini:"read_timeout" yaml:"read_timeout"`
	WriteTimeout string `json:"write_timeout" mapstructure:"write_timeout" ini:"write_timeout" yaml:"write_timeout"`
	Compress     bool   `json:"compress" mapstructure:"compress" ini:"compress" yaml:"compress"`
	Debug        bool   `json:"debug" mapstructure:"debug" ini:"debug" yaml:"debug"`
	Enable       bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type SQLServerConfig struct {
	Host        string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port        uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Database    string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	Username    string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password    string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	Encrypt     bool   `json:"encrypt" mapstructure:"encrypt" ini:"encrypt" yaml:"encrypt"`
	TrustServer bool   `json:"trust_server" mapstructure:"trust_server" ini:"trust_server" yaml:"trust_server"`
	AppName     string `json:"app_name" mapstructure:"app_name" ini:"app_name" yaml:"app_name"`
	Enable      bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type RedisConfig struct {
	Addr        string        `json:"addr" mapstructure:"addr" ini:"addr" yaml:"addr"`
	Addrs       []string      `json:"addrs" mapstructure:"addrs" ini:"addrs" yaml:"addrs"`
	DB          int           `json:"db" mapstructure:"db" ini:"db" yaml:"db"`
	Password    string        `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	Namespace   string        `json:"namespace" mapstructure:"namespace" ini:"namespace" yaml:"namespace"`
	PoolSize    int           `json:"pool_size" mapstructure:"pool_size" ini:"pool_size" yaml:"pool_size"`
	Expiration  time.Duration `json:"expiration" mapstructure:"expiration" ini:"expiration" yaml:"expiration"`
	ClusterMode bool          `json:"cluster_mode" mapstructure:"cluster_mode" ini:"cluster_mode" yaml:"cluster_mode"`

	DialTimeout     time.Duration `json:"dial_timeout" mapstructure:"dial_timeout" ini:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout     time.Duration `json:"read_timeout" mapstructure:"read_timeout" ini:"read_timeout" yaml:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout" mapstructure:"write_timeout" ini:"write_timeout" yaml:"write_timeout"`
	MinIdleConns    int           `json:"min_idle_conns" mapstructure:"min_idle_conns" ini:"min_idle_conns" yaml:"min_idle_conns"`
	MaxRetries      int           `json:"max_retries" mapstructure:"max_retries" ini:"max_retries" yaml:"max_retries"`
	MinRetryBackoff time.Duration `json:"min_retry_backoff" mapstructure:"min_retry_backoff" ini:"min_retry_backoff" yaml:"min_retry_backoff"`
	MaxRetryBackoff time.Duration `json:"max_retry_backoff" mapstructure:"max_retry_backoff" ini:"max_retry_backoff" yaml:"max_retry_backoff"`

	EnableTLS          bool   `json:"enable_tls" mapstructure:"enable_tls" ini:"enable_tls" yaml:"enable_tls"`
	CertFile           string `json:"cert_file" mapstructure:"cert_file" ini:"cert_file" yaml:"cert_file"`
	KeyFile            string `json:"key_file" mapstructure:"key_file" ini:"key_file" yaml:"key_file"`
	CAFile             string `json:"ca_file" mapstructure:"ca_file" ini:"ca_file" yaml:"ca_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify" ini:"insecure_skip_verify" yaml:"insecure_skip_verify"`

	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type ElasticsearchConfig struct {
	Hosts                  string `json:"hosts" mapstructure:"hosts" ini:"hosts" yaml:"hosts"`
	Username               string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password               string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	CloudID                string `json:"cloud_id" mapstructure:"cloud_id" ini:"cloud_id" yaml:"cloud_id"`
	APIKey                 string `json:"api_key" mapstructure:"api_key" ini:"api_key" yaml:"api_key"`
	ServiceToken           string `json:"service_token" mapstructure:"service_token" ini:"service_token" yaml:"service_token"`
	CertificateFingerprint string `json:"certificate_fingerprint" mapstructure:"certificate_fingerprint" ini:"certificate_fingerprint" yaml:"certificate_fingerprint"`
	Enable                 bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type MongoConfig struct {
	Host        string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port        int    `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Username    string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password    string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	Database    string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	AuthSource  string `json:"auth_source" mapstructure:"auth_source" ini:"auth_source" yaml:"auth_source"`
	MaxPoolSize uint64 `json:"max_pool_size" mapstructure:"max_pool_size" ini:"max_pool_size" yaml:"max_pool_size"`
	MinPoolSize uint64 `json:"min_pool_size" mapstructure:"min_pool_size" ini:"min_pool_size" yaml:"min_pool_size"`

	ConnectTimeout         time.Duration `json:"connect_timeout" mapstructure:"connect_timeout" ini:"connect_timeout" yaml:"connect_timeout"`
	ServerSelectionTimeout time.Duration `json:"server_selection_timeout" mapstructure:"server_selection_timeout" ini:"server_selection_timeout" yaml:"server_selection_timeout"`
	MaxConnIdleTime        time.Duration `json:"max_conn_idle_time" mapstructure:"max_conn_idle_time" ini:"max_conn_idle_time" yaml:"max_conn_idle_time"`
	MaxConnecting          uint64        `json:"max_connecting" mapstructure:"max_connecting" ini:"max_connecting" yaml:"max_connecting"`

	ReadConcern  ReadConcern  `json:"read_concern" mapstructure:"read_concern" ini:"read_concern" yaml:"read_concern"`
	WriteConcern WriteConcern `json:"write_concern" mapstructure:"write_concern" ini:"write_concern" yaml:"write_concern"`

	EnableTLS          bool   `json:"enable_tls" mapstructure:"enable_tls" ini:"enable_tls" yaml:"enable_tls"`
	CertFile           string `json:"cert_file" mapstructure:"cert_file" ini:"cert_file" yaml:"cert_file"`
	KeyFile            string `json:"key_file" mapstructure:"key_file" ini:"key_file" yaml:"key_file"`
	CAFile             string `json:"ca_file" mapstructure:"ca_file" ini:"ca_file" yaml:"ca_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify" ini:"insecure_skip_verify" yaml:"insecure_skip_verify"`

	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type ReadConcern string

const (
	ReadConcernLocal        ReadConcern = "local"
	ReadConcernMajority     ReadConcern = "majority"
	ReadConcernAvailable    ReadConcern = "available"
	ReadConcernLinearizable ReadConcern = "linearizable"
	ReadConcernSnapshot     ReadConcern = "snapshot"
)

type WriteConcern string

const (
	WriteConcernMajority  WriteConcern = "majority"
	WriteConcernJournaled WriteConcern = "journaled"
	WriteConcernW0        WriteConcern = "0"
	WriteConcernW1        WriteConcern = "1"
	WriteConcernW2        WriteConcern = "2"
	WriteConcernW3        WriteConcern = "3"
	WriteConcernW4        WriteConcern = "4"
	WriteConcernW5        WriteConcern = "5"
	WriteConcernW6        WriteConcern = "6"
	WriteConcernW7        WriteConcern = "7"
	WriteConcernW8        WriteConcern = "8"
	WriteConcernW9        WriteConcern = "9"
)

// LdapConfig
// For example:
// [ldap]
// host = example.cn
// port = 389
// use_ssl =  false
// bind_dn =  my_bind_dn
// bind_password = mypass
// base_dn = my_base_dn
// search_filter = (sAMAccountName=%s)
type LdapConfig struct {
	Host         string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port         uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	UseSsl       bool   `json:"use_ssl" mapstructure:"use_ssl" ini:"use_ssl" yaml:"use_ssl"`
	BindDN       string `json:"bind_dn" mapstructure:"bind_dn" ini:"bind_dn" yaml:"bind_dn"`
	BindPassword string `json:"bind_password" mapstructure:"bind_password" ini:"bind_password" yaml:"bind_password"`
	BaseDN       string `json:"base_dn" mapstructure:"base_dn" ini:"base_dn" yaml:"base_dn"`
	SearchFilter string `json:"search_filter" mapstructure:"search_filter" ini:"search_filter" yaml:"search_filter"`
	Enable       bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

// InfluxdbConfig is the configuration of influxdb.
// For example:
// [influxdb]
// admin_password = mypass
// admin_token = mytoken
// admin_org = example.com
// bucket = mybucket
type InfluxdbConfig struct {
	// - INFLUXDB_HTTP_AUTH_ENABLED=true
	// - INFLUXDB_ADMIN_USER_PASSWORD=7MXsrGqj3AuEt9rgSq7U
	// - INFLUXDB_ADMIN_USER_TOKEN=7MXsrGqj3AuEt9rgSq7U
	// - INFLUXDB_USER_BUCKET=mybucket
	// - INFLUXDB_ADMIN_ORG=example.com
	Host          string        `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port          uint          `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	AdminPassword string        `json:"admin_password" mapstructure:"admin_password" ini:"admin_password" yaml:"admin_password"`
	AdminToken    string        `json:"admin_token" mapstructure:"admin_token" ini:"admin_token" yaml:"admin_token"`
	AdminOrg      string        `json:"admin_org" mapstructure:"admin_org" ini:"admin_org" yaml:"admin_org"`
	Bucket        string        `json:"bucket" mapstructure:"bucket" ini:"bucket" yaml:"bucket"`
	WriteInterval time.Duration `json:"write_interval" mapstructure:"write_interval" ini:"write_interval" yaml:"write_interval"`
	Enable        bool          `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type MinioConfig struct {
	Endpoint  string `json:"endpoint" mapstructure:"endpoint" ini:"endpoint" yaml:"endpoint"`
	Region    string `json:"region" mapstructure:"region" ini:"region" yaml:"region"`
	AccessKey string `json:"access_key" mapstructure:"access_key" ini:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" mapstructure:"secret_key" ini:"secret_key" yaml:"secret_key"`
	Bucket    string `json:"bucket" mapstructure:"bucket" ini:"bucket" yaml:"bucket"`
	UseSsl    bool   `json:"use_ssl" mapstructure:"use_ssl" ini:"use_ssl" yaml:"use_ssl"`
	Enable    bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type S3Config struct {
	Endpoint        string `json:"endpoint" mapstructure:"endpoint" ini:"endpoint" yaml:"endpoint"`
	Region          string `json:"region" mapstructure:"region" ini:"region" yaml:"region"`
	AccessKeyID     string `json:"access_key_id" mapstructure:"access_key_id" ini:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" mapstructure:"secret_access_key" ini:"secret_access_key" yaml:"secret_access_key"`
	Bucket          string `json:"bucket" mapstructure:"bucket" ini:"bucket" yaml:"bucket"`
	UseSsl          bool   `json:"use_ssl" mapstructure:"use_ssl" ini:"use_ssl" yaml:"use_ssl"`
	Enable          bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type MqttConfig struct {
	Addr               string        `json:"addr" mapstructure:"addr" ini:"addr" yaml:"addr"`
	Username           string        `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password           string        `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	ClientPrefix       string        `json:"client_prefix" mapstructure:"client_prefix" ini:"client_prefix" yaml:"client_prefix"`
	ConnectTimeout     time.Duration `json:"connect_timeout" mapstructure:"connect_timeout" ini:"connect_timeout" yaml:"connect_timeout"`
	Keepalive          time.Duration `json:"keepalive" mapstructure:"keepalive" ini:"keepalive" yaml:"keepalive"`
	CleanSession       bool          `json:"clean_session" mapstructure:"clean_session" ini:"clean_session" yaml:"clean_session"`
	AutoReconnect      bool          `json:"auto_reconnect" mapstructure:"auto_reconnect" ini:"auto_reconnect" yaml:"auto_reconnect"`
	UseTLS             bool          `json:"use_tls" mapstructure:"use_tls" ini:"use_tls" yaml:"use_tls"`
	CertFile           string        `json:"cert_file" mapstructure:"cert_file" ini:"cert_file" yaml:"cert_file"`
	KeyFile            string        `json:"key_file" mapstructure:"key_file" ini:"key_file" yaml:"key_file"`
	InsecureSkipVerify bool          `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify" ini:"insecure_skip_verify" yaml:"insecure_skip_verify"`
	Enable             bool          `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type FeishuConfig struct {
	AppID        string `json:"app_id" mapstructure:"app_id" ini:"app_id" yaml:"app_id"`
	AppSecret    string `json:"app_secret" mapstructure:"app_secret" ini:"app_secret" yaml:"app_secret"`
	MsgAppID     string `json:"msg_app_id" mapstructure:"msg_app_id" ini:"msg_app_id" yaml:"msg_app_id"`
	MsgAppSecret string `json:"msg_app_secret" mapstructure:"msg_app_secret" ini:"msg_app_secret" yaml:"msg_app_secret"`
	Enable       bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

type DebugConfig struct {
	EnableStatsviz bool   `json:"enable_statsviz" mapstructure:"enable_statsviz" ini:"enable_statsviz" yaml:"enable_statsviz"`
	StatsvizListen string `json:"statsviz_listen" mapstructure:"statsviz_listen" ini:"statsviz_listen" yaml:"statsviz_listen"`
	StatsvizPort   int    `json:"statsviz_port" mapstructure:"statsviz_port" ini:"statsviz_port" yaml:"statsviz_port"`

	EnablePprof bool   `json:"enable_pprof" mapstructure:"enable_pprof" ini:"enable_pprof" yaml:"enable_pprof"`
	PprofListen string `json:"pprof_listen" mapstructure:"pprof_listen" ini:"pprof_listen" yaml:"pprof_listen"`
	PprofPort   int    `json:"pprof_port" mapstructure:"pprof_port" ini:"pprof_port" yaml:"pprof_port"`

	EnableGops bool   `json:"enable_gops" mapstructure:"enable_gops" ini:"enable_gops" yaml:"enable_gops"`
	GopsListen string `json:"gops_listen" mapstructure:"gops_listen" ini:"gops_listen" yaml:"gops_listen"`
	GopsPort   int    `json:"gops_port" mapstructure:"gops_port" ini:"gops_port" yaml:"gops_port"`
}
