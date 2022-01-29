package grpcwrap

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/DataWorkbench/common/gtrace"
)

// injectTraceContext inject the trace id to gRPC metadata.
func injectTraceContext(ctx context.Context) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	tid := gtrace.IdFromContext(ctx)
	if tid == "" {
		return ctx
	}
	md.Set(gtrace.IdKey, tid)
	return metadata.NewOutgoingContext(ctx, md)
}

// extractTraceContext extract the trace id from gRPC metadata.
func extractTraceContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	ids := md.Get(gtrace.IdKey)
	if len(ids) != 0 {
		return ids[0]
	}
	return ""
}
