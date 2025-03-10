package config

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
