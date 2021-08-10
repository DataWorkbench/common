package grpcwrap

import "github.com/opentracing/opentracing-go"

// ClientOption is a function that sets some option on the grpc client.
// Deprecated:
type ClientOption func(o *clientOptions)

// clientOptions control behavior of the client.
// Deprecated:
type clientOptions struct {
	tracer opentracing.Tracer
}

// Deprecated: Store opentracing.Tracer to "context.Context" by "gtrace.ContextWithTracer"
func ClientWithTracer(tracer opentracing.Tracer) ClientOption {
	return func(o *clientOptions) {
		o.tracer = tracer
	}
}
