package config

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
