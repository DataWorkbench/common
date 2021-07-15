package grpcwrap

import (
	"context"
	"net"
	"time"

	"github.com/DataWorkbench/glog"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// GServer is an type aliases to make caller don't have to introduce "google.golang.org/grpc"
// into its project go.mod files. And can keep the grpc version in other project consistent
// with this library at all times.
type GServer = grpc.Server

// ServerConfig used to create an new grpc server
type ServerConfig struct {
	// Listening address of the grpc server.
	Address string `json:"address" yaml:"address" env:"ADDRESS" validate:"required"`
}

// Server is an wrapper for gRPC server.
type Server struct {
	lp   *glog.Logger // the parent logger
	cfg  *ServerConfig
	gRPC *grpc.Server
}

// NewServer return a new Server
// NOTICE: Must set glog.loggerT into the ctx by glow.WithContext
func NewServer(ctx context.Context, cfg *ServerConfig, options ...ServerOption) (s *Server, err error) {
	opts := applyServerOptions(options...)
	lp := glog.FromContext(ctx)

	defer func() {
		if err != nil {
			lp.Error().Error("create grpc server error", err).Fire()
		}
	}()

	var srvOpts []grpc.ServerOption
	// Set and add keepalive server parameters
	// TODO: set keepalive parameters by config
	srvOpts = append(srvOpts, grpc.KeepaliveParams(
		keepalive.ServerParameters{
			MaxConnectionIdle:     time.Second * 30,
			MaxConnectionAge:      time.Second * 30,
			MaxConnectionAgeGrace: time.Second * 30,
			Time:                  time.Second * 1,
			Timeout:               time.Second * 10,
		}))

	// Set and add keepalive enforcement policy
	// TODO: set keepalive parameters by config
	srvOpts = append(srvOpts, grpc.KeepaliveEnforcementPolicy(
		keepalive.EnforcementPolicy{
			MinTime:             time.Second * 10,
			PermitWithoutStream: true,
		}))

	// Set and add Unary Server Interceptor
	srvOpts = append(srvOpts, grpc.ChainUnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(opts.tracer),
		loggerUnaryServerInterceptor(lp),
		recoverUnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		basicUnaryServerInterceptor(),
	))

	// TODO: Impls and add Stream Server Interceptor

	s = &Server{
		lp:   lp,
		cfg:  cfg,
		gRPC: grpc.NewServer(srvOpts...),
	}

	// Register the health server that used by k8s health probe.
	grpc_health_v1.RegisterHealthServer(s.gRPC, health.NewServer())

	return s, nil
}

// Register registers a service and its implementation to gRPC server.
// It is called from the IDL generated code. This must be called before
// invoking Serve.
func (s *Server) Register(f func(s *GServer)) {
	f(s.gRPC)
}

// ListenAndServe creates an net listener by config and called  grpc.Server.Serve
func (s *Server) ListenAndServe() error {
	s.lp.Info().String("gRPC server listening", s.cfg.Address).Fire()

	lis, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		s.lp.Error().Error("gRPC server create listen error", err).Fire()
		return err
	}

	reflection.Register(s.gRPC)
	grpc_prometheus.Register(s.gRPC)

	err = s.gRPC.Serve(lis)
	if err != nil {
		s.lp.Error().Error("gRPC server serve error", err).Fire()
	}
	return err
}

// GracefulStop wrapper for grpc.Server.GracefulStop
func (s *Server) GracefulStop() {
	if s == nil {
		return
	}
	s.lp.Info().Msg("waiting for gRPC server stop").Fire()
	s.gRPC.GracefulStop()
	s.lp.Info().Msg("gRPC server stopped").Fire()
}
