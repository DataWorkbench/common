package grpcwrap

import (
	"context"
	"runtime"

	"github.com/DataWorkbench/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DataWorkbench/common/gtrace"
)

// traceUnaryServerInterceptor for trace the request.
// - extract trace id from incoming metadata and store it to context.
// - creates an new logger object with trace id and store it to context.
// - validate argument where are request and reply.
func traceUnaryServerInterceptor(lp *glog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		tid := extractTraceContext(ctx)
		nl := lp.Clone()
		nl.WithFields().AddString(gtrace.IdKey, tid)

		nl.Debug().String("unary receive", info.FullMethod).RawString("request", pbMsgToString(nl, req)).Fire()

		// Validated request parameters
		if err := validateRequestArgument(req, nl); err != nil {
			return nil, err
		}

		ctx = gtrace.ContextWithId(ctx, tid)
		ctx = glog.WithContext(ctx, nl)
		resp, err = handler(ctx, req)
		if err != nil {
			nl.Error().Error("unary handle error", err).Fire()
		} else {
			nl.Debug().RawString("unary reply", pbMsgToString(nl, resp)).Fire()
		}

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
				lg.Error().RawString("error stack trace", string(buf[0:n])).Fire()

				err = status.Errorf(codes.Internal, "unary server panic: %v", r)
			}
		}()

		resp, err = handler(ctx, req)
		panicked = false
		return
	}
}

// wraps for override the context.
type serverStreamWrap struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *serverStreamWrap) Context() context.Context {
	return s.ctx
}

// traceStreamServerInterceptor for trace the request.
// - extract trace id from incoming metadata and store it to context.
// - creates an new logger object with trace id and store it to context.
func traceStreamServerInterceptor(lp *glog.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		tid := extractTraceContext(ctx)

		nl := lp.Clone()
		nl.WithFields().AddString(gtrace.IdKey, tid)

		nl.Debug().
			String("stream receive", info.FullMethod).
			Bool("ClientStream", info.IsClientStream).
			Bool("ServerStream", info.IsServerStream).
			Fire()

		ctx = gtrace.ContextWithId(ctx, tid)
		ctx = glog.WithContext(ctx, nl)
		err := handler(srv, &serverStreamWrap{ServerStream: ss, ctx: ctx})
		if err != nil {
			nl.Error().Error("stream error", err).Fire()
		} else {
			nl.Debug().String("stream done", info.FullMethod).Fire()
		}

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
				lg.Error().RawString("error stack trace", string(buf[0:n])).Fire()

				err = status.Errorf(codes.Internal, "stream server panic: %v", r)
			}
		}()

		err = handler(srv, ss)
		panicked = false
		return
	}
}
