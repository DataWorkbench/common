package grpcwrap

import "github.com/opentracing/opentracing-go"

// ServerOption is a function that sets some option on the grpc client.
// Deprecated:
type ServerOption func(o *serverOptions)

// serverOptions control behavior of the client.
// Deprecated:
type serverOptions struct {
	tracer opentracing.Tracer
}

// Deprecated: Store opentracing.Tracer to "context.Context" by "gtrace.ContextWithTracer"
func ServerWithTracer(tracer opentracing.Tracer) ServerOption {
	return func(o *serverOptions) {
		o.tracer = tracer
	}
}
