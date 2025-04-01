package config

import "time"

const (
	MEMCACHED_SERVERS        = "MEMCACHED_SERVERS"
	MEMCACHED_MAX_IDLE_CONNS = "MEMCACHED_MAX_IDLE_CONNS"
	MEMCACHED_TIMEOUT        = "MEMCACHED_TIMEOUT"
	MEMCACHED_MAX_CACHE_SIZE = "MEMCACHED_MAX_CACHE_SIZE"
	MEMCACHED_ENABLE         = "MEMCACHED_ENABLE"
)

type Memcached struct {
	Servers      []string      `json:"servers" mapstructure:"servers" ini:"servers" yaml:"servers"`
	MaxIdleConns int           `json:"max_idle_conns" mapstructure:"max_idle_conns" ini:"max_idle_conns" yaml:"max_idle_conns"`
	Timeout      time.Duration `json:"timeout" mapstructure:"timeout" ini:"timeout" yaml:"timeout"`
	MaxCacheSize int           `json:"max_cache_size" mapstructure:"max_cache_size" ini:"max_cache_size" yaml:"max_cache_size"`

	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Memcached) setDefault() {
	cv.SetDefault("memcached.servers", "127.0.0.1:11211")
	cv.SetDefault("memcached.max_idle_conns", 2)
	cv.SetDefault("memcached.timeout", 100*time.Millisecond)
	cv.SetDefault("memcached.max_cache_size", 0) // 0 is unlimited
	cv.SetDefault("memcached.enable", false)
}
