package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

	inited  bool
	tempdir string
	mu      sync.RWMutex
	cv      = viper.New()
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

type Config struct {
	Server        `json:"server" mapstructure:"server" ini:"server" yaml:"server"`
	Grpc          `json:"grpc" mapstructure:"grpc" ini:"grpc" yaml:"grpc"`
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
	Etcd          `json:"etcd" mapstructure:"etcd" ini:"etcd" yaml:"etcd"`
	Cassandra     `json:"cassandra" mapstructure:"cassandra" ini:"cassandra" yaml:"cassandra"`
	Feishu        `json:"feishu" mapstructure:"feishu" ini:"feishu" yaml:"feishu"`
	Debug         `json:"debug" mapstructure:"debug" ini:"debug" yaml:"debug"`
}

// setDefault will set config default value
func (*Config) setDefault() {
	new(Server).setDefault()
	new(Grpc).setDefault()
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
	new(Etcd).setDefault()
	new(Cassandra).setDefault()
	new(Feishu).setDefault()
	new(Debug).setDefault()
}

// Init initializes the application configuration
//
// Configuration priority (from highest to lowest):
// 1. Environment variables
// 2. Configuration file
// 3. Default values
func Init() (err error) {
	cv.AutomaticEnv()
	cv.AllowEmptyEnv(true)
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
			if tempdir, err = os.MkdirTemp("", "golib_"); err != nil {
				return errors.Wrap(err, "failed to create temp dir")
			}
			// logger not initialized using fmt.Println instead.
			fmt.Fprintf(os.Stdout, "create temp dir: %s\n", tempdir)
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

func Clean() {
	if err := os.RemoveAll(tempdir); err != nil {
		zap.S().Errorw("failed to remove temp dir", "error", err, "dir", tempdir)
	} else {
		zap.S().Infow("successfully remove temp dir", "dir", tempdir)
	}
}

func Tempdir() string {
	return tempdir
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
	for i := range src.NumField() {
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
