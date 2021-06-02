package grpcwrap

import (
	"context"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc/metadata"
)

const (
	ctxReqIdKey = "rid"
)

// ContextWithRequest set "*glog.Logger" into context.Context and set "reqId" into
// grpc outgoing metadata
func ContextWithRequest(ctx context.Context, l *glog.Logger, reqId string) context.Context {
	if l == nil {
		panic("grpcwrap:ContextWithRequest: logger is nil")
	}
	if reqId == "" {
		panic("grpcwrap:ContextWithRequest: request id is nil")
	}

	// Insert logger to context
	ctx = glog.WithContext(ctx, l)

	// Insert request id to context by grpc metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	md.Set(ctxReqIdKey, reqId)

	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}

// ReqIdFromContext get the request id from rpc request context.
func ReqIdFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ids := md.Get(ctxReqIdKey)
		if len(ids) != 0 {
			return ids[0]
		}
	}
	return "none"
}
