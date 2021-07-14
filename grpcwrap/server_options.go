package grpcwrap

import "github.com/opentracing/opentracing-go"

// ServerOption is a function that sets some option on the grpc client.
type ServerOption func(o *serverOptions)

// serverOptions control behavior of the client.
type serverOptions struct {
	tracer opentracing.Tracer
}

func applyServerOptions(options ...ServerOption) serverOptions {
	opts := serverOptions{}
	for _, option := range options {
		option(&opts)
	}
	if opts.tracer == nil {
		opts.tracer = opentracing.NoopTracer{}
	}
	return opts
}

func ServerWithTracer(tracer opentracing.Tracer) ServerOption {
	return func(o *serverOptions) {
		o.tracer = tracer
	}
}
