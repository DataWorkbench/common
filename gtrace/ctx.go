package gtrace

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

const (
	// The key string format of trace id.
	IdKey = "tid"

	// The request id key in http header.
	HeaderKey = "X-Request-Id"

	// Used when no trace id found or generated failed.
	//DefaultTraceIdValue = "none"
)

type ctxIdKey struct{}

// ContextWithId store the trace id in context.Value.
func ContextWithId(ctx context.Context, tid string) context.Context {
	if tid == "" {
		return ctx
	}
	// Do not store duplicate id.
	if _id, ok := ctx.Value(ctxIdKey{}).(string); ok && _id == tid {
		return ctx
	}
	return context.WithValue(ctx, ctxIdKey{}, tid)
}

// IdFromContext get the trace id from context.Value.
func IdFromContext(ctx context.Context) (id string) {
	var ok bool
	id, ok = ctx.Value(ctxIdKey{}).(string)
	if !ok {
		id = ""
	}
	return
}

type ctxTracerKey struct{}

// ContextWithTracer store the tracer instances to context.Value.
func ContextWithTracer(ctx context.Context, tracer Tracer) context.Context {
	if tracer == nil {
		return ctx
	}
	// Do not store duplicate id.
	if _tracer, ok := ctx.Value(ctxTracerKey{}).(Tracer); ok && _tracer == tracer {
		return ctx
	}
	return context.WithValue(ctx, ctxTracerKey{}, tracer)
}

// TracerFromContext get the tracer instances from context.Value.
func TracerFromContext(ctx context.Context) (tracer Tracer) {
	var ok bool
	tracer, ok = ctx.Value(ctxTracerKey{}).(Tracer)
	if !ok {
		tracer = opentracing.NoopTracer{}
	}
	return
}
