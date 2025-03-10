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
	cv     = viper.New()
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
	cv.AutomaticEnv()
	cv.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	new(Config).setDefault()

	if len(configFile) > 0 {
		cv.SetConfigFile(configFile)
	} else {
		cv.SetConfigName(configName)
		cv.SetConfigType(configType)
	}
	cv.AddConfigPath(".")
	cv.AddConfigPath("/etc/")
	for _, path := range configPaths {
		cv.AddConfigPath(path)
	}

	if err = cv.ReadInConfig(); err != nil {
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
	if err = cv.Unmarshal(App); err != nil {
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
	if err := cv.UnmarshalKey(name, cfg); err != nil {
		zap.S().Warnw("failed to unmarshal config", "name", name, "type", typ, "error", err)
	}

	// Set config value from environment variables.
	envCfg := reflect.New(typ).Interface()
	envPrefix := strings.ToUpper(name) + "_"
	v := reflect.ValueOf(envCfg).Elem()
	t := v.Type()
	for i := range t.NumField() {
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
	return cv.WriteConfigAs(filename)
}

type Config struct {
	Server        `json:"server" mapstructure:"server" ini:"server" yaml:"server"`
	Auth          `json:"auth" mapstructure:"auth" ini:"auth" yaml:"auth"`
	Database      `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	Sqlite        `json:"sqlite" mapstructure:"sqlite" ini:"sqlite" yaml:"sqlite"`
	Postgres      `json:"postgres" mapstructure:"postgres" ini:"postgres" yaml:"postgres"`
	MySQL         `json:"mysql" mapstructure:"mysql" ini:"mysql" yaml:"mysql"`
	SQLServer     `json:"sqlserver" mapstructure:"sqlserver" ini:"sqlserver" yaml:"sqlserver"`
	Clickhouse    `json:"clickhouse" mapstructure:"clickhouse" ini:"clickhouse" yaml:"clickhouse"`
	Redis         `json:"redis" mapstructure:"redis" ini:"redis" yaml:"redis"`
	Elasticsearch `json:"elasticsearch" mapstructure:"elasticsearch" ini:"elasticsearch" yaml:"elasticsearch"`
	Mongo         `json:"mongo" mapstructure:"mongo" ini:"mongo" yaml:"mongo"`
	Kafka         `json:"kafka" mapstructure:"kafka" ini:"kafka" yaml:"kafka"`
	Minio         `json:"minio" mapstructure:"minio" ini:"minio" yaml:"minio"`
	S3            `json:"s3" mapstructure:"s3" ini:"s3" yaml:"s3"`
	Logger        `json:"logger" mapstructure:"logger" ini:"logger" yaml:"logger"`
	Ldap          `json:"ldap" mapstructure:"ldap" ini:"ldap" yaml:"ldap"`
	Influxdb      `json:"influxdb" mapstructure:"influxdb" ini:"influxdb" yaml:"influxdb"`
	Mqtt          `json:"mqtt" mapstructure:"mqtt" ini:"mqtt" yaml:"mqtt"`
	Nats          `json:"nats" mapstructure:"nats" ini:"nats" yaml:"nats"`
	Feishu        `json:"feishu" mapstructure:"feishu" ini:"feishu" yaml:"feishu"`
	Debug         `json:"debug" mapstructure:"debug" ini:"debug" yaml:"debug"`
}

// setDefaultValue will set config default value
func (*Config) setDefault() {
	new(Server).setDefault()
	new(Auth).setDefault()
	new(Logger).setDefault()
	new(Database).setDefault()
	new(Sqlite).setDefault()
	new(Postgres).setDefault()
	new(MySQL).setDefault()
	new(Clickhouse).setDefault()
	new(SQLServer).setDefault()
	new(Redis).setDefault()
	new(Elasticsearch).setDefault()
	new(Mongo).setDefault()
	new(Kafka).setDefault()
	new(Ldap).setDefault()
	new(Influxdb).setDefault()
	new(Minio).setDefault()
	new(S3).setDefault()
	new(Mqtt).setDefault()
	new(Nats).setDefault()
	new(Feishu).setDefault()
	new(Debug).setDefault()
}

const (
	SERVER_DOMAIN      = "SERVER_DOMAIN"
	SERVER_MODE        = "SERVER_MODE"
	SERVER_LISTEN      = "SERVER_LISTEN"
	SERVER_PORT        = "SERVER_PORT"
	SERVER_DB          = "SERVER_DB"
	SERVER_ENABLE_RBAC = "SERVER_ENABLE_RBAC"
)

type Server struct {
	Mode       Mode   `json:"mode" mapstructure:"mode" ini:"mode" yaml:"mode"`
	Listen     string `json:"listen" mapstructure:"listen" ini:"listen" yaml:"listen"`
	Port       int    `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	DB         DB     `json:"db" mapstructure:"db" ini:"db" yaml:"db"`
	Domain     string `json:"domain" mapstructure:"domain" ini:"domain" yaml:"domain"`
	EnableRBAC bool   `json:"enable_rbac" mapstructure:"enable_rbac" ini:"enable_rbac" yaml:"enable_rbac"`
}

func (*Server) setDefault() {
	cv.SetDefault("server.mode", ModeDev)
	cv.SetDefault("server.listen", "")
	cv.SetDefault("server.port", 9000)
	cv.SetDefault("server.db", DBSqlite)
	cv.SetDefault("server.domain", "")
	cv.SetDefault("server.enable_rbac", false)
}

const (
	AUTH_NONE_EXPIRE_TOKEN             = "AUTH_NONE_EXPIRE_TOKEN"
	AUTH_NONE_EXPIRE_USERNAME          = "AUTH_NONE_EXPIRE_USERNAME"
	AUTH_NONE_EXPIRE_PASSWORD          = "AUTH_NONE_EXPIRE_PASSORD"
	AUTH_BASE_AUTH_USERNAME            = "AUTH_BASE_AUTH_USERNAME"
	AUTH_BASE_AUTH_PASSWORD            = "AUTH_BASE_AUTH_PASSWORD"
	AUTH_ACCESS_TOKEN_EXPIRE_DURATION  = "AUTH_ACCESS_TOKEN_EXPIRE_DURATION"
	AUTH_REFRESH_TOKEN_EXPIRE_DURATION = "AUTH_REFRESH_TOKEN_EXPIRE_DURATION"
)

type Auth struct {
	NoneExpireToken            string        `json:"none_expire_token" mapstructure:"none_expire_token" ini:"none_expire_token" yaml:"none_expire_token"`
	NoneExpireUsername         string        `json:"none_expire_username" mapstructure:"none_expire_username" ini:"none_expire_username" yaml:"none_expire_username"`
	NoneExpirePassword         string        `json:"none_expire_passord" mapstructure:"none_expire_passord" ini:"none_expire_passord" yaml:"none_expire_passord"`
	BaseAuthUsername           string        `json:"base_auth_username" mapstructure:"base_auth_username" ini:"base_auth_username" yaml:"base_auth_username"`
	BaseAuthPassword           string        `json:"base_auth_password" mapstructure:"base_auth_password" ini:"base_auth_password" yaml:"base_auth_password"`
	AccessTokenExpireDuration  time.Duration `json:"access_token_expire_duration" mapstructure:"access_token_expire_duration" ini:"access_token_expire_duration" yaml:"access_token_expire_duration"`
	RefreshTokenExpireDuration time.Duration `json:"refresh_token_expire_duration" mapstructure:"refresh_token_expire_duration" ini:"refresh_token_expire_duration" yaml:"refresh_token_expire_duration"`
}

func (*Auth) setDefault() {
	cv.SetDefault("auth.none_expire_token", noneExpireToken)
	cv.SetDefault("auth.none_expire_username", noneExpireUser)
	cv.SetDefault("auth.none_expire_passord", noneExpirePass)
	cv.SetDefault("auth.base_auth_username", baseAuthUsername)
	cv.SetDefault("auth.base_auth_password", baseAuthPassword)
	cv.SetDefault("auth.access_token_expire_duration", "2h")
	cv.SetDefault("auth.refresh_token_expire_duration", "168h")
}

const (
	LOGGER_DIR         = "LOGGER_DIR"
	LOGGER_PREFIX      = "LOGGER_PREFIX"
	LOGGER_FILE        = "LOGGER_FILE"
	LOGGER_LEVEL       = "LOGGER_LEVEL"
	LOGGER_FORMAT      = "LOGGER_FORMAT"
	LOGGER_ENCODER     = "LOGGER_ENCODER"
	LOGGER_MAX_AGE     = "LOGGER_MAX_AGE"
	LOGGER_MAX_SIZE    = "LOGGER_MAX_SIZE"
	LOGGER_MAX_BACKUPS = "LOGGER_MAX_BACKUPS"
)

// Logger represents section "logger" for client-side or server-side configuration,
// and there is only one copy during the application entire lifetime.
type Logger struct {
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

func (*Logger) setDefault() {
	cv.SetDefault("logger.dir", "./logs")
	cv.SetDefault("logger.prefix", "")
	cv.SetDefault("logger.file", "")
	cv.SetDefault("logger.level", "info")
	cv.SetDefault("logger.format", "json")
	cv.SetDefault("logger.encoder", "json")
	cv.SetDefault("logger.max_age", 30)
	cv.SetDefault("logger.max_size", 100)
	cv.SetDefault("logger.max_backups", 1)
}

const (
	DATABASE_SLOW_QUERY_THRESHOLD = "DATABASE_SLOW_QUERY_THRESHOLD"
	DATABASE_MAX_IDLE_CONNS       = "DATABASE_MAX_IDLE_CONNS"
	DATABASE_MAX_OPEN_CONNS       = "DATABASE_MAX_OPEN_CONNS"
	DATABASE_CONN_MAX_LIFETIME    = "DATABASE_CONN_MAX_LIFETIME"
	DATABASE_CONN_MAX_IDLE_TIME   = "DATABASE_CONN_MAX_IDLE_TIME"
)

type Database struct {
	SlowQueryThreshold time.Duration `json:"slow_query_threshold" mapstructure:"slow_query_threshold" ini:"slow_query_threshold" yaml:"slow_query_threshold"`
	MaxIdleConns       int           `json:"max_idle_conns" mapstructure:"max_idle_conns" ini:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns       int           `json:"max_open_conns" mapstructure:"max_open_conns" ini:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime    time.Duration `json:"conn_max_lifetime" mapstructure:"conn_max_lifetime" ini:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime    time.Duration `json:"conn_max_idle_time" mapstructure:"conn_max_idle_time" ini:"conn_max_idle_time" yaml:"conn_max_idle_time"`
}

func (*Database) setDefault() {
	cv.SetDefault("database.slow_query_threshold", 500*time.Millisecond)
	cv.SetDefault("database.max_idle_conns", 10)
	cv.SetDefault("database.max_open_conns", 100)
	cv.SetDefault("database.conn_max_lifetime", 1*time.Hour)
	cv.SetDefault("database.conn_max_idle_time", 10*time.Minute)
}

const (
	SQLITE_PATH      = "SQLITE_PATH"
	SQLITE_DATABASE  = "SQLITE_DATABASE"
	SQLITE_IS_MEMORY = "SQLITE_IS_MEMORY"
	SQLITE_ENABLE    = "SQLITE_ENABLE"
)

type Sqlite struct {
	Path     string `json:"path" mapstructure:"path" ini:"path" yaml:"path"`
	Database string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	IsMemory bool   `json:"is_memory" mapstructure:"is_memory" ini:"is_memory" yaml:"is_memory"`
	Enable   bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Sqlite) setDefault() {
	cv.SetDefault("sqlite.path", "./data.db")
	cv.SetDefault("sqlite.database", "main")
	cv.SetDefault("sqlite.is_memory", false)
	cv.SetDefault("sqlite.enable", true)
}

const (
	POSTGRES_HOST     = "POSTGRES_HOST"
	POSTGRES_PORT     = "POSTGRES_PORT"
	POSTGRES_DATABASE = "POSTGRES_DATABASE"
	POSTGRES_USERNAME = "POSTGRES_USERNAME"
	POSTGRES_PASSWORD = "POSTGRES_PASSWORD"
	POSTGRES_SSLMODE  = "POSTGRES_SSLMODE"
	POSTGRES_TIMEZONE = "POSTGRES_TIMEZONE"
	POSTGRES_ENABLE   = "POSTGRES_ENABLE"
)

type Postgres struct {
	Host     string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port     uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Database string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	Username string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	SSLMode  string `json:"sslmode" mapstructure:"sslmode" ini:"sslmode" yaml:"sslmode"`
	TimeZone string `json:"timezone" mapstructure:"timezone" ini:"timezone" yaml:"timezone"`
	Enable   bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Postgres) setDefault() {
	cv.SetDefault("postgres.host", "127.0.0.1")
	cv.SetDefault("postgres.port", 5432)
	cv.SetDefault("postgres.database", "postgres")
	cv.SetDefault("postgres.username", "postgres")
	cv.SetDefault("postgres.password", "")
	cv.SetDefault("postgres.sslmode", "disable")
	cv.SetDefault("postgres.timezone", "UTC")
	cv.SetDefault("postgres.enable", true)
}

const (
	MYSQL_HOST     = "MYSQL_HOST"
	MYSQL_PORT     = "MYSQL_PORT"
	MYSQL_DATABASE = "MYSQL_DATABASE"
	MYSQL_USERNAME = "MYSQL_USERNAME"
	MYSQL_PASSWORD = "MYSQL_PASSWORD"
	MYSQL_CHARSET  = "MYSQL_CHARSET"
	MYSQL_ENABLE   = "MYSQL_ENABLE"
)

type MySQL struct {
	Host     string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port     uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Database string `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	Username string `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password string `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	Charset  string `json:"charset" mapstructure:"charset" ini:"charset" yaml:"charset"`
	Enable   bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*MySQL) setDefault() {
	cv.SetDefault("mysql.host", "127.0.0.1")
	cv.SetDefault("mysql.port", 3306)
	cv.SetDefault("mysql.database", "")
	cv.SetDefault("mysql.username", "root")
	cv.SetDefault("mysql.password", "toor")
	cv.SetDefault("mysql.charset", "utf8mb4")
	cv.SetDefault("mysql.enable", true)
}

const (
	CLICKHOUSE_HOST          = "CLICKHOUSE_HOST"
	CLICKHOUSE_PORT          = "CLICKHOUSE_PORT"
	CLICKHOUSE_DATABASE      = "CLICKHOUSE_DATABASE"
	CLICKHOUSE_USERNAME      = "CLICKHOUSE_USERNAME"
	CLICKHOUSE_PASSWORD      = "CLICKHOUSE_PASSWORD"
	CLICKHOUSE_DIAL_TIMEOUT  = "CLICKHOUSE_DIAL_TIMEOUT"
	CLICKHOUSE_READ_TIMEOUT  = "CLICKHOUSE_READ_TIMEOUT"
	CLICKHOUSE_WRITE_TIMEOUT = "CLICKHOUSE_WRITE_TIMEOUT"
	CLICKHOUSE_COMPRESS      = "CLICKHOUSE_COMPRESS"
	CLICKHOUSE_DEBUG         = "CLICKHOUSE_DEBUG"
	CLICKHOUSE_ENABLE        = "CLICKHOUSE_ENABLE"
)

type Clickhouse struct {
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

func (*Clickhouse) setDefault() {
	cv.SetDefault("clickhouse.host", "127.0.0.1")
	cv.SetDefault("clickhouse.port", 9000)
	cv.SetDefault("clickhouse.database", "default")
	cv.SetDefault("clickhouse.username", "default")
	cv.SetDefault("clickhouse.password", "")
	cv.SetDefault("clickhouse.dial_timeout", "5s")
	cv.SetDefault("clickhouse.read_timeout", "30s")
	cv.SetDefault("clickhouse.write_timeout", "30s")
	cv.SetDefault("clickhouse.compress", false)
	cv.SetDefault("clickhouse.debug", false)
	cv.SetDefault("clickhouse.enable", false)
}

const (
	SQLSERVER_HOST         = "SQLSERVER_HOST"
	SQLSERVER_PORT         = "SQLSERVER_PORT"
	SQLSERVER_DATABASE     = "SQLSERVER_DATABASE"
	SQLSERVER_USERNAME     = "SQLSERVER_USERNAME"
	SQLSERVER_PASSWORD     = "SQLSERVER_PASSWORD"
	SQLSERVER_ENCRYPT      = "SQLSERVER_ENCRYPT"
	SQLSERVER_TRUST_SERVER = "SQLSERVER_TRUST_SERVER"
	SQLSERVER_APP_NAME     = "SQLSERVER_APP_NAME"
	SQLSERVER_ENABLE       = "SQLSERVER_ENABLE"
)

type SQLServer struct {
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

func (*SQLServer) setDefault() {
	cv.SetDefault("sqlserver.host", "127.0.0.1")
	cv.SetDefault("sqlserver.port", 1433)
	cv.SetDefault("sqlserver.database", "")
	cv.SetDefault("sqlserver.username", "sa")
	cv.SetDefault("sqlserver.password", "")
	cv.SetDefault("sqlserver.encrypt", false)
	cv.SetDefault("sqlserver.trust_server", true)
	cv.SetDefault("sqlserver.app_name", "golib")
	cv.SetDefault("sqlserver.enable", false)
}

const (
	REDIS_ADDR         = "REDIS_ADDR"
	REDIS_ADDRS        = "REDIS_ADDRS"
	REDIS_DB           = "REDIS_DB"
	REDIS_PASSWORD     = "REDIS_PASSWORD"
	REDIS_NAMESPACE    = "REDIS_NAMESPACE"
	REDIS_POOL_SIZE    = "REDIS_POOL_SIZE"
	REDIS_EXPIRATION   = "REDIS_EXPIRATION"
	REDIS_CLUSTER_MODE = "REDIS_CLUSTER_MODE"

	REDIS_DIAL_TIMEOUT      = "REDIS_DIAL_TIMEOUT"
	REDIS_READ_TIMEOUT      = "REDIS_READ_TIMEOUT"
	REDIS_WRITE_TIMEOUT     = "REDIS_WRITE_TIMEOUT"
	REDIS_MIN_IDLE_CONNS    = "REDIS_MIN_IDLE_CONNS"
	REDIS_MAX_RETRIES       = "REDIS_MAX_RETRIES"
	REDIS_MIN_RETRY_BACKOFF = "REDIS_MIN_RETRY_BACKOFF"
	REDIS_MAX_RETRY_BACKOFF = "REDIS_MAX_RETRY_BACKOFF"

	REDIS_ENABLE_TLS           = "REDIS_ENABLE_TLS"
	REDIS_CERT_FILE            = "REDIS_CERT_FILE"
	REDIS_KEY_FILE             = "REDIS_KEY_FILE"
	REDIS_CA_FILE              = "REDIS_CA_FILE"
	REDIS_INSECURE_SKIP_VERIFY = "REDIS_INSECURE_SKIP_VERIFY"

	REDIS_ENABLE = "REDIS_ENABLE"
)

type Redis struct {
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

func (*Redis) setDefault() {
	cv.SetDefault("redis.addr", "127.0.0.1:6379")
	cv.SetDefault("redis.addrs", []string{"127.0.0.1:6379"})
	cv.SetDefault("redis.db", 0)
	cv.SetDefault("redis.password", "")
	cv.SetDefault("redis.pool_size", runtime.NumCPU())
	cv.SetDefault("redis.namespace", APP_NAME)
	cv.SetDefault("redis.expiration", 0)
	cv.SetDefault("redis.cluster_mode", false)

	cv.SetDefault("redis.dial_timeout", 0)
	cv.SetDefault("redis.read_timeout", 0)
	cv.SetDefault("redis.write_timeout", 0)
	cv.SetDefault("redis.min_idle_conns", 0)
	cv.SetDefault("redis.max_retries", 0)
	cv.SetDefault("redis.min_retry_backoff", 0)
	cv.SetDefault("redis.max_retry_backoff", 0)

	cv.SetDefault("redis.enable_tls", false)
	cv.SetDefault("redis.cert_file", "")
	cv.SetDefault("redis.key_file", 0)
	cv.SetDefault("redis.ca_file", "")
	cv.SetDefault("redis.insecure_skip_verify", false)

	cv.SetDefault("redis.enable", false)
}

const (
	ELASTICSEARCH_ADDRS                   = "ELASTICSEARCH_ADDRS"
	ELASTICSEARCH_USERNAME                = "ELASTICSEARCH_USERNAME"
	ELASTICSEARCH_PASSWORD                = "ELASTICSEARCH_PASSWORD"
	ELASTICSEARCH_CLOUD_ID                = "ELASTICSEARCH_CLOUD_ID"
	ELASTICSEARCH_API_KEY                 = "ELASTICSEARCH_API_KEY"
	ELASTICSEARCH_SERVICE_TOKEN           = "ELASTICSEARCH_SERVICE_TOKEN"
	ELASTICSEARCH_CERTIFICATE_FINGERPRINT = "ELASTICSEARCH_CERTIFICATE_FINGERPRINT"
	ELASTICSEARCH_ENABLE                  = "ELASTICSEARCH_ENABLE"
)

type Elasticsearch struct {
	Addrs                  []string `json:"addrs" mapstructure:"addrs" ini:"addrs" yaml:"addrs"`
	Username               string   `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password               string   `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	CloudID                string   `json:"cloud_id" mapstructure:"cloud_id" ini:"cloud_id" yaml:"cloud_id"`
	APIKey                 string   `json:"api_key" mapstructure:"api_key" ini:"api_key" yaml:"api_key"`
	ServiceToken           string   `json:"service_token" mapstructure:"service_token" ini:"service_token" yaml:"service_token"`
	CertificateFingerprint string   `json:"certificate_fingerprint" mapstructure:"certificate_fingerprint" ini:"certificate_fingerprint" yaml:"certificate_fingerprint"`
	Enable                 bool     `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Elasticsearch) setDefault() {
	cv.SetDefault("elasticsearch.addrs", []string{"http://127.0.0.1:9200"})
	cv.SetDefault("elasticsearch.username", "")
	cv.SetDefault("elasticsearch.password", "")
	cv.SetDefault("elasticsearch.cloud_id", "")
	cv.SetDefault("elasticsearch.api_key", "")
	cv.SetDefault("elasticsearch.service_token", "")
	cv.SetDefault("elasticsearch.certificate_fingerprint", "")
	cv.SetDefault("elasticsearch.enable", false)
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

const (
	MONGO_HOST          = "MONGO_HOST"
	MONGO_PORT          = "MONGO_PORT"
	MONGO_USERNAME      = "MONGO_USERNAME"
	MONGO_PASSWORD      = "MONGO_PASSWORD"
	MONGO_DATABASE      = "MONGO_DATABASE"
	MONGO_AUTH_SOURCE   = "MONGO_AUTH_SOURCE"
	MONGO_MAX_POOL_SIZE = "MONGO_MAX_POOL_SIZE"
	MONGO_MIN_POOL_SIZE = "MONGO_MIN_POOL_SIZE"

	MONGO_CONNECT_TIMEOUT          = "MONGO_CONNECT_TIMEOUT"
	MONGO_SERVER_SELECTION_TIMEOUT = "MONGO_SERVER_SELECTION_TIMEOUT"
	MONGO_MAX_CONN_IDLE_TIME       = "MONGO_MAX_CONN_IDLE_TIME"
	MONGO_MAX_CONNECTING           = "MONGO_MAX_CONNECTING"

	MONGO_READ_CONCERN  = "MONGO_READ_CONCERN"
	MONGO_WRITE_CONCERN = "MONGO_WRITE_CONCERN"

	MONGO_ENABLE_TLS           = "MONGO_ENABLE_TLS"
	MONGO_CERT_FILE            = "MONGO_CERT_FILE"
	MONGO_KEY_FILE             = "MONGO_KEY_FILE"
	MONGO_CA_FILE              = "MONGO_CA_FILE"
	MONGO_INSECURE_SKIP_VERIFY = "MONGO_INSECURE_SKIP_VERIFY"

	MONGO_ENABLE = "MONGO_ENABLE"
)

type Mongo struct {
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

func (*Mongo) setDefault() {
	cv.SetDefault("mongo.host", "127.0.0.1")
	cv.SetDefault("mongo.port", 27017)
	cv.SetDefault("mongo.username", "")
	cv.SetDefault("mongo.password", "")
	cv.SetDefault("mongo.database", "")
	cv.SetDefault("mongo.auth_source", "admin")
	cv.SetDefault("mongo.max_pool_size", 0)
	cv.SetDefault("mongo.min_pool_size", 0)

	cv.SetDefault("mongo.connect_timeout", 0)
	cv.SetDefault("mongo.server_selection_timeout", 0)
	cv.SetDefault("mongo.max_conn_idle_time", 0)
	cv.SetDefault("mongo.max_connecting", 0)

	cv.SetDefault("mongo.read_concern", "")
	cv.SetDefault("mongo.write_concern", "")

	cv.SetDefault("mongo.enable_tls", false)
	cv.SetDefault("mongo.cert_file", "")
	cv.SetDefault("mongo.key_file", "")
	cv.SetDefault("mongo.ca_file", "")
	cv.SetDefault("mongo.insecure_skip_verify", false)

	cv.SetDefault("mongo.enable", false)
}

const (
	LDAP_HOST          = "LDAP_HOST"
	LDAP_PORT          = "LDAP_PORT"
	LDAP_USE_SSL       = "LDAP_USE_SSL"
	LDAP_BIND_DN       = "LDAP_BIND_DN"
	LDAP_BIND_PASSWORD = "LDAP_BIND_PASSWORD"
	LDAP_BASE_DN       = "LDAP_BASE_DN"
	LDAP_SEARCH_FILTER = "LDAP_SEARCH_FILTER"
	LDAP_ENABLE        = "LDAP_ENABLE"
)

// Ldap
// For example:
// [ldap]
// host = example.cn
// port = 389
// use_ssl =  false
// bind_dn =  my_bind_dn
// bind_password = mypass
// base_dn = my_base_dn
// search_filter = (sAMAccountName=%s)
type Ldap struct {
	Host         string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port         uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	UseSsl       bool   `json:"use_ssl" mapstructure:"use_ssl" ini:"use_ssl" yaml:"use_ssl"`
	BindDN       string `json:"bind_dn" mapstructure:"bind_dn" ini:"bind_dn" yaml:"bind_dn"`
	BindPassword string `json:"bind_password" mapstructure:"bind_password" ini:"bind_password" yaml:"bind_password"`
	BaseDN       string `json:"base_dn" mapstructure:"base_dn" ini:"base_dn" yaml:"base_dn"`
	SearchFilter string `json:"search_filter" mapstructure:"search_filter" ini:"search_filter" yaml:"search_filter"`
	Enable       bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Ldap) setDefault() {
	cv.SetDefault("ldap.host", "127.0.0.1")
	cv.SetDefault("ldap.port", 389)
	cv.SetDefault("ldap.use_ssl", false)
	cv.SetDefault("ldap.bind_dn", "")
	cv.SetDefault("ldap.bind_password", "")
	cv.SetDefault("ldap.base_dn", "")
	cv.SetDefault("ldap.search_filter", "")
	cv.SetDefault("ldap.enable", false)
}

const (
	INFLUXDB_HOST   = "INFLUXDB_HOST"
	INFLUXDB_PORT   = "INFLUXDB_PORT"
	INFLUXDB_TOKEN  = "INFLUXDB_TOKEN"
	INFLUXDB_ORG    = "INFLUXDB_ORG"
	INFLUXDB_BUCKET = "INFLUXDB_BUCKET"

	INFLUXDB_BATCH_SIZE         = "INFLUXDB_BATCH_SIZE"
	INFLUXDB_FLUSH_INTERVAL     = "INFLUXDB_FLUSH_INTERVAL"
	INFLUXDB_RETRY_INTERVAL     = "INFLUXDB_RETRY_INTERVAL"
	INFLUXDB_MAX_RETRIES        = "INFLUXDB_MAX_RETRIES"
	INFLUXDB_RETRY_BUFFER_LIMIT = "INFLUXDB_RETRY_BUFFER_LIMIT"
	INFLUXDB_MAX_RETRY_INTERVAL = "INFLUXDB_MAX_RETRY_INTERVAL"
	INFLUXDB_MAX_RETRY_TIME     = "INFLUXDB_MAX_RETRY_TIME"
	INFLUXDB_EXPONENTIAL_BASE   = "INFLUXDB_EXPONENTIAL_BASE"
	INFLUXDB_PRECISION          = "INFLUXDB_PRECISION"
	INFLUXDB_USE_GZIP           = "INFLUXDB_USE_GZIP"

	INFLUXDB_ENABLE_TLS           = "INFLUXDB_ENABLE_TLS"
	INFLUXDB_CERT_FILE            = "INFLUXDB_CERT_FILE"
	INFLUXDB_KEY_FILE             = "INFLUXDB_KEY_FILE"
	INFLUXDB_CA_FILE              = "INFLUXDB_CA_FILE"
	INFLUXDB_INSECURE_SKIP_VERIFY = "INFLUXDB_INSECURE_SKIP_VERIFY"

	INFLUXDB_DEFAULT_TAGS = "INFLUXDB_DEFAULT_TAGS" // format：key1=value1,key2=value2
	INFLUXDB_APP_NAME     = "INFLUXDB_APP_NAME"

	INFLUXDB_ENABLE = "INFLUXDB_ENABLE"
)

type Influxdb struct {
	Host   string `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port   uint   `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Token  string `json:"token" mapstructure:"token" ini:"token" yaml:"token"`
	Org    string `json:"org" mapstructure:"org" ini:"org" yaml:"org"`
	Bucket string `json:"bucket" mapstructure:"bucket" ini:"bucket" yaml:"bucket"`

	// Write options
	BatchSize        uint          `json:"batch_size" mapstructure:"batch_size" ini:"batch_size" yaml:"batch_size"`
	FlushInterval    time.Duration `json:"flush_interval" mapstructure:"flush_interval" ini:"flush_interval" yaml:"flush_interval"`
	RetryInterval    time.Duration `json:"retry_interval" mapstructure:"retry_interval" ini:"retry_interval" yaml:"retry_interval"`
	MaxRetries       uint          `json:"max_retries" mapstructure:"max_retries" ini:"max_retries" yaml:"max_retries"`
	RetryBufferLimit uint          `json:"retry_buffer_limit" mapstructure:"retry_buffer_limit" ini:"retry_buffer_limit" yaml:"retry_buffer_limit"`
	MaxRetryInterval time.Duration `json:"max_retry_interval" mapstructure:"max_retry_interval" ini:"max_retry_interval" yaml:"max_retry_interval"`
	MaxRetryTime     time.Duration `json:"max_retry_time" mapstructure:"max_retry_time" ini:"max_retry_time" yaml:"max_retry_time"`
	ExponentialBase  uint          `json:"exponential_base" mapstructure:"exponential_base" ini:"exponential_base" yaml:"exponential_base"`
	Precision        time.Duration `json:"precision" mapstructure:"precision" ini:"precision" yaml:"precision"`
	UseGZip          bool          `json:"use_gzip" mapstructure:"use_gzip" ini:"use_gzip" yaml:"use_gzip"`

	// TLS configuration
	EnableTLS          bool   `json:"enable_tls" mapstructure:"enable_tls" ini:"enable_tls" yaml:"enable_tls"`
	CertFile           string `json:"cert_file" mapstructure:"cert_file" ini:"cert_file" yaml:"cert_file"`
	KeyFile            string `json:"key_file" mapstructure:"key_file" ini:"key_file" yaml:"key_file"`
	CAFile             string `json:"ca_file" mapstructure:"ca_file" ini:"ca_file" yaml:"ca_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify" ini:"insecure_skip_verify" yaml:"insecure_skip_verify"`

	// Advanced options
	DefaultTags map[string]string `json:"default_tags" mapstructure:"default_tags" ini:"default_tags" yaml:"default_tags"`
	AppName     string            `json:"app_name" mapstructure:"app_name" ini:"app_name" yaml:"app_name"`

	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Influxdb) setDefault() {
	cv.SetDefault("influxdb.host", "127.0.0.1")
	cv.SetDefault("influxdb.port", 8086)
	cv.SetDefault("influxdb.token", "")
	cv.SetDefault("influxdb.org", "")
	cv.SetDefault("influxdb.bucket", "")

	cv.SetDefault("influxdb.batch_size", 0)
	cv.SetDefault("influxdb.flush_interval", 0)
	cv.SetDefault("influxdb.retry_interval", 0)
	cv.SetDefault("influxdb.max_retries", 0)
	cv.SetDefault("influxdb.retry_buffer_limit", 0)
	cv.SetDefault("influxdb.max_retry_interval", 0)
	cv.SetDefault("influxdb.max_retry_time", 0)
	cv.SetDefault("influxdb.exponential_base", 0)
	cv.SetDefault("influxdb.precision", 0)
	cv.SetDefault("influxdb.use_gzip", false)

	cv.SetDefault("influxdb.enable_tls", false)
	cv.SetDefault("influxdb.cert_file", "")
	cv.SetDefault("influxdb.key_file", 0)
	cv.SetDefault("influxdb.ca_file", "")
	cv.SetDefault("influxdb.insecure_skip_verify", false)

	cv.SetDefault("influxdb.default_tags", nil)
	cv.SetDefault("influxdb.app_name", "")

	cv.SetDefault("influxdb.enable", false)
}

const (
	MINIO_ENDPOINT   = "MINIO_ENDPOINT"
	MINIO_REGION     = "MINIO_REGION"
	MINIO_ACCESS_KEY = "MINIO_ACCESS_KEY"
	MINIO_SECRET_KEY = "MINIO_SECRET_KEY"
	MINIO_BUCKET     = "MINIO_BUCKET"
	MINIO_USE_SSL    = "MINIO_USE_SSL"
	MINIO_ENABLE     = "MINIO_ENABLE"
)

type Minio struct {
	Endpoint  string `json:"endpoint" mapstructure:"endpoint" ini:"endpoint" yaml:"endpoint"`
	Region    string `json:"region" mapstructure:"region" ini:"region" yaml:"region"`
	AccessKey string `json:"access_key" mapstructure:"access_key" ini:"access_key" yaml:"access_key"`
	SecretKey string `json:"secret_key" mapstructure:"secret_key" ini:"secret_key" yaml:"secret_key"`
	Bucket    string `json:"bucket" mapstructure:"bucket" ini:"bucket" yaml:"bucket"`
	UseSsl    bool   `json:"use_ssl" mapstructure:"use_ssl" ini:"use_ssl" yaml:"use_ssl"`
	Enable    bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Minio) setDefault() {
	cv.SetDefault("minio.endpoint", "127.0.0.1:9000")
	cv.SetDefault("minio.region", "")
	cv.SetDefault("minio.access_key", "")
	cv.SetDefault("minio.secret_key", "")
	cv.SetDefault("minio.bucket", "")
	cv.SetDefault("minio.use_ssl", false)
	cv.SetDefault("minio.enable", false)
}

const (
	S3_ENDPOINT          = "S3_ENDPOINT"
	S3_REGION            = "S3_REGION"
	S3_ACCESS_KEY_ID     = "S3_ACCESS_KEY_ID"
	S3_SECRET_ACCESS_KEY = "S3_SECRET_ACCESS_KEY"
	S3_BUCKET            = "S3_BUCKET"
	S3_USE_SSL           = "S3_USE_SSL"
	S3_ENABLE            = "S3_ENABLE"
)

type S3 struct {
	Endpoint        string `json:"endpoint" mapstructure:"endpoint" ini:"endpoint" yaml:"endpoint"`
	Region          string `json:"region" mapstructure:"region" ini:"region" yaml:"region"`
	AccessKeyID     string `json:"access_key_id" mapstructure:"access_key_id" ini:"access_key_id" yaml:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key" mapstructure:"secret_access_key" ini:"secret_access_key" yaml:"secret_access_key"`
	Bucket          string `json:"bucket" mapstructure:"bucket" ini:"bucket" yaml:"bucket"`
	UseSsl          bool   `json:"use_ssl" mapstructure:"use_ssl" ini:"use_ssl" yaml:"use_ssl"`
	Enable          bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*S3) setDefault() {
	cv.SetDefault("s3.endpoint", "")
	cv.SetDefault("s3.region", "")
	cv.SetDefault("s3.access_key_id", "")
	cv.SetDefault("s3.secret_access_key", "")
	cv.SetDefault("s3.bucket", "")
	cv.SetDefault("s3.use_ssl", false)
	cv.SetDefault("s3.enable", false)
}

const (
	MQTT_ADDR                 = "MQTT_ADDR"
	MQTT_USERNAME             = "MQTT_USERNAME"
	MQTT_PASSWORD             = "MQTT_PASSWORD"
	MQTT_CLIENT_PREFIX        = "MQTT_CLIENT_PREFIX"
	MQTT_CONNECT_TIMEOUT      = "MQTT_CONNECT_TIMEOUT"
	MQTT_KEEPALIVE            = "MQTT_KEEPALIVE"
	MQTT_CLEAN_SESSION        = "MQTT_CLEAN_SESSION"
	MQTT_AUTO_RECONNECT       = "MQTT_AUTO_RECONNECT"
	MQTT_USE_TLS              = "MQTT_USE_TLS"
	MQTT_CERT_FILE            = "MQTT_CERT_FILE"
	MQTT_KEY_FILE             = "MQTT_KEY_FILE"
	MQTT_INSECURE_SKIP_VERIFY = "MQTT_INSECURE_SKIP_VERIFY"
	MQTT_ENABLE               = "MQTT_ENABLE"
)

type Mqtt struct {
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

func (*Mqtt) setDefault() {
	cv.SetDefault("mqtt.addr", "127.0.0.1:1883")
	cv.SetDefault("mqtt.username", "")
	cv.SetDefault("mqtt.password", "")
	cv.SetDefault("mqtt.client_prefix", "")
	cv.SetDefault("mqtt.connect_timeout", 10*time.Second)
	cv.SetDefault("mqtt.keepalive", 1*time.Minute)
	cv.SetDefault("mqtt.clean_session", true)
	cv.SetDefault("mqtt.auto_reconnect", true)
	cv.SetDefault("mqtt.use_tls", false)
	cv.SetDefault("mqtt.cert_file", "")
	cv.SetDefault("mqtt.key_file", "")
	cv.SetDefault("mqtt.insecure_skip_verify", true)
	cv.SetDefault("mqtt.enable", false)
}

const (
	NATS_ADDRS       = "NATS_ADDRS"
	NATS_CLIENT_NAME = "NATS_CLIENT_NAME"
	NATS_USERNAME    = "NATS_USERNAME"
	NATS_PASSWORD    = "NATS_PASSWORD"
	NATS_TOKEN       = "NATS_TOKEN"
	NATS_CREDENTIALS = "NATS_CREDENTIALS"
	NATS_NKEY_FILE   = "NATS_NKEY_FILE"

	NATS_MAX_RECONNECTS       = "NATS_MAX_RECONNECTS"
	NATS_RECONNECT_WAIT       = "NATS_RECONNECT_WAIT"
	NATS_RECONNECT_JITTER     = "NATS_RECONNECT_JITTER"
	NATS_RECONNECT_JITTER_TLS = "NATS_RECONNECT_JITTER_TLS"

	NATS_CONNECT_TIMEOUT       = "NATS_CONNECT_TIMEOUT"
	NATS_PING_INTERVAL         = "NATS_PING_INTERVAL"
	NATS_MAX_PINGS_OUTSTANDING = "NATS_MAX_PINGS_OUTSTANDING"

	NATS_ENABLE_TLS           = "NATS_ENABLE_TLS"
	NATS_CERT_FILE            = "NATS_CERT_FILE"
	NATS_KEY_FILE             = "NATS_KEY_FILE"
	NATS_CA_FILE              = "NATS_CA_FILE"
	NATS_INSECURE_SKIP_VERIFY = "NATS_INSECURE_SKIP_VERIFY"

	NATS_ENABLE = "NATS_ENABLE"
)

type Nats struct {
	Addrs           []string `json:"addrs" mapstructure:"addrs" ini:"addrs" yaml:"addrs"`
	ClientName      string   `json:"client_name" mapstructure:"client_name" ini:"client_name" yaml:"client_name"`
	Username        string   `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password        string   `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	Token           string   `json:"token" mapstructure:"token" ini:"token" yaml:"token"`
	CredentialsFile string   `json:"credentials_file" mapstructure:"credentials_file" ini:"credentials_file" yaml:"credentials_file"`
	NKeyFile        string   `json:"nkey_file" mapstructure:"nkey_file" ini:"nkey_file" yaml:"nkey_file"`

	MaxReconnects      int           `json:"max_reconnects" mapstructure:"max_reconnects" ini:"max_reconnects" yaml:"max_reconnects"`
	ReconnectWait      time.Duration `json:"reconnect_wait" mapstructure:"reconnect_wait" ini:"reconnect_wait" yaml:"reconnect_wait"`
	ReconnectJitter    time.Duration `json:"reconnect_jitter" mapstructure:"reconnect_jitter" ini:"reconnect_jitter" yaml:"reconnect_jitter"`
	ReconnectJitterTLS time.Duration `json:"reconnect_jitter_tls" mapstructure:"reconnect_jitter_tls" ini:"reconnect_jitter_tls" yaml:"reconnect_jitter_tls"`

	ConnectTimeout      time.Duration `json:"connect_timeout" mapstructure:"connect_timeout" ini:"connect_timeout" yaml:"connect_timeout"`
	PingInterval        time.Duration `json:"ping_interval" mapstructure:"ping_interval" ini:"ping_interval" yaml:"ping_interval"`
	MaxPingsOutstanding int           `json:"max_pings_outstanding" mapstructure:"max_pings_outstanding" ini:"max_pings_outstanding" yaml:"max_pings_outstanding"`

	EnableTLS          bool   `json:"enable_tls" mapstructure:"enable_tls" ini:"enable_tls" yaml:"enable_tls"`
	CertFile           string `json:"cert_file" mapstructure:"cert_file" ini:"cert_file" yaml:"cert_file"`
	KeyFile            string `json:"key_file" mapstructure:"key_file" ini:"key_file" yaml:"key_file"`
	CAFile             string `json:"ca_file" mapstructure:"ca_file" ini:"ca_file" yaml:"ca_file"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify" mapstructure:"insecure_skip_verify" ini:"insecure_skip_verify" yaml:"insecure_skip_verify"`

	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Nats) setDefault() {
	cv.SetDefault("nats.addrs", []string{"nats://127.0.0.1:4222"})
	cv.SetDefault("nats.client_name", "")
	cv.SetDefault("nats.username", "")
	cv.SetDefault("nats.password", "")
	cv.SetDefault("nats.token", "")
	cv.SetDefault("nats.credentials_file", "")
	cv.SetDefault("nats.nkey_file", "")

	cv.SetDefault("nats.max_reconnects", 5)
	cv.SetDefault("nats.reconnect_wait", 1*time.Second)
	cv.SetDefault("nats.reconnect_jitter", 0)
	cv.SetDefault("nats.reconnect_jitter_tls", 0)

	cv.SetDefault("nats.connect_timeout", 2*time.Second)
	cv.SetDefault("nats.ping_interval", 2*time.Minute)
	cv.SetDefault("nats.max_pings_outstanding", 2)

	cv.SetDefault("nats.enable_tls", false)
	cv.SetDefault("nats.cert_file", "")
	cv.SetDefault("nats.key_file", "")
	cv.SetDefault("nats.ca_file", "")
	cv.SetDefault("nats.insecure_skip_verify", false)

	cv.SetDefault("nats.enable", false)
}

const (
	FEISHU_APP_ID     = "FEISHU_APP_ID"
	FEISHU_APP_SECRET = "FEISHU_APP_SECRET"
	FEISHU_ENABLE     = "FEISHU_ENABLE"
)

type Feishu struct {
	AppID     string `json:"app_id" mapstructure:"app_id" ini:"app_id" yaml:"app_id"`
	AppSecret string `json:"app_secret" mapstructure:"app_secret" ini:"app_secret" yaml:"app_secret"`
	Enable    bool   `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Feishu) setDefault() {
	cv.SetDefault("feishu.app_id", "")
	cv.SetDefault("feishu.app_secret", "")
	cv.SetDefault("feishu.enable", false)
}

const (
	DEBUG_ENABLE_STATSVIZ = "DEBUG_ENABLE_STATSVIZ"
	DEBUG_ENABLE_PPROF    = "DEBUG_ENABLE_PPROF"
	DEBUG_ENABLE_GOPS     = "DEBUG_ENABLE_GOPS"

	DEBUG_STATSVIZ_LISTEN = "DEBUG_STATSVIZ_LISTEN"
	DEBUG_PPROF_LISTEN    = "DEBUG_PPROF_LISTEN"
	DEBUG_GOPS_LISTEN     = "DEBUG_GOPS_LISTEN"

	DEBUG_STATSVIZ_PORT = "DEBUG_STATSVIZ_PORT"
	DEBUG_PPROF_PORT    = "DEBUG_PPROF_PORT"
	DEBUG_GOPS_PORT     = "DEBUG_GOPS_PORT"
)

type Debug struct {
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

func (*Debug) setDefault() {
	cv.SetDefault("debug.enable_statsviz", false)
	cv.SetDefault("debug.statsviz_listen", "127.0.0.1")
	cv.SetDefault("debug.statsviz_port", 10000)

	cv.SetDefault("debug.enable_pprof", false)
	cv.SetDefault("debug.pprof_listen", "127.0.0.1")
	cv.SetDefault("debug.gops_listen", "127.0.0.1")

	cv.SetDefault("debug.enable_gops", false)
	cv.SetDefault("debug.pprof_port", 10001)
	cv.SetDefault("debug.gops_port", 10002)
}
