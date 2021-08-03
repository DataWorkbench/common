package grpcwrap

import (
	"context"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc/metadata"

	"github.com/DataWorkbench/common/trace"
)

// injectTraceContext inject the trace id to gRPC metadata.
func injectTraceContext(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	tid := trace.IdFromContext(ctx)
	if tid == "" {
		return ctx
	}
	md.Set(trace.IdKey, tid)
	return metadata.NewOutgoingContext(ctx, md)
}

// extractTraceContext extract the trace id from gRPC metadata.
func extractTraceContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	ids := md.Get(trace.IdKey)
	if len(ids) != 0 {
		return ids[0]
	}
	return ""
}

// ContextWithRequest set "*glog.Logger" into context.Context and set "reqId" into
// grpc outgoing metadata
//
// Deprecated: use trace.ContextWithId(ctx, tid) and glog.WithContext(ctx, nl) instead.
//
func ContextWithRequest(ctx context.Context, l *glog.Logger, reqId string) context.Context {
	if l == nil {
		panic("grpcwrap:ContextWithRequest: logger is nil")
	}
	if reqId == "" {
		panic("grpcwrap:ContextWithRequest: request id is nil")
	}
	ctx = glog.WithContext(ctx, l)
	ctx = trace.ContextWithId(ctx, reqId)
	return ctx
}
