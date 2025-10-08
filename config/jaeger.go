package config

import (
	"time"

	"github.com/forbearing/gst/types/consts"
)

const (
	JAEGER_ENABLE                  = "JAEGER_ENABLE"
	JAEGER_SERVICE_NAME            = "JAEGER_SERVICE_NAME"
	JAEGER_EXPORTER_TYPE           = "JAEGER_EXPORTER_TYPE"
	JAEGER_OTLP_ENDPOINT           = "JAEGER_OTLP_ENDPOINT"
	JAEGER_OTLP_HEADERS            = "JAEGER_OTLP_HEADERS"
	JAEGER_OTLP_INSECURE           = "JAEGER_OTLP_INSECURE"
	JAEGER_SAMPLER_TYPE            = "JAEGER_SAMPLER_TYPE"
	JAEGER_SAMPLER_PARAM           = "JAEGER_SAMPLER_PARAM"
	JAEGER_LOG_SPANS               = "JAEGER_LOG_SPANS"
	JAEGER_MAX_TAG_VALUE_LEN       = "JAEGER_MAX_TAG_VALUE_LEN"
	JAEGER_BUFFER_FLUSH_INTERVAL   = "JAEGER_BUFFER_FLUSH_INTERVAL"
	JAEGER_REPORTER_QUEUE_SIZE     = "JAEGER_REPORTER_QUEUE_SIZE"
	JAEGER_REPORTER_FLUSH_INTERVAL = "JAEGER_REPORTER_FLUSH_INTERVAL"
)

type ExportType string

const (
	ExportTypeOtlpHttp ExportType = "otlp-http"
	ExportTypeOtlpGrpc ExportType = "otlp-grpc"
)

type SamplerType string

const (
	SamplerTypeConst         SamplerType = "const"
	SamplerTypeProbabilistic SamplerType = "probabilistic"
	SamplerTypeRateLimiting  SamplerType = "ratelimiting"
)

type Jaeger struct {
	// Enable controls whether Jaeger tracing is enabled
	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`

	// ServiceName is the name of the service for tracing
	ServiceName string `json:"service_name" mapstructure:"service_name" ini:"service_name" yaml:"service_name"`

	// ExporterType defines the exporter type (otlp-http, otlp-grpc)
	// Note: "jaeger" exporter is deprecated since OpenTelemetry dropped support in July 2023
	// Jaeger officially accepts and recommends using OTLP instead
	ExporterType ExportType `json:"exporter_type" mapstructure:"exporter_type" ini:"exporter_type" yaml:"exporter_type"`

	// OTLPEndpoint is the OTLP endpoint for HTTP/gRPC
	OTLPEndpoint string `json:"otlp_endpoint" mapstructure:"otlp_endpoint" ini:"otlp_endpoint" yaml:"otlp_endpoint"`

	// OTLPHeaders are the headers to send with OTLP requests
	OTLPHeaders map[string]string `json:"otlp_headers" mapstructure:"otlp_headers" ini:"otlp_headers" yaml:"otlp_headers"`

	// OTLPInsecure controls whether to use insecure connection for OTLP
	OTLPInsecure bool `json:"otlp_insecure" mapstructure:"otlp_insecure" ini:"otlp_insecure" yaml:"otlp_insecure"`

	// SamplerType defines the sampling strategy (const, probabilistic, ratelimiting, remote)
	SamplerType SamplerType `json:"sampler_type" mapstructure:"sampler_type" ini:"sampler_type" yaml:"sampler_type"`

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

func (j *Jaeger) setDefault() {
	cv.SetDefault("jaeger.enable", false)
	cv.SetDefault("jaeger.service_name", consts.FrameworkName)
	cv.SetDefault("jaeger.exporter_type", ExportTypeOtlpHttp)
	cv.SetDefault("jaeger.otlp_endpoint", "localhost:4318") // Default OTLP HTTP endpoint
	cv.SetDefault("jaeger.otlp_headers", map[string]string{})
	cv.SetDefault("jaeger.otlp_insecure", true) // Default to insecure for local development
	cv.SetDefault("jaeger.sampler_type", SamplerTypeConst)
	cv.SetDefault("jaeger.sampler_param", 1.0)
	cv.SetDefault("jaeger.log_spans", false)
	cv.SetDefault("jaeger.max_tag_value_len", 256)
	cv.SetDefault("jaeger.buffer_flush_interval", time.Second)
	cv.SetDefault("jaeger.reporter_queue_size", 100)
	cv.SetDefault("jaeger.reporter_flush_interval", time.Second)
}
