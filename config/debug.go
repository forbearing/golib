package config

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
