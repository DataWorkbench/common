package grpcwrap

import (
	"context"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc/metadata"
)

const (
	ctxReqIdKey = "rid"
)

// ContextWithRequest set "*glo.Logger" into context.Context and set "reqId" into
// grpc outgoing metadata
func ContextWithRequest(ctx context.Context, l *glog.Logger, reqId string) context.Context {
	if l != nil {
		ctx = glog.WithContext(ctx, l)
	}
	if reqId != "" {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(ctxReqIdKey, reqId))
	}
	return ctx
}

func reqIdFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ids := md.Get(ctxReqIdKey)
		if len(ids) != 0 {
			return ids[0]
		}
	}
	return "none"
}
