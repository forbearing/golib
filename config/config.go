package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/spf13/viper"
)

const (
	APP_NAME         = "myproject"
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
	mu          sync.Mutex
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

	if App.RedisConfig.Expiration < time.Minute || App.RedisConfig.Expiration > 24*time.Hour {
		App.RedisConfig.Expiration = 8 * time.Hour
	}

	return nil
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

	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.pool_size", runtime.NumCPU())
	viper.SetDefault("redis.namespace", "myproject")
	viper.SetDefault("redis.expiration", "8h")
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
	viper.SetDefault("mongo.max_pool_size", 100)
	viper.SetDefault("mongo.min_pool_size", 0)
	viper.SetDefault("mongo.connect_timeout", 30*time.Second)
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
	FeishuConfig        `json:"feishu" mapstructure:"feishu" ini:"feishu" yaml:"feishu"`
	MqttConfig          `json:"mqtt" mapstructure:"mqtt" ini:"mqtt" yaml:"mqtt"`
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
	Host       string        `json:"host" mapstructure:"host"`
	Port       uint          `json:"port" mapstructure:"port"`
	DB         int           `json:"db" mapstructure:"db"`
	Password   string        `json:"password" mapstructure:"password"`
	PoolSize   int           `json:"pool_size" mapstructure:"pool_size"`
	Namespace  string        `json:"namespace" mapstructure:"namespace" ini:"namespace" yaml:"namespace"`
	Expiration time.Duration `json:"expiration" mapstructure:"expiration" ini:"expiration" yaml:"expiration"`
	Enable     bool          `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
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
	Host           string        `json:"host" mapstructure:"host" ini:"host" yaml:"host"`
	Port           int           `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Username       string        `json:"username" mapstructure:"username" ini:"username" yaml:"username"`
	Password       string        `json:"password" mapstructure:"password" ini:"password" yaml:"password"`
	Database       string        `json:"database" mapstructure:"database" ini:"database" yaml:"database"`
	AuthSource     string        `json:"auth_source" mapstructure:"auth_source" ini:"auth_source" yaml:"auth_source"`
	MaxPoolSize    uint64        `json:"max_pool_size" mapstructure:"max_pool_size" ini:"max_pool_size" yaml:"max_pool_size"`
	MinPoolSize    uint64        `json:"min_pool_size" mapstructure:"min_pool_size" ini:"min_pool_size" yaml:"min_pool_size"`
	ConnectTimeout time.Duration `json:"connect_timeout" mapstructure:"connect_timeout" ini:"connect_timeout" yaml:"connect_timeout"`
	Enable         bool          `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

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
