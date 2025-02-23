package config

const (
	SERVER_DOMAIN               = "SERVER_DOMAIN"
	SERVER_MODE                 = "SERVER_MODE"
	SERVER_LISTEN               = "SERVER_LISTEN"
	SERVER_PORT                 = "SERVER_PORT"
	SERVER_DB                   = "SERVER_DB"
	SERVER_SLOW_QUERY_THRESHOLD = "SERVER_SLOW_QUERY_THRESHOLD"
	SERVER_ENABLE_RBAC          = "SERVER_ENABLE_RBAC"
)

const (
	AUTH_NONE_EXPIRE_TOKEN             = "AUTH_NONE_EXPIRE_TOKEN"
	AUTH_NONE_EXPIRE_USER              = "AUTH_NONE_EXPIRE_USER"
	AUTH_NONE_EXPIRE_PASS              = "AUTH_NONE_EXPIRE_PASS"
	AUTH_BASE_AUTH_USERNAME            = "AUTH_BASE_AUTH_USERNAME"
	AUTH_BASE_AUTH_PASSWORD            = "AUTH_BASE_AUTH_PASSWORD"
	AUTH_ACCESS_TOKEN_EXPIRE_DURATION  = "AUTH_TOKEN_ACCESS_EXPIRE_DURATION"
	AUTH_REFRESH_TOKEN_EXPIRE_DURATION = "AUTH_TOKEN_REFRESH_EXPIRE_DURATION"
)

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

const (
	DATABASE_SLOW_QUERY_THRESHOLD = "DATABASE_SLOW_QUERY_THRESHOLD"
	DATABASE_MAX_IDLE_CONNS       = "DATABASE_MAX_IDLE_CONNS"
	DATABASE_MAX_OPEN_CONNS       = "DATABASE_MAX_OPEN_CONNS"
	DATABASE_CONN_MAX_LIFETIME    = "DATABASE_CONN_MAX_LIFETIME"
	DATABASE_CONN_MAX_IDLE_TIME   = "DATABASE_CONN_MAX_IDLE_TIME"
)

const (
	SQLITE_PATH      = "SQLITE_PATH"
	SQLITE_DATABASE  = "SQLITE_DATABASE"
	SQLITE_IS_MEMORY = "SQLITE_IS_MEMORY"
	SQLITE_ENABLE    = "SQLITE_ENABLE"
)

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

const (
	MYSQL_HOST     = "MYSQL_HOST"
	MYSQL_PORT     = "MYSQL_PORT"
	MYSQL_DATABASE = "MYSQL_DATABASE"
	MYSQL_USERNAME = "MYSQL_USERNAME"
	MYSQL_PASSWORD = "MYSQL_PASSWORD"
	MYSQL_CHARSET  = "MYSQL_CHARSET"
	MYSQL_ENABLE   = "MYSQL_ENABLE"
)

const (
	REDIS_HOST       = "REDIS_HOST"
	REDIS_PORT       = "REDIS_PORT"
	REDIS_DB         = "REDIS_DB"
	REDIS_PASSWORD   = "REDIS_PASSWORD"
	REDIS_POOL_SIZE  = "REDIS_POOL_SIZE"
	REDIS_NAMESPACE  = "REDIS_NAMESPACE"
	REDIS_EXPIRATION = "REDIS_EXPIRATION"
	REDIS_ENABLE     = "REDIS_ENABLE"
)

const (
	ELASTICSEARCH_HOSTS                   = "ELASTICSEARCH_HOSTS"
	ELASTICSEARCH_USERNAME                = "ELASTICSEARCH_USERNAME"
	ELASTICSEARCH_PASSWORD                = "ELASTICSEARCH_PASSWORD"
	ELASTICSEARCH_CLOUD_ID                = "ELASTICSEARCH_CLOUD_ID"
	ELASTICSEARCH_API_KEY                 = "ELASTICSEARCH_API_KEY"
	ELASTICSEARCH_SERVICE_TOKEN           = "ELASTICSEARCH_SERVICE_TOKEN"
	ELASTICSEARCH_CERTIFICATE_FINGERPRINT = "ELASTICSEARCH_CERTIFICATE_FINGERPRINT"
	ELASTICSEARCH_ENABLE                  = "ELASTICSEARCH_ENABLE"
)

const (
	MONGO_HOST            = "MONGO_HOST"
	MONGO_PORT            = "MONGO_PORT"
	MONGO_USERNAME        = "MONGO_USERNAME"
	MONGO_PASSWORD        = "MONGO_PASSWORD"
	MONGO_DATABASE        = "MONGO_DATABASE"
	MONGO_AUTH_SOURCE     = "MONGO_AUTH_SOURCE"
	MONGO_MAX_POOL_SIZE   = "MONGO_MAX_POOL_SIZE"
	MONGO_MIN_POOL_SIZE   = "MONGO_MIN_POOL_SIZE"
	MONGO_CONNECT_TIMEOUT = "MONGO_CONNECT_TIMEOUT"
	MONGO_ENABLE          = "MONGO_ENABLE"
)

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

const (
	INFLUXDB_HOST           = "INFLUXDB_HOST"
	INFLUXDB_PORT           = "INFLUXDB_PORT"
	INFLUXDB_ADMIN_PASSWORD = "INFLUXDB_ADMIN_PASSWORD"
	INFLUXDB_ADMIN_TOKEN    = "INFLUXDB_ADMIN_TOKEN"
	INFLUXDB_ADMIN_ORG      = "INFLUXDB_ADMIN_ORG"
	INFLUXDB_BUCKET         = "INFLUXDB_BUCKET"
	INFLUXDB_WRITE_INTERVAL = "INFLUXDB_WRITE_INTERVAL"
	INFLUXDB_ENABLE         = "INFLUXDB_ENABLE"
)

const (
	MINIO_ENDPOINT   = "MINIO_ENDPOINT"
	MINIO_REGION     = "MINIO_REGION"
	MINIO_ACCESS_KEY = "MINIO_ACCESS_KEY"
	MINIO_SECRET_KEY = "MINIO_SECRET_KEY"
	MINIO_BUCKET     = "MINIO_BUCKET"
	MINIO_USE_SSL    = "MINIO_USE_SSL"
	MINIO_ENABLE     = "MINIO_ENABLE"
)

const (
	S3_ENDPOINT          = "S3_ENDPOINT"
	S3_REGION            = "S3_REGION"
	S3_ACCESS_KEY_ID     = "S3_ACCESS_KEY_ID"
	S3_SECRET_ACCESS_KEY = "S3_SECRET_ACCESS_KEY"
	S3_BUCKET            = "S3_BUCKET"
	S3_USE_SSL           = "S3_USE_SSL"
	S3_ENABLE            = "S3_ENABLE"
)

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

const (
	FEISHU_APP_ID         = "FEISHU_APP_ID"
	FEISHU_APP_SECRET     = "FEISHU_APP_SECRET"
	FEISHU_MSG_APP_ID     = "FEISHU_MSG_APP_ID"
	FEISHU_MSG_APP_SECRET = "FEISHU_MSG_APP_SECRET"
	FEISHU_ENABLE         = "FEISHU_ENABLE"
)
