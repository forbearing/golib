package config

import (
	"time"

	"github.com/forbearing/golib/types/consts"
)

const (
	JAEGER_ENABLE                  = "JAEGER_ENABLE"
	JAEGER_SERVICE_NAME            = "JAEGER_SERVICE_NAME"
	JAEGER_AGENT_ENDPOINT          = "JAEGER_AGENT_ENDPOINT"
	JAEGER_COLLECTOR_URL           = "JAEGER_COLLECTOR_URL"
	JAEGER_SAMPLER_TYPE            = "JAEGER_SAMPLER_TYPE"
	JAEGER_SAMPLER_PARAM           = "JAEGER_SAMPLER_PARAM"
	JAEGER_LOG_SPANS               = "JAEGER_LOG_SPANS"
	JAEGER_MAX_TAG_VALUE_LEN       = "JAEGER_MAX_TAG_VALUE_LEN"
	JAEGER_BUFFER_FLUSH_INTERVAL   = "JAEGER_BUFFER_FLUSH_INTERVAL"
	JAEGER_REPORTER_QUEUE_SIZE     = "JAEGER_REPORTER_QUEUE_SIZE"
	JAEGER_REPORTER_FLUSH_INTERVAL = "JAEGER_REPORTER_FLUSH_INTERVAL"
)

type Jaeger struct {
	// Enable controls whether Jaeger tracing is enabled
	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`

	// ServiceName is the name of the service for tracing
	ServiceName string `json:"service_name" mapstructure:"service_name" ini:"service_name" yaml:"service_name"`

	// AgentEndpoint is the Jaeger agent endpoint (UDP)
	AgentEndpoint string `json:"agent_endpoint" mapstructure:"agent_endpoint" ini:"agent_endpoint" yaml:"agent_endpoint"`

	// CollectorURL is the Jaeger collector URL (HTTP)
	CollectorURL string `json:"collector_url" mapstructure:"collector_url" ini:"collector_url" yaml:"collector_url"`

	// SamplerType defines the sampling strategy (const, probabilistic, ratelimiting, remote)
	SamplerType string `json:"sampler_type" mapstructure:"sampler_type" ini:"sampler_type" yaml:"sampler_type"`

	// SamplerParam is the parameter for the sampling strategy
	SamplerParam float64 `json:"sampler_param" mapstructure:"sampler_param" ini:"sampler_param" yaml:"sampler_param"`

	// LogSpans controls whether to log spans to the logger
	LogSpans bool `json:"log_spans" mapstructure:"log_spans" ini:"log_spans" yaml:"log_spans"`

	// MaxTagValueLen is the maximum length of tag values
	MaxTagValueLen int `json:"max_tag_value_len" mapstructure:"max_tag_value_len" ini:"max_tag_value_len" yaml:"max_tag_value_len"`

	// BufferFlushInterval is the interval for flushing the buffer
	BufferFlushInterval time.Duration `json:"buffer_flush_interval" mapstructure:"buffer_flush_interval" ini:"buffer_flush_interval" yaml:"buffer_flush_interval"`

	// ReporterQueueSize is the size of the reporter queue
	ReporterQueueSize int `json:"reporter_queue_size" mapstructure:"reporter_queue_size" ini:"reporter_queue_size" yaml:"reporter_queue_size"`

	// ReporterFlushInterval is the interval for flushing the reporter
	ReporterFlushInterval time.Duration `json:"reporter_flush_interval" mapstructure:"reporter_flush_interval" ini:"reporter_flush_interval" yaml:"reporter_flush_interval"`
}

func (*Jaeger) setDefault() {
	cv.SetDefault("jaeger.enable", false)
	cv.SetDefault("jaeger.service_name", consts.FrameworkName)
	cv.SetDefault("jaeger.agent_endpoint", "localhost:6831")
	cv.SetDefault("jaeger.collector_url", "http://localhost:14268/api/traces")
	cv.SetDefault("jaeger.sampler_type", "const")
	cv.SetDefault("jaeger.sampler_param", 1.0)
	cv.SetDefault("jaeger.log_spans", false)
	cv.SetDefault("jaeger.max_tag_value_len", 256)
	cv.SetDefault("jaeger.buffer_flush_interval", time.Second)
	cv.SetDefault("jaeger.reporter_queue_size", 100)
	cv.SetDefault("jaeger.reporter_flush_interval", time.Second)
}
