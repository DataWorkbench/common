package grpcwrap

import "github.com/opentracing/opentracing-go"

// ClientOption is a function that sets some option on the grpc client.
type ClientOption func(o *clientOptions)

// clientOptions control behavior of the client.
type clientOptions struct {
	tracer opentracing.Tracer
}

func applyClientOptions(options ...ClientOption) clientOptions {
	opts := clientOptions{}
	for _, option := range options {
		option(&opts)
	}
	if opts.tracer == nil {
		opts.tracer = opentracing.NoopTracer{}
	}
	return opts
}

func ClientWithTracer(tracer opentracing.Tracer) ClientOption {
	return func(o *clientOptions) {
		o.tracer = tracer
	}
}
