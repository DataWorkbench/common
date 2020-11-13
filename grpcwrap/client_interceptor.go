package grpcwrap

import (
	"context"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc"
)

// basicUnaryClientInterceptor do validate the argument and print log
func basicUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logger := glog.FromContext(ctx)

		logger.Debug().String("invoker method", method).RawString("request", pbMsgToString(logger, req)).Fire()

		// Validated request parameters
		if err := validateRequestArgument(req, logger); err != nil {
			return err
		}

		// Invoker to the grpc server
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			logger.Error().Error("invoker error", err).Fire()
			return err
		}

		logger.Debug().RawString("receive reply", pbMsgToString(logger, reply)).Fire()

		// Validated reply parameters
		if err := validateReplyArgument(reply, logger); err != nil {
			return err
		}
		return nil
	}
}
