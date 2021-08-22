package rediswrap

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"

	"github.com/DataWorkbench/common/gtrace"
)

var traceComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "redis"}

// hookTrace implements redis.Hook to trace span.
type hookTrace struct {
	tracer gtrace.Tracer
}

func (h *hookTrace) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}
	span := h.tracer.StartSpan(
		"Redis"+strings.ToUpper(cmd.Name()),
		opentracing.ChildOf(parentCtx),
		ext.SpanKindRPCClient,
		traceComponentTag,
	)
	span.LogFields(tracerLog.String("cmd", cmd.String()))
	return opentracing.ContextWithSpan(ctx, span), nil
}

func (h *hookTrace) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if err := cmd.Err(); err != nil {
		ext.Error.Set(span, true)
		span.LogFields(tracerLog.Error(err))
	}
	span.Finish()
	return nil
}

func (h *hookTrace) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}
	span := h.tracer.StartSpan(
		"RedisPipeline",
		opentracing.ChildOf(parentCtx),
		ext.SpanKindRPCClient,
		traceComponentTag,
	)
	return opentracing.ContextWithSpan(ctx, span), nil
}

func (h *hookTrace) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	for i := 0; i < len(cmds); i++ {
		if err := cmds[i].Err(); err != nil {
			ext.Error.Set(span, true)
			span.LogFields(tracerLog.Error(err))
		}
	}
	span.Finish()
	return nil
}
