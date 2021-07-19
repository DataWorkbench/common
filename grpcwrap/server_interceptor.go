package grpcwrap

import (
	"context"
	"runtime"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ctxUnaryServerInterceptor create an new logger object with requestId.
// You can get logger by glog.FromContext(cxt) after.
func ctxUnaryServerInterceptor(lp *glog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		reqId := ReqIdFromContext(ctx)

		// Copy a new logger
		nl := lp.Clone()
		nl.WithFields().AddString(ctxReqIdKey, reqId)

		ctx = ContextWithRequest(ctx, nl, reqId)

		resp, err = handler(ctx, req)

		// Close the logger instances
		_ = nl.Close()
		return resp, err
	}
}

// recoverUnaryServerInterceptor returns a new unary server interceptor for panic recovery.
func recoverUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		panicked := true
		defer func() {
			if r := recover(); r != nil || panicked {
				lg := glog.FromContext(ctx)
				lg.Error().Any("unary server panic", r).Fire()

				buf := make([]byte, 2048)
				n := runtime.Stack(buf, true)
				lg.Error().RawString("error trace", string(buf[0:n])).Fire()

				err = status.Errorf(codes.Internal, "unary server panic: %v", r)
			}
		}()

		resp, err = handler(ctx, req)
		panicked = false
		return
	}
}

// basicUnaryServerInterceptor do validate the argument and print log
func basicUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		lg := glog.FromContext(ctx)

		lg.Debug().String("unary receive", info.FullMethod).RawString("request", pbMsgToString(lg, req)).Fire()

		// Validated request parameters
		if err := validateRequestArgument(req, lg); err != nil {
			return nil, err
		}
		reply, err := handler(ctx, req)

		if err != nil {
			lg.Error().Error("unary serve error", err).Fire()
			return nil, err
		}

		lg.Debug().RawString("unary reply", pbMsgToString(lg, reply)).Fire()
		return reply, err
	}
}

type ctxServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *ctxServerStream) Context() context.Context {
	return s.ctx
}

// ctxStreamServerInterceptor create an new logger object with requestId.
// You can get logger by glog.FromContext(cxt) after.
func ctxStreamServerInterceptor(lp *glog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()

		reqId := ReqIdFromContext(ctx)

		nl := lp.Clone()
		nl.WithFields().AddString(ctxReqIdKey, reqId)

		ctx = ContextWithRequest(ctx, nl, reqId)

		err := handler(srv, &ctxServerStream{ServerStream: ss, ctx: ctx})

		// Close the logger instances
		_ = nl.Close()
		return err
	}
}

// recoverStreamServerInterceptor returns a new stream server interceptor for panic recovery.
func recoverStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		panicked := true
		defer func() {
			if r := recover(); r != nil || panicked {
				ctx := ss.Context()
				lg := glog.FromContext(ctx)
				lg.Error().Any("stream server panic", r).Fire()

				buf := make([]byte, 2048)
				n := runtime.Stack(buf, true)
				lg.Error().RawString("error trace", string(buf[0:n])).Fire()

				err = status.Errorf(codes.Internal, "stream server panic: %v", r)
			}
		}()

		err = handler(srv, ss)
		panicked = false
		return
	}
}

// basicStreamServerInterceptor do print log.
func basicStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		lg := glog.FromContext(ctx)

		lg.Debug().String("stream receive", info.FullMethod).Bool("ClientStream", info.IsClientStream).Bool("ServerStream", info.IsServerStream).Fire()

		err := handler(srv, ss)
		if err != nil {
			lg.Error().Error("stream server error", err).Fire()
			return err
		}

		lg.Debug().String("stream done", info.FullMethod).Fire()
		return nil
	}
}
