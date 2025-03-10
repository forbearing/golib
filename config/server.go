package config

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
