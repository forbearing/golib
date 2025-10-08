package util

import (
	"context"

	"github.com/forbearing/gst/provider/otel"
	"go.opentelemetry.io/otel/codes"
)

// TraceFunction is a helper to trace function execution
func TraceFunction(ctx context.Context, functionName string, fn func(context.Context) error) error {
	ctx, span := otel.StartSpan(ctx, functionName)
	defer span.End()

	err := fn(ctx)
	if err != nil {
		otel.RecordError(span, err)
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

// TraceFunctionWithResult is a helper to trace function execution with result
func TraceFunctionWithResult[T any](ctx context.Context, functionName string, fn func(context.Context) (T, error)) (T, error) {
	ctx, span := otel.StartSpan(ctx, functionName)
	defer span.End()

	result, err := fn(ctx)
	if err != nil {
		otel.RecordError(span, err)
		return result, err
	}

	span.SetStatus(codes.Ok, "")
	return result, nil
}
