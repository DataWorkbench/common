package grpcwrap

import (
	"context"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc"
)

// basicUnaryClientInterceptor do validate the argument and print log.
func basicUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		lg := glog.FromContext(ctx)

		lg.Debug().String("unary invoker", method).RawString("request", pbMsgToString(lg, req)).Fire()

		// Validated request parameters
		if err := validateRequestArgument(req, lg); err != nil {
			return err
		}

		// Invoker to the grpc server
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

// basicStreamClientInterceptor do print log.
func basicStreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (stream grpc.ClientStream, err error) {
		lg := glog.FromContext(ctx)

		lg.Debug().String("stream invoker", method).Bool("ClientStreams", desc.ClientStreams).Bool("ServerStreams", desc.ServerStreams).Fire()

		stream, err = streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			lg.Error().Error("invoker error", err).Fire()
			return nil, err
		}

		lg.Debug().String("stream done", method).Fire()

		return stream, nil
	}
}
