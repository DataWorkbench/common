package gormwrap

import "github.com/opentracing/opentracing-go"

// Option is a function that sets some option on the gorm client.
type Option func(o *Options)

// Options control behavior of the client.
type Options struct {
	tracer opentracing.Tracer
}

func WithTracer(tracer opentracing.Tracer) Option {
	return func(o *Options) {
		o.tracer = tracer
	}
}

func applyOptions(options ...Option) Options {
	opts := Options{}
	for _, option := range options {
		option(&opts)
	}
	if opts.tracer == nil {
		opts.tracer = opentracing.NoopTracer{}
	}
	return opts
}
