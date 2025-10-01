package jaeger

import (
	"context"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/gst/config"
	"github.com/forbearing/gst/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
	mu             sync.Mutex
	initialized    bool

	ErrJaegerIsDisabled = errors.New("jaeger is disabled")
)

// Init initializes the Jaeger tracer
func Init() error {
	cfg := config.App.Jaeger
	if !cfg.Enable {
		logger.Jaeger.Info("jaeger tracing is disabled")
		return nil
	}

	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return nil
	}

	// Create Jaeger exporter
	exporter, err := createJaegerExporter(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create jaeger exporter")
	}

	// Create resource
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create resource")
	}

	// Create sampler
	sampler := createSampler(cfg)

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(cfg.BufferFlushInterval),
			sdktrace.WithMaxExportBatchSize(cfg.ReporterQueueSize),
		),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer
	tracer = otel.Tracer(cfg.ServiceName)

	// Store tracer provider for cleanup
	tracerProvider = tp

	initialized = true
	logger.Jaeger.Info("jaeger tracing initialized",
		zap.String("service_name", cfg.ServiceName),
		zap.String("agent_endpoint", cfg.AgentEndpoint),
		zap.String("collector_url", cfg.CollectorURL),
		zap.String("sampler_type", cfg.SamplerType),
		zap.Float64("sampler_param", cfg.SamplerParam),
	)

	return nil
}

// Close closes the Jaeger tracer
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if !initialized || tracerProvider == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tracerProvider.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shutdown tracer provider")
	}

	initialized = false
	logger.Jaeger.Info("jaeger tracer closed")
	return nil
}

// GetTracer returns the global tracer
func GetTracer() trace.Tracer {
	if !initialized {
		return trace.NewNoopTracerProvider().Tracer("noop")
	}
	return tracer
}

// IsEnabled returns whether Jaeger tracing is enabled
func IsEnabled() bool {
	return config.App.Jaeger.Enable && initialized
}

// StartSpan starts a new span with the given name and options
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if !IsEnabled() {
		return ctx, trace.SpanFromContext(ctx)
	}
	return tracer.Start(ctx, name, opts...)
}

// SpanFromContext returns the span from the context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// createJaegerExporter creates a Jaeger exporter based on configuration
func createJaegerExporter(cfg config.Jaeger) (sdktrace.SpanExporter, error) {
	if cfg.CollectorURL != "" {
		// Use HTTP collector
		return jaeger.New(jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(cfg.CollectorURL),
		))
	}

	// Use UDP agent
	return jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(cfg.AgentEndpoint),
	))
}

// createSampler creates a sampler based on configuration
func createSampler(cfg config.Jaeger) sdktrace.Sampler {
	switch cfg.SamplerType {
	case "const":
		if cfg.SamplerParam >= 1.0 {
			return sdktrace.AlwaysSample()
		}
		return sdktrace.NeverSample()
	case "probabilistic":
		return sdktrace.TraceIDRatioBased(cfg.SamplerParam)
	case "ratelimiting":
		// Note: OpenTelemetry doesn't have built-in rate limiting sampler
		// This would need to be implemented separately
		return sdktrace.TraceIDRatioBased(cfg.SamplerParam)
	default:
		return sdktrace.AlwaysSample()
	}
}

// AddSpanTags adds tags to the current span
func AddSpanTags(span trace.Span, tags map[string]interface{}) {
	if span == nil || !span.IsRecording() {
		return
	}

	for key, value := range tags {
		switch v := value.(type) {
		case string:
			span.SetAttributes(attribute.String(key, v))
		case int:
			span.SetAttributes(attribute.Int(key, v))
		case int64:
			span.SetAttributes(attribute.Int64(key, v))
		case float64:
			span.SetAttributes(attribute.Float64(key, v))
		case bool:
			span.SetAttributes(attribute.Bool(key, v))
		default:
			span.SetAttributes(attribute.String(key, "unsupported_type"))
		}
	}
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	if span == nil || !span.IsRecording() {
		return
	}
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordError records an error in the current span
func RecordError(span trace.Span, err error) {
	if span == nil || !span.IsRecording() || err == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}
