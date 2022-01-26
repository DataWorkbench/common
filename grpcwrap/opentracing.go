package grpcwrap

import (
	"fmt"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var excludeTraceMethod = map[string]bool{
	fmt.Sprintf("/%s/Check", grpc_health_v1.Health_ServiceDesc.ServiceName): true,
	fmt.Sprintf("/%s/Watch", grpc_health_v1.Health_ServiceDesc.ServiceName): true,
}

// traceSpanInclusionFunc is used to filter grpc method that don't need to be traced.
// Is a type of otgrpc.SpanInclusionFunc
func traceSpanInclusionFunc(parentSpanCtx opentracing.SpanContext, method string, req, resp interface{}) bool {
	return !excludeTraceMethod[method]
}
