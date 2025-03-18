package config

import "time"

type (
	Mode string
)

const (
	ModeProd = "prod"
	ModeStg  = "stg"
	ModeDev  = "dev"
)

const (
	SERVER_DOMAIN      = "SERVER_DOMAIN"
	SERVER_MODE        = "SERVER_MODE"
	SERVER_LISTEN      = "SERVER_LISTEN"
	SERVER_PORT        = "SERVER_PORT"
	SERVER_ENABLE_RBAC = "SERVER_ENABLE_RBAC"

	// Circuit breaker related environment variables
	SERVER_CIRCUIT_BREAKER_NAME         = "SERVER_CIRCUIT_BREAKER_NAME"
	SERVER_CIRCUIT_BREAKER_MAX_REQUESTS = "SERVER_CIRCUIT_BREAKER_MAX_REQUESTS"
	SERVER_CIRCUIT_BREAKER_INTERVAL     = "SERVER_CIRCUIT_BREAKER_INTERVAL"
	SERVER_CIRCUIT_BREAKER_TIMEOUT      = "SERVER_CIRCUIT_BREAKER_TIMEOUT"
	SERVER_CIRCUIT_BREAKER_FAILURE_RATE = "SERVER_CIRCUIT_BREAKER_FAILURE_RATE"
	SERVER_CIRCUIT_BREAKER_MIN_REQUESTS = "SERVER_CIRCUIT_BREAKER_MIN_REQUESTS"
	SERVER_CIRCUIT_BREAKER_ENABLE       = "SERVER_CIRCUIT_BREAKER_ENABLE"
)

type Server struct {
	Mode       Mode   `json:"mode" mapstructure:"mode" ini:"mode" yaml:"mode"`
	Listen     string `json:"listen" mapstructure:"listen" ini:"listen" yaml:"listen"`
	Port       int    `json:"port" mapstructure:"port" ini:"port" yaml:"port"`
	Domain     string `json:"domain" mapstructure:"domain" ini:"domain" yaml:"domain"`
	EnableRBAC bool   `json:"enable_rbac" mapstructure:"enable_rbac" ini:"enable_rbac" yaml:"enable_rbac"`

	CircuitBreaker CircuitBreaker `json:"circuit_breaker" mapstructure:"circuit_breaker" ini:"circuit_breaker" yaml:"circuit_breaker"`
}

type CircuitBreaker struct {
	Name        string        `json:"name" mapstructure:"name" ini:"name" yaml:"name"`
	MaxRequests uint32        `json:"max_requests" mapstructure:"max_requests" ini:"max_requests" yaml:"max_requests"`
	Interval    time.Duration `json:"interval" mapstructure:"interval" ini:"interval" yaml:"interval"`
	Timeout     time.Duration `json:"timeout" mapstructure:"timeout" ini:"timeout" yaml:"timeout"`
	FailureRate float64       `json:"failure_rate" mapstructure:"failure_rate" ini:"failure_rate" yaml:"failure_rate"`
	MinRequests uint32        `json:"min_requests" mapstructure:"min_requests" ini:"min_requests" yaml:"min_requests"`
	Enable      bool          `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`
}

func (*Server) setDefault() {
	cv.SetDefault("server.mode", ModeDev)
	cv.SetDefault("server.listen", "")
	cv.SetDefault("server.port", 8080)
	cv.SetDefault("server.domain", "")
	cv.SetDefault("server.enable_rbac", false)

	// Circuit breaker defaults
	cv.SetDefault("server.circuit_breaker.name", "backend-server")
	cv.SetDefault("server.circuit_breaker.max_requests", uint32(100))
	cv.SetDefault("server.circuit_breaker.interval", 10*time.Second)
	cv.SetDefault("server.circuit_breaker.timeout", 30*time.Second)
	cv.SetDefault("server.circuit_breaker.failure_rate", 0.5)
	cv.SetDefault("server.circuit_breaker.min_requests", uint32(10))
	cv.SetDefault("server.circuit_breaker.enable", true)
}
