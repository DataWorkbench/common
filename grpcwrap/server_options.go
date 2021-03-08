package grpcwrap

import "github.com/opentracing/opentracing-go"

// ClientOption is a function that sets some option on the grpc client.
type ServerOption func(o *ServerOptions)

// ClientOptions control behavior of the client.
type ServerOptions struct {
	tracer opentracing.Tracer
}

func applyServerOptions(options ...ServerOption) ServerOptions {
	opts := ServerOptions{}
	for _, option := range options {
		option(&opts)
	}
	if opts.tracer == nil {
		opts.tracer = opentracing.NoopTracer{}
	}
	return opts
}

func ServerWithTracer(tracer opentracing.Tracer) ServerOption {
	return func(o *ServerOptions) {
		o.tracer = tracer
	}
}
