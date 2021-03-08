package grpcwrap

import "github.com/opentracing/opentracing-go"

// ClientOption is a function that sets some option on the grpc client.
type ClientOption func(o *ClientOptions)

// ClientOptions control behavior of the client.
type ClientOptions struct {
	tracer opentracing.Tracer
}

func applyClientOptions(options ...ClientOption) ClientOptions {
	opts := ClientOptions{}
	for _, option := range options {
		option(&opts)
	}
	if opts.tracer == nil {
		opts.tracer = opentracing.NoopTracer{}
	}
	return opts
}

func ClientWithTracer(tracer opentracing.Tracer) ClientOption {
	return func(o *ClientOptions) {
		o.tracer = tracer
	}
}
