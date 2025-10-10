package config

import (
	"time"

	"github.com/forbearing/gst/types/consts"
)

const (
	OTEL_ENABLE                  = "OTEL_ENABLE"                  //nolint:staticcheck
	OTEL_SERVICE_NAME            = "OTEL_SERVICE_NAME"            //nolint:staticcheck
	OTEL_EXPORTER_TYPE           = "OTEL_EXPORTER_TYPE"           //nolint:staticcheck
	OTEL_OTLP_ENDPOINT           = "OTEL_OTLP_ENDPOINT"           //nolint:staticcheck
	OTEL_OTLP_HEADERS            = "OTEL_OTLP_HEADERS"            //nolint:staticcheck
	OTEL_OTLP_INSECURE           = "OTEL_OTLP_INSECURE"           //nolint:staticcheck
	OTEL_SAMPLER_TYPE            = "OTEL_SAMPLER_TYPE"            //nolint:staticcheck
	OTEL_SAMPLER_PARAM           = "OTEL_SAMPLER_PARAM"           //nolint:staticcheck
	OTEL_LOG_SPANS               = "OTEL_LOG_SPANS"               //nolint:staticcheck
	OTEL_MAX_TAG_VALUE_LEN       = "OTEL_MAX_TAG_VALUE_LEN"       //nolint:staticcheck
	OTEL_BUFFER_FLUSH_INTERVAL   = "OTEL_BUFFER_FLUSH_INTERVAL"   //nolint:staticcheck
	OTEL_REPORTER_QUEUE_SIZE     = "OTEL_REPORTER_QUEUE_SIZE"     //nolint:staticcheck
	OTEL_REPORTER_FLUSH_INTERVAL = "OTEL_REPORTER_FLUSH_INTERVAL" //nolint:staticcheck
)

type ExportType string

const (
	ExportTypeOtlpHTTP ExportType = "otlp-http"
	ExportTypeOtlpGRPC ExportType = "otlp-grpc"
)

type SamplerType string

const (
	SamplerTypeConst         SamplerType = "const"
	SamplerTypeProbabilistic SamplerType = "probabilistic"
	SamplerTypeRateLimiting  SamplerType = "ratelimiting"
)

// OTEL represents OpenTelemetry tracing configuration using OTLP exporters.
// This configuration supports sending traces to Jaeger, Uptrace, or other OTLP-compatible backends.
type OTEL struct {
	// Enable controls whether OpenTelemetry tracing is enabled
	Enable bool `json:"enable" mapstructure:"enable" ini:"enable" yaml:"enable"`

	// ServiceName is the name of the service for tracing
	ServiceName string `json:"service_name" mapstructure:"service_name" ini:"service_name" yaml:"service_name"`

	// ExporterType defines the exporter type (otlp-http, otlp-grpc)
	// Note: "jaeger" exporter is deprecated since OpenTelemetry dropped support in July 2023
	// Jaeger officially accepts and recommends using OTLP instead
	ExporterType ExportType `json:"exporter_type" mapstructure:"exporter_type" ini:"exporter_type" yaml:"exporter_type"`

	// OTLPEndpoint is the OTLP endpoint URL (e.g., http://localhost:4318/v1/traces for HTTP, localhost:4317 for gRPC)
	OTLPEndpoint string `json:"otlp_endpoint" mapstructure:"otlp_endpoint" ini:"otlp_endpoint" yaml:"otlp_endpoint"`

	// OTLPHeaders are additional headers to send with OTLP requests
	OTLPHeaders map[string]string `json:"otlp_headers" mapstructure:"otlp_headers" ini:"otlp_headers" yaml:"otlp_headers"`

	// OTLPInsecure controls whether to use insecure connection for OTLP
	OTLPInsecure bool `json:"otlp_insecure" mapstructure:"otlp_insecure" ini:"otlp_insecure" yaml:"otlp_insecure"`

	// SamplerType defines the sampling strategy
	SamplerType SamplerType `json:"sampler_type" mapstructure:"sampler_type" ini:"sampler_type" yaml:"sampler_type"`

	// SamplerParam is the parameter for the sampler (e.g., sampling rate for probabilistic)
	SamplerParam float64 `json:"sampler_param" mapstructure:"sampler_param" ini:"sampler_param" yaml:"sampler_param"`

	// LogSpans controls whether to log spans
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

func (o *OTEL) setDefault() {
	if o.ServiceName == "" {
		o.ServiceName = consts.FrameworkName
	}
	if o.ExporterType == "" {
		o.ExporterType = ExportTypeOtlpHTTP
	}
	if o.OTLPEndpoint == "" {
		o.OTLPEndpoint = "http://localhost:4318/v1/traces"
	}
	if o.SamplerType == "" {
		o.SamplerType = SamplerTypeConst
	}
	if o.SamplerParam == 0 {
		o.SamplerParam = 1
	}
}
