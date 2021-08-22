package gormwrap

import "github.com/opentracing/opentracing-go"

// Option is a function that sets some option on the gorm client.
// Deprecated:
type Option func(o *Options)

// Options control behavior of the client.
// Deprecated:
type Options struct {
	tracer opentracing.Tracer
}

// Deprecated: Store opentracing.Tracer to "context.Context" by "gtrace.ContextWithTracer"
func WithTracer(tracer opentracing.Tracer) Option {
	return func(o *Options) {
		o.tracer = tracer
	}
}
