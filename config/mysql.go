package config

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
	cv.SetDefault("mysql.password", "")
	cv.SetDefault("mysql.charset", "utf8mb4")
	cv.SetDefault("mysql.enable", true)
}
