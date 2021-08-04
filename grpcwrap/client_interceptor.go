package grpcwrap

import (
	"context"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc"
)

// traceUnaryClientInterceptor for trace the request.
// - inject trace id to outgoing metadata.
// - validate argument where are request and reply.
func traceUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		lg := glog.FromContext(ctx)

		lg.Debug().String("unary invoker", method).RawString("request", pbMsgToString(lg, req)).Fire()

		// Validated request parameters
		if err := validateRequestArgument(req, lg); err != nil {
			return err
		}

		ctx = injectTraceContext(ctx)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			lg.Error().Error("invoker error", err).Fire()
			return err
		}

		lg.Debug().RawString("receive reply", pbMsgToString(lg, reply)).Fire()

		// Validated reply parameters
		if err := validateReplyArgument(reply, lg); err != nil {
			return err
		}
		return nil
	}
}

// traceStreamClientInterceptor for trace the request.
// - inject trace id to metadata.
func traceStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		lg := glog.FromContext(ctx)

		lg.Debug().String("stream invoker", method).Bool("ClientStreams", desc.ClientStreams).Bool("ServerStreams", desc.ServerStreams).Fire()

		ctx = injectTraceContext(ctx)
		stream, err = streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			lg.Error().Error("invoker error", err).Fire()
			return nil, err
		}

		lg.Debug().String("stream done", method).Fire()

		return stream, nil
	}
}
